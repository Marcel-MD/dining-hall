package domain

type Distribution struct {
	Order

	CookingTime    int             `json:"cooking_time"`
	CookingDetails []CookingDetail `json:"cooking_details"`
}

type CookingDetail struct {
	FoodId int `json:"food_id"`
	CookId int `json:"cook_id"`
}

type DistributionResponse struct {
	OrderId        int64           `json:"order_id"`
	IsReady        bool            `json:"is_ready"`
	EstimatedWait  float64         `json:"estimated_waiting_time"`
	Priority       int             `json:"priority"`
	MaxWait        float64         `json:"max_wait"`
	CreatedTime    int64           `json:"created_time"`
	RegisteredTime int64           `json:"registered_time"`
	PreparedTime   int64           `json:"prepared_time"`
	CookingTime    int             `json:"cooking_time"`
	CookingDetails []CookingDetail `json:"cooking_details"`
}
