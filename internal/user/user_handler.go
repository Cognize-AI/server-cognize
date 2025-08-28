package user

import (
	"net/http"
	"strconv"

	"github.com/Cognize-AI/server-cognize/logger"
	"github.com/Cognize-AI/server-cognize/models"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Handler struct {
	Service
}

func NewHandler(s Service) *Handler {
	return &Handler{s}
}

func (h *Handler) Me(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		logger.Logger.Warn("Failed to get user from context : UNAUTHORIZED")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	currentUser := user.(models.User)

	res, err := h.Service.Me(c, currentUser)
	if err != nil {
		logger.Logger.Error("user not found: ", zap.String("id", strconv.Itoa(int(currentUser.ID))))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": res})
}
