package infrastructure

import (
    "database/sql"
    "ecommerce-go/domain"
    "errors"
)

// CartRepository defines CRUD operations for Cart
type CartRepository interface {
    CreateCart(userID int64) (*domain.Cart, error)
    GetCartByUserID(userID int64) (*domain.Cart, error)
    AddItemToCart(cartID int64, productID int64, quantity int, price float64) (*domain.CartItem, error)
    RemoveItemFromCart(cartID int64, productID int64) error
    UpdateCartItem(cartID int64, productID int64, quantity int) (*domain.CartItem, error)
    GetCartItems(cartID int64) ([]*domain.CartItem, error)
    ClearCart(cartID int64) error
    DeleteCart(cartID int64) error
    UpdateCartTotal(cartID int64) error
}

// cartRepo is the concrete implementation
type cartRepo struct {
    db Repository
}

// NewCartRepository creates a new CartRepository
func NewCartRepository(db Repository) CartRepository {
    return &cartRepo{db: db}
}

// CreateCart creates a new cart for a user
func (r *cartRepo) CreateCart(userID int64) (*domain.Cart, error) {
    query := "INSERT INTO Cart (UserID, TotalAmount, IsActive) VALUES (?, 0.00, true)"
    result, err := r.db.Exec(query, userID)
    if err != nil {
        return nil, err
    }

    cartID, err := result.LastInsertId()
    if err != nil {
        return nil, err
    }

    return &domain.Cart{
        ID:          cartID,
        UserID:      userID,
        TotalAmount: 0,
        Items:       make([]domain.CartItem, 0),
        IsActive:    true,
    }, nil
}

// GetCartByUserID retrieves a cart by user ID
func (r *cartRepo) GetCartByUserID(userID int64) (*domain.Cart, error) {
    var cart domain.Cart

    query := "SELECT ID, UserID, TotalAmount, IsActive FROM Cart WHERE UserID = ? AND IsActive = true"
    err := r.db.QueryRow(query, userID).Scan(&cart.ID, &cart.UserID, &cart.TotalAmount, &cart.IsActive)

    if err == sql.ErrNoRows {
        return nil, errors.New("cart not found")
    }
    if err != nil {
        return nil, err
    }

    // Get cart items
    items, err := r.GetCartItems(cart.ID)
    if err != nil {
        return nil, err
    }

    cart.Items = *convertCartItemsToSlice(items)
    return &cart, nil
}

// AddItemToCart adds a product to cart
func (r *cartRepo) AddItemToCart(cartID int64, productID int64, quantity int, price float64) (*domain.CartItem, error) {
    // Check if item already exists
    query := "SELECT ID, Quantity FROM CartItem WHERE CartID = ? AND ProductID = ?"
    var itemID int64
    var existingQty int

    err := r.db.QueryRow(query, cartID, productID).Scan(&itemID, &existingQty)

    if err == nil {
        // Item exists, update quantity
        return r.UpdateCartItem(cartID, productID, existingQty+quantity)
    }

    if err != sql.ErrNoRows {
        return nil, err
    }

    // Item doesn't exist, insert new
    insertQuery := "INSERT INTO CartItem (CartID, ProductID, Quantity, Price) VALUES (?, ?, ?, ?)"
    result, err := r.db.Exec(insertQuery, cartID, productID, quantity, price)
    if err != nil {
        return nil, err
    }

    itemID, err = result.LastInsertId()
    if err != nil {
        return nil, err
    }

    // Update cart total
    r.UpdateCartTotal(cartID)

    return &domain.CartItem{
        ID:        itemID,
        CartID:    cartID,
        ProductID: productID,
        Quantity:  quantity,
        Price:     price,
    }, nil
}

// RemoveItemFromCart removes a product from cart
func (r *cartRepo) RemoveItemFromCart(cartID int64, productID int64) error {
    query := "DELETE FROM CartItem WHERE CartID = ? AND ProductID = ?"
    result, err := r.db.Exec(query, cartID, productID)
    if err != nil {
        return err
    }

    rowsAffected, err := result.RowsAffected()
    if err != nil {
        return err
    }

    if rowsAffected == 0 {
        return errors.New("item not found in cart")
    }

    // Update cart total
    r.UpdateCartTotal(cartID)

    return nil
}

