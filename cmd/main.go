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

	// Роуты без авторизации
	r.POST("/register", handlers.Register)
	r.POST("/login", handlers.Login)

	r.GET("/products", handlers.ListProducts) // публичный
	r.GET("/products/:id", handlers.GetProductByID)

	auth := r.Group("/")
	auth.Use(middleware.JWTAuth())
	{
		auth.POST("/products", handlers.CreateProduct)
		auth.PUT("/products/:id", handlers.UpdateProduct)
		auth.DELETE("/products/:id", handlers.DeleteProduct)
	}

	log.Println("Server started on :8080")
	r.Run(":8080")
}
