package middleware

import (
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"strings"
	"time"

	"github.com/Cognize-AI/server-cognize/config"
	"github.com/Cognize-AI/server-cognize/models"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

var _config config.Config

func init() {
	var err error
	_config, err = config.LoadConfig(".")
	if err != nil {
		panic(err)
	}
}

func RequireAuth(c *gin.Context) {
	tokenString, err := c.Cookie("Authorization")
	if err != nil {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}

		// Split the header to get the token part
		// Example: "Bearer eyJhbGci..." -> ["Bearer", "eyJhbGci..."]
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header must be in format 'Bearer <token>'"})
			c.Abort()
			return
		}
		tokenString = parts[1]
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Unauthorized",
			})
			c.Abort()
			return
		}
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, err
		}
		return []byte(_config.JwtSecret), nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if float64(time.Now().Unix()) > claims["exp"].(float64) {
			c.AbortWithStatus(http.StatusUnauthorized)
		}

		var user models.User
		config.DB.First(&user, claims["id"])

		if user.ID == 0 {
			c.AbortWithStatus(http.StatusUnauthorized)
		}

		c.Set("user", user)

		c.Next()
	} else {
		c.AbortWithStatus(http.StatusUnauthorized)
	}
	c.Next()
}

func RequireAPIKey(c *gin.Context) {
	apiKey := c.GetHeader("Cognize-API-Key")
	if apiKey == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "API key required in 'Cognize-API-Key' header"})
		c.Abort()
		return
	}

	hash := sha256.Sum256([]byte(apiKey))
	hashed := hex.EncodeToString(hash[:])

	var key models.Key
	if err := config.DB.Where("hash = ?", hashed).First(&key).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid API key"})
		c.Abort()
		return
	}

	// (Optional) check if key is active / expired
	//if !key.Active {
	//	c.JSON(http.StatusUnauthorized, gin.H{"error": "API key is inactive"})
	//	c.Abort()
	//	return
	//}

	// Attach key info (or associated user) to context
	c.Set("apiKey", key)
	c.Next()
}
