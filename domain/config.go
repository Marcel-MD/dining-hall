package domain

type Config struct {
	TimeUnit    int `json:"time_unit"`
	NrOfTables  int `json:"nr_of_tables"`
	NrOfWaiters int `json:"nr_of_waiters"`

	MaxOrderItemsCount     int     `json:"max_order_items_count"`
	MaxTableFreeTime       int     `json:"max_table_free_time"`
	MaxWaitTimeCoefficient float64 `json:"max_wait_time_coefficient"`
	MaxPickupTime          int     `json:"max_pickup_time"`

	KitchenUrl      string `json:"kitchen_url"`
	FoodOrderingUrl string `json:"food_ordering_url"`
	DiningHallPort  string `json:"dining_hall_port"`
	DiningHallUrl   string `json:"dining_hall_url"`

	RestaurantName string `json:"restaurant_name"`
	RestaurantId   int    `json:"restaurant_id"`
}

var cfg Config = Config{
	TimeUnit:    250,
	NrOfTables:  10,
	NrOfWaiters: 4,

	MaxOrderItemsCount:     10,
	MaxTableFreeTime:       20,
	MaxWaitTimeCoefficient: 1.3,
	MaxPickupTime:          5,

	KitchenUrl:      "http://kitchen:8081",
	FoodOrderingUrl: "http://food-ordering:8090",
	DiningHallPort:  "8080",
	DiningHallUrl:   "http://dining-hall:8080",

	RestaurantName: "Mujik's Pizza",
	RestaurantId:   1,
}

func SetConfig(c Config) {
	cfg = c
}
