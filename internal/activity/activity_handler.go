package activity

import (
	"net/http"

	"github.com/Cognize-AI/client-cognize/logger"
	"github.com/Cognize-AI/client-cognize/util"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	Service
}

func NewHandler(s Service) *Handler {
	return &Handler{s}
}

func (h *Handler) CreateActivity(c *gin.Context) {
	currentUser, valid := util.GetCurrentUser(c)
	if !valid {
		return
	}

	var req CreateActivityReq
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Logger.Error("Failed to parse json")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := h.Service.CreateActivity(c, req, currentUser)
	if err != nil {
		logger.Logger.Error("Failed to create activity")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": res})
}
