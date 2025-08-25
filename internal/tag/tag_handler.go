package tag

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

func (h *Handler) CreateTag(c *gin.Context) {
	currentUser, valid := util.GetCurrentUser(c)
	if !valid {
		return
	}

	var req CreateTagReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	res, err := h.Service.CreateTag(c, req, currentUser)
	if err != nil {
		logger.Logger.Error("Error creating tag: ", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": res})
}

func (h *Handler) AddTag(c *gin.Context) {
	currentUser, valid := util.GetCurrentUser(c)
	if !valid {
		return
	}

	var req AddTagReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	err := h.Service.AddTag(c, req, currentUser)
	if err != nil {
		logger.Logger.Error("Error adding tag: ", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": "ok"})
}

func (h *Handler) GetAllTags(c *gin.Context) {
	currentUser, valid := util.GetCurrentUser(c)
	if !valid {
		return
	}

	tags, err := h.Service.GetAllTags(c, currentUser)
	if err != nil {
		logger.Logger.Error("Error getting tags: ", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": tags})
}

func (h *Handler) DeleteTag(c *gin.Context) {
	currentUser, valid := util.GetCurrentUser(c)
	if !valid {
		return
	}

	var req DeleteTagReq
	if err := c.ShouldBindUri(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.Service.DeleteTag(c, req, currentUser)

	if err != nil {
		logger.Logger.Error("Error deleting tag: ", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": "ok"})
}

func (h *Handler) EditTag(c *gin.Context) {
	currentUser, valid := util.GetCurrentUser(c)
	if !valid {
		return
	}

	var req EditTagReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	res, err := h.Service.EditTag(c, req, currentUser)
	if err != nil {
		logger.Logger.Error("Error editing tag: ", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": res})
}

func (h *Handler) RemoveTagAssociation(c *gin.Context) {
	currentUser, valid := util.GetCurrentUser(c)
	if !valid {
		return
	}

	var req RemoveTagReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	err := h.Service.RemoveTagAssociation(c, req, currentUser)
	if err != nil {
		logger.Logger.Error("Error removing tag association: ", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": "ok"})
}
