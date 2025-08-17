package oauth

import (
	"log"
	"net/http"

	"github.com/Cognize-AI/client-cognize/config"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	Service
}

func NewHandler(s Service) *Handler {
	return &Handler{s}
}

func (h *Handler) GetRedirectURL(c *gin.Context) {
	res, err := h.Service.GetRedirectURL(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": res})
}

func (h *Handler) HandleGoogleCallback(c *gin.Context) {
	code := c.Query("code")
	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No code in request"})
		return
	}

	res, err := h.Service.HandleGoogleCallback(c, &HandleGoogleCallbackReq{Code: code})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	_config, err := config.LoadConfig(".")
	if err != nil {
		log.Fatal(err)
	}

	if _config.Environment != "dev" {
		c.SetSameSite(http.SameSiteNoneMode)
		c.SetCookie("Authorization", res.Token, 3600*24*30, "/", "", true, true)
	} else {
		c.SetSameSite(http.SameSiteLaxMode)
		c.SetCookie("Authorization", res.Token, 3600*24*30, "", "", false, true)
	}
	c.JSON(http.StatusOK, gin.H{"data": res})
}
