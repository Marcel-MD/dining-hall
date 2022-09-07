package dto

type Distribution struct {
	Order

	CookingTime    int             `json:"cooking_time"`
	CookingDetails []CookingDetail `json:"cooking_details"`
}
