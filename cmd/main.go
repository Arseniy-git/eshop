package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"eshop/config"
	"eshop/internal/db"
	"eshop/internal/handlers"
	"eshop/internal/middleware"
)

func main() {

	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}
	// Инициализируем подключение к БД
	db.Init(config.GetDSN())

	r := gin.Default()
	r.LoadHTMLGlob("internal/templates/*.html")
	r.Use(middleware.JWTFromCookie())
	r.Use(middleware.AddLoginStatus())

	// Роуты без авторизации
	//r.GET("/", handlers.ShowHomePage)
	r.GET("/login", handlers.ShowLoginPage)
	r.GET("/register", handlers.ShowRegisterPage)

	r.POST("/register", handlers.Register)
	r.POST("/login", handlers.Login)

	r.GET("/", handlers.ListProducts) // публичный
	r.GET("/products/:id", handlers.GetProductByID)

	r.GET("/products/create", handlers.ShowCreateProductPage)
	r.GET("/logout", handlers.Logout)

	auth := r.Group("/")
	auth.Use(middleware.JWTAuth())
	{
		auth.POST("/products", handlers.CreateProduct)
		auth.PUT("/products/:id", handlers.UpdateProduct)
		auth.DELETE("/products/:id", handlers.DeleteProduct)
		auth.GET("/my-products", handlers.ListMyProducts)
		r.POST("/products/create", handlers.CreateProduct)
		r.GET("/products/edit/:id", handlers.ShowEditProductPage)
		r.POST("/products/edit/:id", handlers.EditProduct)
		r.POST("/products/delete/:id", handlers.DeleteProduct)

	}

	log.Println("Server started on :8080")
	r.Run(":8080")
}
