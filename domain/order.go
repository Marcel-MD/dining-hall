package domain

import "time"

var timeUnit = 1000

func SetTimeUnit(unit int) {
	timeUnit = unit
}

type Order struct {
	OrderId    int     `json:"order_id"`
	TableId    int     `json:"table_id"`
	WaiterId   int     `json:"waiter_id"`
	Items      []int   `json:"items"`
	Priority   int     `json:"priority"`
	MaxWait    float64 `json:"max_wait"`
	PickUpTime int64   `json:"pick_up_time"`
}

func (o Order) CalculateRating() int {
	orderTime := float64(time.Now().UnixMilli() - o.PickUpTime)
	maxWaitTime := o.MaxWait * float64(timeUnit)

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
