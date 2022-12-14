package domain

import (
	"math/rand"
	"sync/atomic"
	"time"

	"github.com/rs/zerolog/log"
)

const (
	free    = "free"
	ready   = "ready"
	waiting = "waiting"
)

var OrderId int64

type Table struct {
	Id           int
	Menu         Menu
	State        string
	CurrentOrder Order
	SendChan     chan<- Order
	ReceiveChan  chan Order
	RatingChan   chan<- int
}

func NewTable(id int, menu Menu, orderChan chan<- Order, ratingChan chan<- int) *Table {
	return &Table{
		Id:          id,
		Menu:        menu,
		State:       free,
		SendChan:    orderChan,
		ReceiveChan: make(chan Order),
		RatingChan:  ratingChan,
	}
}

func (t *Table) Run() {
	for {
		t.waitFree()
		t.sendOrder()
		t.receiveOrder()
	}
}

func (t *Table) nextState() {
	if t.State == free {
		t.State = ready
	} else if t.State == ready {
		t.State = waiting
	} else if t.State == waiting {
		t.State = free
	}
}

func (t *Table) waitFree() {
	if t.State != free {
		return
	}

	for {
		freeTime := time.Duration(cfg.TimeUnit*(rand.Intn(cfg.MaxTableFreeTime)+1)) * time.Millisecond
		time.Sleep(freeTime)

		if atomic.LoadInt64(&NrOfWaitingFoods) > int64(cfg.MaxNrOfWaitingFoods) {
			log.Warn().Int("table_id", t.Id).Int64("waiting_foods", atomic.LoadInt64(&NrOfWaitingFoods)).Msg("Table is waiting for free kitchen")
			continue
		}

		t.nextState()
		log.Debug().Int("table_id", t.Id).Msg("Table has been occupied")
		return
	}
}

func (t *Table) sendOrder() {
	if t.State != ready {
		return
	}

	foodCount := rand.Intn(cfg.MaxOrderItemsCount) + 1

	for i := 0; i < 3; i++ {
		if foodCount > 5 {
			foodCount = rand.Intn(cfg.MaxOrderItemsCount) + 1
		} else {
			break
		}
	}

	order := Order{
		OrderId: atomic.AddInt64(&OrderId, 1),
		TableId: t.Id,
		Items:   make([]int, foodCount),
	}

	order.Priority = (cfg.MaxOrderItemsCount - foodCount) / (cfg.MaxOrderItemsCount / 5)

	maxTime := 0
	for i := 0; i < foodCount; i++ {
		order.Items[i] = rand.Intn(t.Menu.FoodsCount) + 1
		prepTime := t.Menu.Foods[i].PreparationTime
		if prepTime > maxTime {
			maxTime = prepTime
		}
	}

	order.MaxWait = float64(maxTime) * cfg.MaxWaitTimeCoefficient

	t.CurrentOrder = order
	t.SendChan <- order
	t.nextState()

	log.Debug().Int("table_id", t.Id).Int64("order_id", order.OrderId).Msg("Table placed new order")
}

func (t *Table) receiveOrder() {
	if t.State != waiting {
		return
	}

	for order := range t.ReceiveChan {
		if order.TableId != t.Id || order.OrderId != t.CurrentOrder.OrderId {
			log.Err(nil).Int("table_id", t.Id).Int64("order_id", order.OrderId).Msg("Table received wrong order")
			continue
		}

		rating := order.CalculateRating()
		t.RatingChan <- rating
		t.nextState()

		log.Debug().Int("table_id", t.Id).Int64("order_id", order.OrderId).Int("rating", rating).Msg("Table received order")
		return
	}
}
