package application

import (
	infrastructure "ecommerce-go/infrastructure/mysql"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
)

func GetAllProductsHandler(productRepo infrastructure.ProductRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		product, err := productRepo.GetAll()
		if err != nil {
			http.Error(w, "Failed to fetch product", http.StatusInternalServerError)
			return
		}

		if product == nil {
			http.Error(w, "Product not found", http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(product)
	}
}

func GetProductHandler(productRepo infrastructure.ProductRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		idParam := r.URL.Query().Get("id")
		if strings.TrimSpace(idParam) == "" {
			http.Error(w, "Missing 'id' query parameter", http.StatusBadRequest)
			return
		}

		id, err := strconv.ParseInt(idParam, 10, 64)
		if err != nil {
			http.Error(w, "Invalid 'id' parameter", http.StatusBadRequest)
			return
		}

		product, err := productRepo.GetByID(id)
		if err != nil {
			http.Error(w, "Failed to fetch product", http.StatusInternalServerError)
			return
		}

		if product == nil {
			http.Error(w, "Product not found", http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(product)
	}
}
