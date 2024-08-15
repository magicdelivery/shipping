package model

type Customer struct {
	ID        string           `json:"id"`
	Name      string           `json:"name"`
	Address   *ShippingAddress `json:"address"`
	CreatedAt int64            `json:"created_at"`
}

type ShippingAddress struct {
	ID        string  `json:"id"`
	City      string  `json:"city"`
	Street    string  `json:"street"`
	Longitude float64 `json:"longitude"`
	Latitude  float64 `json:"latitude"`
}
