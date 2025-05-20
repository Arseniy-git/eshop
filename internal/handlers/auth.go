package handlers

import (
	"eshop/internal/db"
	"eshop/internal/models"
	"eshop/internal/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ShowLoginPage(c *gin.Context) {
	userID, exists := c.Get("userID")
	isLoggedIn := exists && userID != nil
	c.HTML(http.StatusOK, "login.html", gin.H{"Title": "Login",
		"IsLoggedIn": isLoggedIn})
}

func ShowRegisterPage(c *gin.Context) {
	userID, exists := c.Get("userID")
	isLoggedIn := exists && userID != nil
	c.HTML(http.StatusOK, "register.html", gin.H{"Title": "Register",
		"IsLoggedIn": isLoggedIn})
}

// func Register(c *gin.Context) {
// 	var user models.User
// 	if err := c.BindJSON(&user); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
// 		return
// 	}

// 	hash, _ := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
// 	_, err := db.DB.Exec("INSERT INTO users (username, password_hash) VALUES (?, ?)", user.Name, hash)
// 	if err != nil {
// 		c.JSON(http.StatusConflict, gin.H{"error": "username already exists"})
// 		return
// 	}

// 	c.Status(http.StatusCreated)
// }

// func Login(c *gin.Context) {
// 	var req models.User
// 	if err := c.BindJSON(&req); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
// 		return
// 	}

// 	var storedHash string
// 	var userID int
// 	err := db.DB.QueryRow("SELECT id, password_hash FROM users WHERE username = ?", req.Name).
// 		Scan(&userID, &storedHash)
// 	if err != nil {
// 		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
// 		return
// 	}

// 	if bcrypt.CompareHashAndPassword([]byte(storedHash), []byte(req.Password)) != nil {
// 		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
// 		return
// 	}

// 	// Генерируем JWT
// 	token, err := utils.GenerateJWT(userID)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not generate token"})
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{"token": token})
// }

func Register(c *gin.Context) {
	name := c.PostForm("name")
	email := c.PostForm("email")
	password := c.PostForm("password")

	hashed, err := utils.HashPassword(password)
	if err != nil {
		c.String(http.StatusInternalServerError, "Error hashing password")
		return
	}

	_, err = db.DB.Exec("INSERT INTO users (name, email, password_hash) VALUES (?, ?, ?)", name, email, hashed)
	if err != nil {
		c.String(http.StatusInternalServerError, "User already exists or DB error")
		return
	}

	c.Redirect(http.StatusSeeOther, "/login")
}

func Login(c *gin.Context) {
	email := c.PostForm("email")
	password := c.PostForm("password")

	var user models.User
	row := db.DB.QueryRow("SELECT id, password_hash FROM users WHERE email = ?", email)
	if err := row.Scan(&user.ID, &user.PasswordHash); err != nil {
		c.String(http.StatusUnauthorized, "Invalid email")
		return
	}

	if !utils.CheckPasswordHash(password, user.PasswordHash) {
		c.String(http.StatusUnauthorized, "Invalid password")
		return
	}

	token, err := utils.GenerateJWT(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not generate token"})
		return
	}

	//c.SetCookie("auth", token, 3600*24, "/", "", false, true)
	c.SetCookie("Authorization", token, 3600, "/", "", false, true)
	c.Redirect(http.StatusSeeOther, "/products/create")
}

func Logout(c *gin.Context) {
	// Удаляем cookie с JWT
	c.SetCookie("Authorization", "", -1, "/", "", false, true)

	// Редирект на главную страницу
	c.Redirect(http.StatusSeeOther, "/")
}
