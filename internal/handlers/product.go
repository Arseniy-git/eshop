package handlers

import (
	"database/sql"
	"eshop/internal/db"
	"eshop/internal/forms"
	"eshop/internal/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// func HomePage(c *gin.Context) {
// 	rows, err := db.DB.Query(`
// 		SELECT id, title, description, price, quantity
// 		FROM products
// 	`)
// 	if err != nil {
// 		c.HTML(http.StatusInternalServerError, "error.html", gin.H{"error": "Database error"})
// 		return
// 	}
// 	defer rows.Close()

// 	var products []models.Product
// 	for rows.Next() {
// 		var p models.Product
// 		if err := rows.Scan(&p.ID, &p.Title, &p.Description, &p.Price, &p.Quantity); err != nil {
// 			c.HTML(http.StatusInternalServerError, "error.html", gin.H{"error": "Parsing error"})
// 			return
// 		}
// 		products = append(products, p)
// 	}

// 	c.HTML(http.StatusOK, "home.html", gin.H{
// 		"Products": products,
// 	})
// }

// GET /products

// func ShowHomePage(c *gin.Context) {
// 	c.HTML(http.StatusOK, "home.html", gin.H{
// 		"Products": []models.Product{},
// 	})
// }

func ShowCreateProductPage(c *gin.Context) {
	userID, exists := c.Get("userID")
	isLoggedIn := exists && userID != nil
	c.HTML(http.StatusOK, "create_product.html", gin.H{
		"IsLoggedIn": isLoggedIn})
}

func ListProducts(c *gin.Context) {
	userID, exists := c.Get("userID")
	isLoggedIn := exists && userID != nil
	rows, err := db.DB.Query("SELECT id, title, description, price, quantity, user_id FROM products")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
		return
	}
	defer rows.Close()

	var products []models.Product
	for rows.Next() {
		var p models.Product
		if err := rows.Scan(&p.ID, &p.Title, &p.Description, &p.Price, &p.Quantity, &p.UserID); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "scan error"})
			return
		}
		products = append(products, p)
	}

	//c.JSON(http.StatusOK, products)
	c.HTML(http.StatusOK, "home.html", gin.H{
		"Products":   products,
		"IsLoggedIn": isLoggedIn,
	})
}

// POST /products
// func CreateProduct(c *gin.Context) {
// 	userID := c.GetInt("userID")
// 	var p models.Product
// 	if err := c.BindJSON(&p); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid data"})
// 		return
// 	}

// 	p.UserID = userID

// 	res, err := db.DB.Exec("INSERT INTO products (title, description, price, quantity, user_id) VALUES (?, ?, ?, ?, ?)",
// 		p.Title, p.Description, p.Price, p.Quantity, p.UserID)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "insert failed"})
// 		return
// 	}

// 	id, _ := res.LastInsertId()
// 	p.ID = int(id)
// 	c.JSON(http.StatusCreated, p)
// }

func CreateProduct(c *gin.Context) {

	userIDInterface, exists := c.Get("userID")

	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	userID := userIDInterface.(int)

	title := c.PostForm("title")
	description := c.PostForm("description")
	priceStr := c.PostForm("price")
	quantityStr := c.PostForm("quantity")

	// Парсим числа
	price, err := strconv.ParseFloat(priceStr, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid price"})
		return
	}
	quantity, err := strconv.Atoi(quantityStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid quantity"})
		return
	}

	// Добавляем продукт
	_, err = db.DB.Exec(`
		INSERT INTO products (title, description, price, quantity, user_id)
		VALUES (?, ?, ?, ?, ?)`,
		title, description, price, quantity, userID,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not insert product"})
		return
	}

	c.Redirect(http.StatusSeeOther, "/")
}

// PUT /products/:id
// func UpdateProduct(c *gin.Context) {
// 	id := c.Param("id")
// 	var p models.Product
// 	if err := c.BindJSON(&p); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid data"})
// 		return
// 	}

// 	_, err := db.DB.Exec("UPDATE products SET title=?, description=?, price=?, quantity=? WHERE id=?",
// 		p.Title, p.Description, p.Price, p.Quantity, id)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "update failed"})
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{"message": "updated"})
// }

