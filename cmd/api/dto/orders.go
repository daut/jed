package dto

type OrderCreateRequest struct {
	FirstName string         `json:"first_name"`
	LastName  string         `json:"last_name"`
	Email     string         `json:"email"`
	Phone     string         `json:"phone"`
	Address   string         `json:"address"`
	City      string         `json:"city"`
	Products  []OrderProduct `json:"products"`
}

type OrderProduct struct {
	ProductID int32 `json:"product_id"`
	Quantity  int32 `json:"quantity"`
}
