package domain

type Basket struct {
	ID          int64
	UserID      int64
	TotalAmount float64
	Items       []Product
}
