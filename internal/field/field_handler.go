package field

import (
	"net/http"

	"github.com/Cognize-AI/client-cognize/logger"
	"github.com/Cognize-AI/client-cognize/util"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Handler struct {
	Service
}

func NewHandler(s Service) *Handler {
	return &Handler{s}
}

func (h *Handler) CreateField(c *gin.Context) {
	currentUser, valid := util.GetCurrentUser(c)
	if !valid {
		return
	}

	var req CreateFieldReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		logger.Logger.Error("CreateField ShouldBindJSON", zap.Error(err))
		return
	}

	res, err := h.Service.CreateField(c, req, currentUser)
	if err != nil {
		logger.Logger.Error("CreateField", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": res})
}

func (h *Handler) InsertFieldVal(c *gin.Context) {
	currentUser, valid := util.GetCurrentUser(c)
	if !valid {
		return
	}

	var req InsertFieldValReq
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Logger.Error("InsertFieldVal ShouldBindJSON", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := h.Service.InsertFieldVal(c, req, currentUser)
	if err != nil {
		logger.Logger.Error("InsertFieldVal", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": res})
}

func (h *Handler) GetFields(c *gin.Context) {
	currentUser, valid := util.GetCurrentUser(c)
	if !valid {
		return
	}

	res, err := h.Service.GetFields(c, currentUser)
	if err != nil {
		logger.Logger.Error("GetFields", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": res})
}

func (h *Handler) UpdateFieldDefinition(c *gin.Context) {
	currentUser, valid := util.GetCurrentUser(c)
	if !valid {
		return
	}

	var req UpdateFieldDef
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Logger.Error("UpdateFieldDefinition ShouldBindJSON", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.Service.UpdateFieldDefinition(c, req, currentUser)
	if err != nil {
		logger.Logger.Error("UpdateFieldDefinition", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": "success"})
}
