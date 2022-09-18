package domain

import (
	"bytes"
	"encoding/json"
	"math/rand"
	"net/http"
	"time"

	"github.com/rs/zerolog/log"
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
			pickupTime := time.Duration(cfg.TimeUnit*(rand.Intn(cfg.MaxPickupTime)+1)) * time.Millisecond
			time.Sleep(pickupTime)

			order.PickUpTime = time.Now().Unix()
			order.WaiterId = w.Id

			jsonBody, err := json.Marshal(order)
			if err != nil {
				log.Fatal().Err(err).Msg("Error marshalling order")
			}
			contentType := "application/json"

			_, err = http.Post(cfg.KitchenUrl+"/order", contentType, bytes.NewReader(jsonBody))
			if err != nil {
				log.Fatal().Err(err).Msg("Error sending order to kitchen")
			}

			log.Debug().Int("waiter_id", w.Id).Int64("order_id", order.OrderId).Msg("Waiter sent order to kitchen")

		case distribution := <-w.DistributionChan:
			order := distribution.Order
			log.Debug().Int("waiter_id", w.Id).Int64("order_id", order.OrderId).Int("cooking_time", distribution.CookingTime).Float64("max_wait", distribution.MaxWait).Msgf("Waiter received distribution")
			w.TablesChans[order.TableId] <- order
		}
	}
}
