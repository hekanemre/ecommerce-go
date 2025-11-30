package main

import (
	applicationCart "ecommerce-go/application/cart"
	applicationProduct "ecommerce-go/application/product"
	applicationUser "ecommerce-go/application/user"
	"ecommerce-go/config"
	infrastructure "ecommerce-go/infrastructure/mysql"
	"log"
	"net/http"
)

func main() {

	conf, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	dbRepo, err := infrastructure.NewMySQLRepository(
		conf.DB.User,
		conf.DB.Password,
		conf.DB.Host,
		conf.DB.Port,
		conf.DB.Database,
	)
	if err != nil {
		log.Fatal(err)
	}
	defer dbRepo.Close()

	productRepo := infrastructure.NewProductRepository(dbRepo)
	userRepo := infrastructure.NewUserRepository(dbRepo)
	cartRepo := infrastructure.NewCartRepository(dbRepo)

	// Product routes
	http.HandleFunc("/product", applicationProduct.GetProductHandler(productRepo))
	http.HandleFunc("/products", applicationProduct.GetAllProductsHandler(productRepo))
	http.HandleFunc("/product/create", applicationProduct.CreateProductHandler(productRepo))
	http.HandleFunc("/product/update", applicationProduct.UpdateProductHandler(productRepo))
	http.HandleFunc("/product/delete", applicationProduct.DeleteProductHandler(productRepo))

	// User routes
	http.HandleFunc("/auth/login", func(w http.ResponseWriter, r *http.Request) {
		applicationUser.HandleUserLogin(w, r, userRepo)
	})
	http.HandleFunc("/auth/signup", func(w http.ResponseWriter, r *http.Request) {
		applicationUser.HandleUserSignUp(w, r, userRepo)
	})

	// Cart routes
	http.HandleFunc("/cart/create", applicationCart.CreateCartHandler(cartRepo))
	http.HandleFunc("/cart/get", applicationCart.GetCartHandler(cartRepo))
	http.HandleFunc("/cart/add", applicationCart.AddToCartHandler(cartRepo))
	http.HandleFunc("/cart/remove", applicationCart.RemoveFromCartHandler(cartRepo))
	http.HandleFunc("/cart/update", applicationCart.UpdateCartItemHandler(cartRepo))
	http.HandleFunc("/cart/clear", applicationCart.ClearCartHandler(cartRepo))

	log.Println("Server running at http://localhost:9000")
	log.Fatal(http.ListenAndServe(":9000", nil))
}
