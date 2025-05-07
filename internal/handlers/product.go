package handlers

import (
	"database/sql"
	"eshop/internal/db"
	"eshop/internal/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GET /products
func ListProducts(c *gin.Context) {
	rows, err := db.DB.Query("SELECT id, title, description, price, quantity FROM products")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
		return
	}
	defer rows.Close()

	var products []models.Product
	for rows.Next() {
		var p models.Product
		if err := rows.Scan(&p.ID, &p.Title, &p.Description, &p.Price, &p.Quantity); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "scan error"})
			return
		}
		products = append(products, p)
	}

	c.JSON(http.StatusOK, products)
}

// POST /products
func CreateProduct(c *gin.Context) {
	userID := c.GetInt("userID")
	var p models.Product
	if err := c.BindJSON(&p); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid data"})
		return
	}

	p.UserID = userID

	res, err := db.DB.Exec("INSERT INTO products (title, description, price, quantity, user_id) VALUES (?, ?, ?, ?, ?)",
		p.Title, p.Description, p.Price, p.Quantity, p.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "insert failed"})
		return
	}

	id, _ := res.LastInsertId()
	p.ID = int(id)
	c.JSON(http.StatusCreated, p)
}

// PUT /products/:id
func UpdateProduct(c *gin.Context) {
	id := c.Param("id")
	var p models.Product
	if err := c.BindJSON(&p); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid data"})
		return
	}

	_, err := db.DB.Exec("UPDATE products SET title=?, description=?, price=?, quantity=? WHERE id=?",
		p.Title, p.Description, p.Price, p.Quantity, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "update failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "updated"})
}

// DELETE /products/:id
func DeleteProduct(c *gin.Context) {
	id := c.Param("id")
	_, err := db.DB.Exec("DELETE FROM products WHERE id=?", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "delete failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "deleted"})
}

func GetProductByID(c *gin.Context) {
	id := c.Param("id")

	var product models.Product

	err := db.DB.QueryRow(`
		SELECT id, title, description, price, quantity
		FROM products
		WHERE id = ?
	`, id).Scan(
		&product.ID,
		&product.Title,
		&product.Description,
		&product.Price,
		&product.Quantity,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		}
		return
	}

	c.JSON(http.StatusOK, product)
}
