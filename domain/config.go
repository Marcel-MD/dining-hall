package domain

type Config struct {
	TimeUnit    int `json:"time_unit"`
	NrOfTables  int `json:"nr_of_tables"`
	NrOfWaiters int `json:"nr_of_waiters"`

	MaxOrderItemsCount     int     `json:"max_order_items_count"`
	MaxTableFreeTime       int     `json:"max_table_free_time"`
	MaxWaitTimeCoefficient float64 `json:"max_wait_time_coefficient"`
	MaxPickupTime          int     `json:"max_pickup_time"`

	KitchenUrl string `json:"kitchen_url"`
}

var cfg Config = Config{
	TimeUnit:    1000,
	NrOfTables:  10,
	NrOfWaiters: 4,

	MaxOrderItemsCount:     5,
	MaxTableFreeTime:       20,
	MaxWaitTimeCoefficient: 1.3,
	MaxPickupTime:          5,

	KitchenUrl: "http://kitchen:8081",
}

func SetConfig(c Config) {
	cfg = c
}
