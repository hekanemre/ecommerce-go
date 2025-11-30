package application

import (
	"ecommerce-go/domain"
	infrastructure "ecommerce-go/infrastructure/mysql"
	"encoding/json"
	"net/http"
	"strconv"
)

type AddToCartRequest struct {
	UserID    int64   `json:"user_id"`
	ProductID int64   `json:"product_id"`
	Quantity  int     `json:"quantity"`
	Price     float64 `json:"price"`
}

type RemoveFromCartRequest struct {
	UserID    int64 `json:"user_id"`
	ProductID int64 `json:"product_id"`
}

type UpdateCartRequest struct {
	UserID    int64 `json:"user_id"`
	ProductID int64 `json:"product_id"`
	Quantity  int   `json:"quantity"`
}

type CartItemResponse struct {
	ID        int64   `json:"id"`
	ProductID int64   `json:"product_id"`
	Quantity  int     `json:"quantity"`
	Price     float64 `json:"price"`
	Subtotal  float64 `json:"subtotal"`
}

type CartResponse struct {
	CartID      int64              `json:"cart_id"`
	UserID      int64              `json:"user_id"`
	Items       []CartItemResponse `json:"items"`
	TotalAmount float64            `json:"total_amount"`
	TotalItems  int                `json:"total_items"`
}

// CreateCartHandler - Create a new cart for user
func CreateCartHandler(repo infrastructure.CartRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		userID := r.URL.Query().Get("user_id")
		if userID == "" {
			http.Error(w, "User ID required", http.StatusBadRequest)
			return
		}

		uid, err := strconv.ParseInt(userID, 10, 64)
		if err != nil {
			http.Error(w, "Invalid user ID", http.StatusBadRequest)
			return
		}

		// Create cart
		cart, err := repo.CreateCart(uid)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		response := CartResponse{
			CartID:      cart.ID,
			UserID:      cart.UserID,
			Items:       make([]CartItemResponse, 0),
			TotalAmount: 0,
			TotalItems:  0,
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(response)
	}
}

// GetCartHandler - Get user's cart
func GetCartHandler(repo infrastructure.CartRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		userID := r.URL.Query().Get("user_id")
		if userID == "" {
			http.Error(w, "User ID required", http.StatusBadRequest)
			return
		}

		uid, err := strconv.ParseInt(userID, 10, 64)
		if err != nil {
			http.Error(w, "Invalid user ID", http.StatusBadRequest)
			return
		}

		// Get cart by user ID
		cart, err := repo.GetCartByUserID(uid)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		// Convert to response
		response := buildCartResponse(cart)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}

// AddToCartHandler - Add product to cart
func AddToCartHandler(repo infrastructure.CartRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var req AddToCartRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request", http.StatusBadRequest)
			return
		}

		// Validate input
		if req.UserID == 0 || req.ProductID == 0 || req.Quantity <= 0 || req.Price <= 0 {
			http.Error(w, "Invalid product data", http.StatusBadRequest)
			return
		}

		// Get user's cart
		cart, err := repo.GetCartByUserID(req.UserID)
		if err != nil {
			http.Error(w, "Cart not found", http.StatusNotFound)
			return
		}

		// Add item to cart
		_, err = repo.AddItemToCart(cart.ID, req.ProductID, req.Quantity, req.Price)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Get updated cart
		updatedCart, err := repo.GetCartByUserID(req.UserID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		response := buildCartResponse(updatedCart)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}
}

// RemoveFromCartHandler - Remove product from cart
func RemoveFromCartHandler(repo infrastructure.CartRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var req RemoveFromCartRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request", http.StatusBadRequest)
			return
		}

		// Validate input
		if req.UserID == 0 || req.ProductID == 0 {
			http.Error(w, "User ID and Product ID required", http.StatusBadRequest)
			return
		}

		// Get user's cart
		cart, err := repo.GetCartByUserID(req.UserID)
		if err != nil {
			http.Error(w, "Cart not found", http.StatusNotFound)
			return
		}

		// Remove item from cart
		err = repo.RemoveItemFromCart(cart.ID, req.ProductID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Get updated cart
		updatedCart, err := repo.GetCartByUserID(req.UserID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		response := buildCartResponse(updatedCart)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}

// UpdateCartItemHandler - Update product quantity in cart
func UpdateCartItemHandler(repo infrastructure.CartRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var req UpdateCartRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request", http.StatusBadRequest)
			return
		}

		// Validate input
		if req.UserID == 0 || req.ProductID == 0 || req.Quantity <= 0 {
			http.Error(w, "Invalid request data", http.StatusBadRequest)
			return
		}

		// Get user's cart
		cart, err := repo.GetCartByUserID(req.UserID)
		if err != nil {
			http.Error(w, "Cart not found", http.StatusNotFound)
			return
		}

		// Update cart item
		_, err = repo.UpdateCartItem(cart.ID, req.ProductID, req.Quantity)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Get updated cart
		updatedCart, err := repo.GetCartByUserID(req.UserID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		response := buildCartResponse(updatedCart)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}

// ClearCartHandler - Clear all items from cart
func ClearCartHandler(repo infrastructure.CartRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		userID := r.URL.Query().Get("user_id")
		if userID == "" {
			http.Error(w, "User ID required", http.StatusBadRequest)
			return
		}

		uid, err := strconv.ParseInt(userID, 10, 64)
		if err != nil {
			http.Error(w, "Invalid user ID", http.StatusBadRequest)
			return
		}

		// Get user's cart
		cart, err := repo.GetCartByUserID(uid)
		if err != nil {
			http.Error(w, "Cart not found", http.StatusNotFound)
			return
		}

		// Clear cart
		err = repo.ClearCart(cart.ID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		response := CartResponse{
			CartID:      cart.ID,
			UserID:      uid,
			Items:       make([]CartItemResponse, 0),
			TotalAmount: 0,
			TotalItems:  0,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}

// Helper function to build cart response
func buildCartResponse(cart *domain.Cart) CartResponse {
	items := make([]CartItemResponse, 0)
	totalItems := 0

	for _, item := range cart.Items {
		items = append(items, CartItemResponse{
			ID:        item.ID,
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			Price:     item.Price,
			Subtotal:  float64(item.Quantity) * item.Price,
		})
		totalItems += item.Quantity
	}

	return CartResponse{
		CartID:      cart.ID,
		UserID:      cart.UserID,
		Items:       items,
		TotalAmount: cart.TotalAmount,
		TotalItems:  totalItems,
	}
}
