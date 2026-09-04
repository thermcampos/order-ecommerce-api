package products

type createProductDto struct {
	Name         string `json:"name"`
	PriceInCents int32  `json:"priceInCents"`
	Quantity     int32  `json:"quantity"`
}
