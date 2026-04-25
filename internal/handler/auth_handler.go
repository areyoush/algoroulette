package handler

import (
	"net/http"
	"os"
	"strings"
	"time"
	"log"
	
	"github.com/areyoush/algoroulette/internal/model"
	"github.com/areyoush/algoroulette/internal/repository"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"	
)

type AuthHandler struct {
	repo *repository.UserRepository
}

func NewAuthHandler(repo *repository.UserRepository) *AuthHandler {
	return &AuthHandler{repo: repo}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var body struct {
		Email		string	`json:"email"`
		Password	string	`json:"password"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	if strings.Contains(body.Email, " ") {
    	c.JSON(http.StatusBadRequest, gin.H{"error": "email cannot contain spaces"})
     	return
	}

	if !strings.Contains(body.Email, "@") {
    	c.JSON(http.StatusBadRequest, gin.H{"error": "email must contain @"})
     	return
	}
	
	parts := strings.Split(body.Email, "@")
	if len(parts) != 2 || parts[0] == "" {
    	c.JSON(http.StatusBadRequest, gin.H{"error": "email must have a name before @"})
     	return
	}
	
	if !strings.Contains(parts[1], ".") || parts[1] == "." {
    	c.JSON(http.StatusBadRequest, gin.H{"error": "email must have a valid domain"})
     	return
	}
	
	if len(body.Password) < 8 || len(body.Password) > 72 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Password must be between 8 and 72 characters" })
		return
	}
	
	hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to hash password"})
		return
	}
	
	user := &model.User{
		Email:		body.Email,
		Password:	string(hash),
	}
	
	if err := h.repo.Create(user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "registration failed. please check your details or try logging in"})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"id": user.ID, "email": user.Email})

}


func (h *AuthHandler) Login(c *gin.Context) {
	var body struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// find user by email
	user, err := h.repo.GetByEmail(body.Email)
	if err != nil || user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	// compare password with hash 
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}


	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(7 * 24 * time.Hour).Unix(),
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate token"})
		return
	}
	
	isProd := os.Getenv("ENV") == "production"
	c.SetCookie("token", tokenString, 7*24*3600, "/", "", isProd, true)
	
	c.JSON(http.StatusOK, gin.H{"message": "logged in successfully"})
}

func (h *AuthHandler) Logout(c *gin.Context) {
	isProd := os.Getenv("ENV") == "production"
	
	tokenString, err := c.Cookie("token")
	
	if err == nil && tokenString != "" {
		expiry := time.Now().Add(7 * 24 * time.Hour)
		
		if err := h.repo.DenylistToken(tokenString, expiry); err != nil {
			log.Printf("Could not denylist token: %v", err)
		}
	}
	
	
	c.SetCookie("token", "", -1, "/", "", isProd, true)
	c.JSON(http.StatusOK, gin.H{"message": "logged out successfully"})
}