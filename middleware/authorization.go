package middleware

import (
	"fmt"
	"os"
	"time"

	"github.com/JohnnyOhms/projectx/config"
	"github.com/JohnnyOhms/projectx/entity"
	_ "github.com/JohnnyOhms/projectx/model"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func Authorization(ctx *gin.Context) {
	// get the cookie from req body
	tokenString, err := ctx.Cookie("Authorization")
	if err != nil {
		ctx.AbortWithStatusJSON(500, gin.H{"error": "Cookie not found"})
		return
	}

	// decode/validate the JWT token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("SECRET")), nil
	})

	if err != nil || !token.Valid {
		ctx.AbortWithStatusJSON(500, gin.H{"error": "Failed to authorize"})
		return
	}

	// Extract claims
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		ctx.AbortWithStatusJSON(500, gin.H{"error": "Invalid token claims"})
		return
	}

	// Check for token expiration
	expirationTime := time.Unix(int64(claims["exp"].(float64)), 0)
	if time.Now().After(expirationTime) {
		ctx.AbortWithStatusJSON(500, gin.H{"error": "Token has expired"})
		return
	}

	// Find the user with token sub
	var user entity.User
	if err := config.DB.First(&user, "user_id = ?", claims["sub"]).Error; err != nil {
		ctx.AbortWithStatusJSON(500, gin.H{"error": "No user found during authorization"})
		return
	}

	// Attach user to the context
	ctx.Set("user", user)

	// Continue with the next middleware/handler
	ctx.Next()
}