// DELETE /products/:id
// func DeleteProduct(c *gin.Context) {
// 	id := c.Param("id")
// 	_, err := db.DB.Exec("DELETE FROM products WHERE id=?", id)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "delete failed"})
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{"message": "deleted"})
// }

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

func UpdateProduct(c *gin.Context) {
	userID := c.GetInt("userID")
	id := c.Param("id")

	// Проверим, принадлежит ли товар текущему пользователю
	var ownerID int
	err := db.DB.QueryRow("SELECT user_id FROM products WHERE id = ?", id).Scan(&ownerID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}
	if ownerID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "You can only update your own products"})
		return
	}

	// Прочитаем обновлённые данные
	var p models.Product
	if err := c.ShouldBindJSON(&p); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	_, err = db.DB.Exec(`
		UPDATE products
		SET title = ?, description = ?, price = ?, quantity = ?
		WHERE id = ?
	`, p.Title, p.Description, p.Price, p.Quantity, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Update failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Product updated"})
}

// func DeleteProduct(c *gin.Context) {
// 	userID := c.GetInt("userID")
// 	id := c.Param("id")

// 	// Проверим владельца
// 	var ownerID int
// 	err := db.DB.QueryRow("SELECT user_id FROM products WHERE id = ?", id).Scan(&ownerID)
// 	if err != nil {
// 		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
// 		return
// 	}
// 	if ownerID != userID {
// 		c.JSON(http.StatusForbidden, gin.H{"error": "You can only delete your own products"})
// 		return
// 	}

// 	_, err = db.DB.Exec("DELETE FROM products WHERE id = ?", id)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Delete failed"})
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{"message": "Product deleted"})
// }

func ListMyProducts(c *gin.Context) {
	userID, exists := c.Get("userID")
	isLoggedIn := exists && userID != nil

	rows, err := db.DB.Query(`
		SELECT id, title, description, price, quantity, user_id
		FROM products
		WHERE user_id = ?
	`, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}
	defer rows.Close()

	products := []models.Product{}
	for rows.Next() {
		var p models.Product
		if err := rows.Scan(&p.ID, &p.Title, &p.Description, &p.Price, &p.Quantity, &p.UserID); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse product"})
			return
		}
		products = append(products, p)
	}

	// c.JSON(http.StatusOK, products)
	c.HTML(http.StatusOK, "my_products.html", gin.H{
		"Products":   products,
		"IsLoggedIn": isLoggedIn,
	})
}

// ////////////////
func ShowEditProductPage(c *gin.Context) {
	id := c.Param("id")
	var product models.Product
	err := db.DB.QueryRow("SELECT id, title, description, price, quantity FROM products WHERE id = ?", id).
		Scan(&product.ID, &product.Title, &product.Description, &product.Price, &product.Quantity)

	if err != nil {
		c.String(http.StatusNotFound, "Product not found")
		return
	}

	userID, _ := c.Get("userID")
	isLoggedIn := userID != nil

	c.HTML(http.StatusOK, "edit_product.html", gin.H{
		"Product":    product,
		"IsLoggedIn": isLoggedIn,
	})
}

func EditProduct(c *gin.Context) {
	id := c.Param("id")

	var form forms.ProductForm
	if err := c.ShouldBind(&form); err != nil {
		c.String(http.StatusBadRequest, "Invalid form data")
		return
	}

	_, err := db.DB.Exec(`
		UPDATE products
		SET title = ?, description = ?, price = ?, quantity = ?
		WHERE id = ?`,
		form.Title, form.Description, form.Price, form.Quantity, id)

	if err != nil {
		c.String(http.StatusInternalServerError, "Failed to update product")
		return
	}

	c.Redirect(http.StatusSeeOther, "/my-products")
}

func DeleteProduct(c *gin.Context) {
	id := c.Param("id")
	_, err := db.DB.Exec("DELETE FROM products WHERE id = ?", id)

	if err != nil {
		c.String(http.StatusInternalServerError, "Failed to delete product")
		return
	}

	c.Redirect(http.StatusFound, "/my-products")
}
