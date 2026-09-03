package orders

type orderItem struct {
	ProductID int64 `json:"productId"`
	Quantity  int32 `json:"quantity"`
}

type createOrderParam struct {
	CustomerID int64       `json:"customerId"`
	Items      []orderItem `json:"items"`
}