// UpdateCartItem updates the quantity of a cart item
func (r *cartRepo) UpdateCartItem(cartID int64, productID int64, quantity int) (*domain.CartItem, error) {
    if quantity <= 0 {
        return nil, errors.New("quantity must be greater than 0")
    }

    query := "UPDATE CartItem SET Quantity = ? WHERE CartID = ? AND ProductID = ?"
    result, err := r.db.Exec(query, quantity, cartID, productID)
    if err != nil {
        return nil, err
    }

    rowsAffected, err := result.RowsAffected()
    if err != nil {
        return nil, err
    }

    if rowsAffected == 0 {
        return nil, errors.New("item not found in cart")
    }

    // Get updated item
    var item domain.CartItem
    getQuery := "SELECT ID, CartID, ProductID, Quantity, Price FROM CartItem WHERE CartID = ? AND ProductID = ?"
    err = r.db.QueryRow(getQuery, cartID, productID).Scan(&item.ID, &item.CartID, &item.ProductID, &item.Quantity, &item.Price)
    if err != nil {
        return nil, err
    }

    // Update cart total
    r.UpdateCartTotal(cartID)

    return &item, nil
}

// GetCartItems retrieves all items in a cart
func (r *cartRepo) GetCartItems(cartID int64) ([]*domain.CartItem, error) {
    query := "SELECT ID, CartID, ProductID, Quantity, Price FROM CartItem WHERE CartID = ?"
    rows, err := r.db.Query(query, cartID)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var items []*domain.CartItem
    for rows.Next() {
        var item domain.CartItem
        err := rows.Scan(&item.ID, &item.CartID, &item.ProductID, &item.Quantity, &item.Price)
        if err != nil {
            return nil, err
        }
        items = append(items, &item)
    }

    if err := rows.Err(); err != nil {
        return nil, err
    }

    return items, nil
}

// ClearCart removes all items from a cart
func (r *cartRepo) ClearCart(cartID int64) error {
    query := "DELETE FROM CartItem WHERE CartID = ?"
    _, err := r.db.Exec(query, cartID)
    if err != nil {
        return err
    }

    // Reset cart total
    updateQuery := "UPDATE Cart SET TotalAmount = 0.00 WHERE ID = ?"
    _, err = r.db.Exec(updateQuery, cartID)
    return err
}

// DeleteCart deletes a cart
func (r *cartRepo) DeleteCart(cartID int64) error {
    // First delete all cart items
    query := "DELETE FROM CartItem WHERE CartID = ?"
    _, err := r.db.Exec(query, cartID)
    if err != nil {
        return err
    }

    // Then delete the cart
    deleteQuery := "DELETE FROM Cart WHERE ID = ?"
    result, err := r.db.Exec(deleteQuery, cartID)
    if err != nil {
        return err
    }

    rowsAffected, err := result.RowsAffected()
    if err != nil {
        return err
    }

    if rowsAffected == 0 {
        return errors.New("cart not found")
    }

    return nil
}

// UpdateCartTotal updates the total amount of a cart
func (r *cartRepo) UpdateCartTotal(cartID int64) error {
    query := "SELECT COALESCE(SUM(Quantity * Price), 0) FROM CartItem WHERE CartID = ?"
    var total float64

    err := r.db.QueryRow(query, cartID).Scan(&total)
    if err != nil {
        return err
    }

    updateQuery := "UPDATE Cart SET TotalAmount = ? WHERE ID = ?"
    _, err = r.db.Exec(updateQuery, total, cartID)
    return err
}

// Helper function to convert cart items
func convertCartItemsToSlice(items []*domain.CartItem) *[]domain.CartItem {
    result := make([]domain.CartItem, 0)
    for _, item := range items {
        result = append(result, *item)
    }
    return &result
}