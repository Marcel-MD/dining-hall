package domain

type Restaurant struct {
	RestaurantId int     `json:"restaurant_id"`
	Name         string  `json:"name"`
	Address      string  `json:"address"`
	MenuItems    int     `json:"menu_items"`
	Menu         []Food  `json:"menu"`
	Rating       float64 `json:"rating"`
}
