package router

import (
	"net/http"

	"github.com/Cognize-AI/client-cognize/internal/oauth"
	"github.com/Cognize-AI/client-cognize/internal/user"
	"github.com/Cognize-AI/client-cognize/middleware"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

var r *gin.Engine

func InitRouter(
	userHandler *user.Handler,
	oauthHandler *oauth.Handler,
) {
	r = gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOriginFunc: func(origin string) bool {
			return origin == "http://localhost:5173"
		},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Welcome to cognize")
	})

	userRouter := r.Group("/user")
	{
		userRouter.GET("/me", middleware.RequireAuth, userHandler.Me)
	}

	oAuthRouter := r.Group("/oauth")
	{
		oAuthRouter.GET("/google/redirect-uri", oauthHandler.GetRedirectURL)
		oAuthRouter.GET("/google/callback", oauthHandler.HandleGoogleCallback)
	}
}

func Start(addr string) error {
	return r.Run(addr)
}
