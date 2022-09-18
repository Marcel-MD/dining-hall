package domain

import (
	"time"

	"github.com/rs/zerolog/log"
)

type Order struct {
	OrderId    int64   `json:"order_id"`
	TableId    int     `json:"table_id"`
	WaiterId   int     `json:"waiter_id"`
	Items      []int   `json:"items"`
	Priority   int     `json:"priority"`
	MaxWait    float64 `json:"max_wait"`
	PickUpTime int64   `json:"pick_up_time"`
}

func (o Order) CalculateRating() int {
	orderTime := float64((time.Now().Unix() - o.PickUpTime) * 1000 / int64(cfg.TimeUnit))
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
