package application

import (
	"ecommerce-go/domain"
	infrastructure "ecommerce-go/infrastructure/mysql"
	"encoding/json"
	"net/http"
	"strconv"
)

type CreateProductRequest struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
}

type UpdateProductRequest struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
}

type ProductResponse struct {
	ID          int64   `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
}

func CreateProductHandler(repo infrastructure.ProductRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var req CreateProductRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request", http.StatusBadRequest)
			return
		}

		// Validate input
		if req.Name == "" || req.Price <= 0 {
			http.Error(w, "Invalid product data", http.StatusBadRequest)
			return
		}

		// Create product
		product := &domain.Product{
			Name:        req.Name,
			Description: req.Description,
			Price:       req.Price,
		}

		id, err := repo.Create(product)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		response := ProductResponse{
			ID:          id,
			Name:        product.Name,
			Description: product.Description,
			Price:       product.Price,
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(response)
	}
}

func UpdateProductHandler(repo infrastructure.ProductRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Get product ID from query parameter
		productID := r.URL.Query().Get("id")
		if productID == "" {
			http.Error(w, "Product ID required", http.StatusBadRequest)
			return
		}

		id, err := strconv.ParseInt(productID, 10, 64)
		if err != nil {
			http.Error(w, "Invalid product ID", http.StatusBadRequest)
			return
		}

		var req UpdateProductRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request", http.StatusBadRequest)
			return
		}

		// Validate input
		if req.Name == "" || req.Price <= 0 {
			http.Error(w, "Invalid product data", http.StatusBadRequest)
			return
		}

		// Update product
		product := &domain.Product{
			ID:          id,
			Name:        req.Name,
			Description: req.Description,
			Price:       req.Price,
		}

		err = repo.Update(product)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(product)
	}
}

// DeleteProductHandler - Delete a product
func DeleteProductHandler(repo infrastructure.ProductRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Get product ID from query parameter
		productID := r.URL.Query().Get("id")
		if productID == "" {
			http.Error(w, "Product ID required", http.StatusBadRequest)
			return
		}

		id, err := strconv.ParseInt(productID, 10, 64)
		if err != nil {
			http.Error(w, "Invalid product ID", http.StatusBadRequest)
			return
		}

		// Delete product
		err = repo.Delete(id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNoContent)
	}
}
