package domain

type Cart struct {
	ID          int64
	UserID      int64
	TotalAmount float64
	Items       []CartItem
	IsActive    bool
}

type CartItem struct {
	ID        int64
	CartID    int64
	ProductID int64
	Quantity  int
	Price     float64
}
