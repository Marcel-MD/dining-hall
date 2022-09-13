package main

import (
	"encoding/json"
	"io"
	"net/http"
	"os"

	"github.com/Marcel-MD/dining-hall/domain"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	cfg := config()
	domain.SetConfig(cfg)

	menu := domain.GetMenu()
	newOrderChan := make(chan domain.Order)
	ratingChan := make(chan int)
	tablesChans := make([]chan domain.Order, 0)
	waitersChans := make([]chan domain.Distribution, 0)

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

	r := mux.NewRouter()
	r.HandleFunc("/distribution", func(w http.ResponseWriter, r *http.Request) {
		var distribution domain.Distribution
		err := json.NewDecoder(r.Body).Decode(&distribution)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		waiterId := distribution.WaiterId
		waitersChans[waiterId] <- distribution

		w.WriteHeader(http.StatusOK)
	}).Methods("POST")

	http.ListenAndServe(":8080", r)
}

func rating(ratingChan <-chan int) {
	nrOfRatings := 0
	totalRating := 0

	for {
		rating := <-ratingChan
		nrOfRatings++
		totalRating += rating
		log.Info().Int("rating", rating).Float64("avg_rating", float64(totalRating)/float64(nrOfRatings)).Msg("Received rating")
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
