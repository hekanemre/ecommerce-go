package main

import (
	application "ecommerce-go/application/product"
	infrastructure "ecommerce-go/infrastructure/mysql"
	"log"
	"net/http"
)

func main() {

	dbRepo, err := infrastructure.NewMySQLRepository("appuser", "appuser123", "localhost", "3306", "ECommercial")
	if err != nil {
		log.Fatal(err)
	}
	defer dbRepo.Close()

	productRepo := infrastructure.NewProductRepository(dbRepo)

	// Register HTTP handler
	http.HandleFunc("/product", application.GetProductHandler(productRepo))
	http.HandleFunc("/products", application.GetAllProductsHandler(productRepo))

	log.Println("Server running at http://localhost:9000")
	log.Fatal(http.ListenAndServe(":9000", nil))
}
