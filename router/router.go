package router

import (
	"net/http"

	"github.com/Cognize-AI/client-cognize/internal/card"
	"github.com/Cognize-AI/client-cognize/internal/list"
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
	listHandler *list.Handler,
	cardHandler *card.Handler,
) {
	r = gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOriginFunc: func(origin string) bool {
			return origin == "http://localhost:3000" || origin == "https://client-cognize.vercel.app"
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

	listRouter := r.Group("/list")
	{
		listRouter.GET("/create-default", middleware.RequireAuth, listHandler.CreateDefaultLists)
		listRouter.GET("/all", middleware.RequireAuth, listHandler.GetLists)
	}

	cardRouter := r.Group("/card")
	{
		cardRouter.POST("/create", middleware.RequireAuth, cardHandler.CreateCard)
		cardRouter.POST("/move", middleware.RequireAuth, cardHandler.MoveCard)
		cardRouter.DELETE("/:id", middleware.RequireAuth, cardHandler.DeleteCard)
		cardRouter.PUT("/:id", middleware.RequireAuth, cardHandler.UpdateCard)
	}
}

func Start(addr string) error {
	return r.Run(addr)
}
