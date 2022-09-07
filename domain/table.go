package domain

import (
	"math/rand"

	"github.com/Marcel-MD/dining-hall/dto"
	"github.com/rs/zerolog/log"
)

const (
	free    = "free"
	ready   = "ready"
	waiting = "waiting"
)

const (
	maxFoodCount     = 5
	orderProbability = 0.25
)

type Table struct {
	Id           int
	Menu         Menu
	State        string
	CurrentOrder dto.Order
	SendChan     chan<- dto.Order
	ReceiveChan  <-chan dto.Order
	RatingChan   chan<- int
}

func NewTable(id int, menu Menu, orderChan chan<- dto.Order, ratingChan chan<- int) Table {
	return Table{
		Id:          id,
		Menu:        menu,
		State:       free,
		SendChan:    orderChan,
		ReceiveChan: make(<-chan dto.Order),
		RatingChan:  ratingChan,
	}
}

func (t *Table) NextState() {
	if t.State == free {
		t.State = ready
	} else if t.State == ready {
		t.State = waiting
	} else if t.State == waiting {
		t.State = free
	}
}

func (t *Table) SendOrder() {
	if t.State != ready {
		return
	}

	foodCount := rand.Intn(maxFoodCount)

	order := dto.Order{
		OrderId:  rand.Intn(1000) + 1,
		TableId:  t.Id,
		Items:    make([]int, foodCount),
		Priority: maxFoodCount - foodCount,
	}

	maxTime := 0
	for i := 0; i < foodCount; i++ {
		order.Items[i] = rand.Intn(t.Menu.FoodsCount) + 1
		prepTime := t.Menu.Foods[i].PreparationTime
		if prepTime > maxTime {
			maxTime = prepTime
		}
	}

	order.MaxWait = float64(maxTime) * 1.3

	t.CurrentOrder = order
	t.NextState()
	t.SendChan <- order

	log.Info().Int("table_id", t.Id).Int("order_id", order.OrderId).Msg("Sent order")
}

func (t *Table) ReceiveOrder() {

	for order := range t.ReceiveChan {
		if order.TableId != t.Id || order.OrderId != t.CurrentOrder.OrderId {
			log.Err(nil).Int("table_id", t.Id).Int("order_id", order.OrderId).Msg("Received wrong order")
			continue
		}

		// TODO: calculate rating
	}
}
