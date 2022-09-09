package domain

import (
	"bytes"
	"encoding/json"
	"net/http"
	"time"

	"github.com/rs/zerolog/log"
)

const (
	kitchenPath = "http://localhost:8081/order"
)

type Waiter struct {
	Id               int
	CurrentOrder     Order
	DistributionChan chan Distribution
	OrderChan        <-chan Order
	TablesChans      []chan Order
}

func NewWaiter(id int, orderChan <-chan Order, tablesChans []chan Order) Waiter {
	return Waiter{
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
			order.PickUpTime = time.Now()
			order.WaiterId = w.Id

			jsonBody, err := json.Marshal(order)
			if err != nil {
				log.Fatal().Err(err).Msg("Error marshalling order")
			}
			contentType := "application/json"

			log.Info().Int("waiter_id", w.Id).Int("order_id", order.OrderId).Msg("Received order from table")

			_, err = http.Post(kitchenPath, contentType, bytes.NewReader(jsonBody))
			if err != nil {
				log.Fatal().Err(err).Msg("Error sending order to kitchen")
			}

		case distribution := <-w.DistributionChan:
			order := distribution.Order
			log.Info().Int("waiter_id", w.Id).Int("order_id", order.OrderId).Msg("Received distribution from chef")
			w.TablesChans[order.TableId] <- order
		}
	}
}
