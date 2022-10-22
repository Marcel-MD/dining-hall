package main

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"strconv"
	"sync/atomic"
	"time"

	"github.com/Marcel-MD/dining-hall/domain"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

const v2Id = 69

func main() {
	cfg := config()
	domain.SetConfig(cfg)

	menu := domain.GetMenu()
	newOrderChan := make(chan domain.Order)
	ratingChan := make(chan int)
	tablesChans := make([]chan domain.Order, 0)
	waitersChans := make([]chan domain.Distribution, 0)

	distributionAwaitingPickup := make(map[int64]domain.DistributionResponse)

	for i := 0; i < cfg.NrOfTables; i++ {
		table := domain.NewTable(i, menu, newOrderChan, ratingChan)
		tablesChans = append(tablesChans, table.ReceiveChan)
		go table.Run()
	}

	for i := 0; i < cfg.NrOfWaiters; i++ {
		waiter := domain.NewWaiter(i, newOrderChan, tablesChans)
		waitersChans = append(waitersChans, waiter.DistributionChan)
		go waiter.Run()
	}

	go rating(ratingChan)

	// Register restaurant
	restaurant := domain.Restaurant{
		RestaurantId: cfg.RestaurantId,
		Name:         cfg.RestaurantName,
		Address:      cfg.DiningHallUrl,
		MenuItems:    menu.FoodsCount,
		Menu:         menu.Foods,
	}

	jsonBody, err := json.Marshal(restaurant)
	if err != nil {
		log.Fatal().Err(err).Msg("Error marshalling restaurant")
	}
	contentType := "application/json"

	_, err = http.Post(cfg.FoodOrderingUrl+"/register", contentType, bytes.NewReader(jsonBody))
	if err != nil {
		log.Fatal().Err(err).Msg("Error registering restaurant")
	}

	r := mux.NewRouter()
	r.HandleFunc("/distribution", func(w http.ResponseWriter, r *http.Request) {
		var distribution domain.Distribution
		err := json.NewDecoder(r.Body).Decode(&distribution)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		atomic.AddInt64(&domain.NrOfWaitingFoods, -int64(len(distribution.Order.Items)))

		if distribution.WaiterId == v2Id {

			distributionResponse, ok := distributionAwaitingPickup[distribution.OrderId]
			if !ok {
				log.Error().Int64("order_id", distribution.OrderId).Msg("Order not found")
				http.Error(w, "Order not found", http.StatusBadRequest)
				return
			}

			distributionResponse.CookingDetails = distribution.CookingDetails
			distributionResponse.CookingTime = distribution.CookingTime
			distributionResponse.IsReady = true

			distributionAwaitingPickup[distribution.OrderId] = distributionResponse
			log.Info().Int64("order_id", distribution.OrderId).Msg("Distribution ready for pickup")
		} else {
			waiterId := distribution.WaiterId
			waitersChans[waiterId] <- distribution
		}

		w.WriteHeader(http.StatusOK)
	}).Methods("POST")

	// V2 Orders
	r.HandleFunc("/v2/order", func(w http.ResponseWriter, r *http.Request) {
		var order domain.Order
		err := json.NewDecoder(r.Body).Decode(&order)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		order.TableId = v2Id
		order.WaiterId = v2Id
		order.OrderId = atomic.AddInt64(&domain.OrderId, 1)
		order.PickUpTime = time.Now().UnixMilli()

		jsonBody, err := json.Marshal(order)
		if err != nil {
			log.Fatal().Err(err).Msg("Error marshalling order")
		}
		contentType := "application/json"

		_, err = http.Post(cfg.KitchenUrl+"/order", contentType, bytes.NewReader(jsonBody))
		if err != nil {
			log.Fatal().Err(err).Msg("Error sending order to kitchen")
		}

		log.Debug().Int64("order_id", order.OrderId).Msg("v2 Order sent to kitchen")

		atomic.AddInt64(&domain.NrOfWaitingFoods, int64(len(order.Items)))

		wait := estimateWait(cfg, menu, order)
		if wait > order.MaxWait {
			wait = order.MaxWait
		}

		distributionResponse := domain.DistributionResponse{
			OrderId:        order.OrderId,
			IsReady:        false,
			Priority:       order.Priority,
			MaxWait:        order.MaxWait,
			EstimatedWait:  wait,
			CreatedTime:    order.CreatedTime,
			RegisteredTime: time.Now().UnixMilli(),
		}

		distributionAwaitingPickup[order.OrderId] = distributionResponse

		orderResponse := domain.OrderResponseData{
			OrderId:        order.OrderId,
			RestaurantId:   cfg.RestaurantId,
			EstimatedWait:  wait,
			CreatedTime:    order.CreatedTime,
			RegisteredTime: time.Now().UnixMilli(),
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(orderResponse)
	}).Methods("POST")

	r.HandleFunc("/v2/order/{order_id}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		orderIdStr := vars["order_id"]
		orderId, err := strconv.ParseInt(orderIdStr, 10, 64)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		distribution, ok := distributionAwaitingPickup[orderId]
		if !ok {
			http.Error(w, "Distribution not found", http.StatusNotFound)
			log.Warn().Int64("order_id", orderId).Msg("Distribution not found")
			return
		}

		if distribution.IsReady {
			delete(distributionAwaitingPickup, orderId)
		}

		log.Info().Int64("order_id", orderId).Msg("Distribution picked up")

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(distribution)
	}).Methods("GET")

	r.HandleFunc("/v2/rating", func(w http.ResponseWriter, r *http.Request) {
		var rating domain.OrderRating
		err := json.NewDecoder(r.Body).Decode(&rating)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		log.Info().Int("order_id", rating.OrderId).Int("rating", rating.Rating).Msg("Rating received from food ordering")

		ratingChan <- rating.Rating

		ratingResponse := domain.RatingResponse{
			RestaurantId:        cfg.RestaurantId,
			RestaurantAvgRating: avgRating,
			PreparedOrders:      nrOfRatings,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(ratingResponse)
	}).Methods("POST")

	http.ListenAndServe(":"+cfg.DiningHallPort, r)
}

func estimateWait(cfg domain.Config, menu domain.Menu, order domain.Order) float64 {

	a, c := 0, 0

	for item := range order.Items {
		food := menu.Foods[item]

		if food.CookingApparatus != "" {
			c += food.PreparationTime
		} else {
			a += food.PreparationTime
		}
	}

	b, d := cfg.SumProficiencies, cfg.NrApparatuses

	e := int(atomic.LoadInt64(&domain.NrOfWaitingFoods))

	f := len(order.Items)

	wait := float64((a/b + c/d) * (e + f) / f)

	log.Debug().Int64("order_id", order.OrderId).Float64("estimated_wait", wait).Float64("max_wait", order.MaxWait).Msg("Estimated wait")

	return wait
}

var avgRating float64
var nrOfRatings int
var totalRating int

func rating(ratingChan <-chan int) {

	for {
		rating := <-ratingChan
		nrOfRatings++
		totalRating += rating
		avgRating = float64(totalRating) / float64(nrOfRatings)
		log.Info().Int("rating", rating).Float64("avg_rating", avgRating).Msg("Received rating")
	}
}

func config() domain.Config {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	log.Logger = log.With().Caller().Logger()

	file, err := os.Open("config/cfg.json")
	if err != nil {
		log.Fatal().Err(err).Msg("Error opening menu.json")
	}
	defer file.Close()

	byteValue, _ := io.ReadAll(file)
	var cfg domain.Config
	json.Unmarshal(byteValue, &cfg)

	return cfg
}
