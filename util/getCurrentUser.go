package util

import (
	"net/http"

	"github.com/Cognize-AI/client-cognize/logger"
	"github.com/Cognize-AI/client-cognize/models"
	"github.com/gin-gonic/gin"
)

func GetCurrentUser(c *gin.Context) (models.User, bool) {
	user, exists := c.Get("user")
	if !exists {
		logger.Logger.Warn("Failed to get user from context : UNAUTHORIZED")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return models.User{}, false
	}
	return user.(models.User), true
}
