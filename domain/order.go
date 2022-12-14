package domain

import (
	"time"

	"github.com/rs/zerolog/log"
)

type Order struct {
	OrderId     int64   `json:"order_id"`
	TableId     int     `json:"table_id"`
	WaiterId    int     `json:"waiter_id"`
	Items       []int   `json:"items"`
	Priority    int     `json:"priority"`
	MaxWait     float64 `json:"max_wait"`
	PickUpTime  int64   `json:"pick_up_time"`
	CreatedTime int64   `json:"created_time"`
}

type OrderResponseData struct {
	OrderId        int64   `json:"order_id"`
	RestaurantId   int     `json:"restaurant_id"`
	EstimatedWait  float64 `json:"estimated_waiting_time"`
	CreatedTime    int64   `json:"created_time"`
	RegisteredTime int64   `json:"registered_time"`
}

func (o Order) CalculateRating() int {
	orderTime := float64((time.Now().UnixMilli() - o.PickUpTime) / int64(cfg.TimeUnit))
	maxWaitTime := o.MaxWait

	log.Debug().Int64("order_id", o.OrderId).Float64("order_time", orderTime).Float64("max_wait", maxWaitTime).Msg("Calculating rating")

	if orderTime < maxWaitTime {
		return 5
	}

	if orderTime < maxWaitTime*1.1 {
		return 4
	}

	if orderTime < maxWaitTime*1.2 {
		return 3
	}

	if orderTime < maxWaitTime*1.3 {
		return 2
	}

	if orderTime < maxWaitTime*1.4 {
		return 1
	}

	return 0
}
