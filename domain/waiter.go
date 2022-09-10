package domain

import (
	"bytes"
	"encoding/json"
	"math/rand"
	"net/http"
	"time"

	"github.com/rs/zerolog/log"
)

const (
	kitchenPath   = "http://localhost:8081/order"
	maxPickupTime = 5
)

type Waiter struct {
	Id               int
	CurrentOrder     Order
	DistributionChan chan Distribution
	OrderChan        <-chan Order
	TablesChans      []chan Order
}

func NewWaiter(id int, orderChan <-chan Order, tablesChans []chan Order) *Waiter {
	return &Waiter{
		Id:               id,
		DistributionChan: make(chan Distribution),
		OrderChan:        orderChan,
		TablesChans:      tablesChans,
	}
}

func (w *Waiter) Run() {
	for {
		select {
		case order := <-w.OrderChan:
			pickupTime := time.Duration(timeUnit*(rand.Intn(maxPickupTime)+1)) * time.Millisecond
			time.Sleep(pickupTime)

			order.PickUpTime = time.Now().UnixMilli()
			order.WaiterId = w.Id

			jsonBody, err := json.Marshal(order)
			if err != nil {
				log.Fatal().Err(err).Msg("Error marshalling order")
			}
			contentType := "application/json"

			_, err = http.Post(kitchenPath, contentType, bytes.NewReader(jsonBody))
			if err != nil {
				log.Fatal().Err(err).Msg("Error sending order to kitchen")
			}

			log.Info().Int("waiter_id", w.Id).Int("order_id", order.OrderId).Msg("Waiter sent order to kitchen")

		case distribution := <-w.DistributionChan:
			order := distribution.Order
			log.Info().Int("waiter_id", w.Id).Int("order_id", order.OrderId).Int("cooking_time", distribution.CookingTime).Float64("max_wait", distribution.MaxWait).Msgf("Waiter received distribution")
			w.TablesChans[order.TableId] <- order
		}
	}
}
