package domain

import "time"

const (
	timeUnit = 250
)

type Order struct {
	OrderId    int       `json:"order_id"`
	TableId    int       `json:"table_id"`
	WaiterId   int       `json:"waiter_id"`
	Items      []int     `json:"items"`
	Priority   int       `json:"priority"`
	MaxWait    float64   `json:"max_wait"`
	PickUpTime time.Time `json:"pick_up_time"`
}

func (o Order) CalculateRating() int {
	orderTime := float64(time.Since(o.PickUpTime))
	maxWaitTime := o.MaxWait * float64(timeUnit) * float64(time.Millisecond)

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
