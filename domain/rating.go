package domain

type OrderRating struct {
	RestaurantId      int `json:"restaurant_id"`
	OrderId           int `json:"order_id"`
	Rating            int `json:"rating"`
	EstimatedWaitTime int `json:"estimated_waiting_time"`
	WaitTime          int `json:"waiting_time"`
}

type RatingResponse struct {
	RestaurantId        int     `json:"restaurant_id"`
	RestaurantAvgRating float64 `json:"restaurant_avg_rating"`
	PreparedOrders      int     `json:"prepared_orders"`
}
