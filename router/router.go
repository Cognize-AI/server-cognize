package router

import (
	"net/http"

	"github.com/Cognize-AI/client-cognize/internal/activity"
	"github.com/Cognize-AI/client-cognize/internal/card"
	"github.com/Cognize-AI/client-cognize/internal/field"
	"github.com/Cognize-AI/client-cognize/internal/keys"
	"github.com/Cognize-AI/client-cognize/internal/list"
	"github.com/Cognize-AI/client-cognize/internal/oauth"
	"github.com/Cognize-AI/client-cognize/internal/tag"
	"github.com/Cognize-AI/client-cognize/internal/user"
	"github.com/Cognize-AI/client-cognize/logger"
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
	tagHandler *tag.Handler,
	keyHandler *keys.Handler,
	fieldHandler *field.Handler,
	activityHandler *activity.Handler,
) {
	r = gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOriginFunc: func(origin string) bool {
			return origin == "http://localhost:3000" || origin == "https://client-cognize.vercel.app" || origin == "https://cognize.live" || origin == "https://www.cognize.live"
		},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	r.GET("/", func(c *gin.Context) {
		logger.Logger.Info("lol")
		c.String(http.StatusOK, "Welcome To Cognize!")
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
		cardRouter.GET("/:id", middleware.RequireAuth, cardHandler.GetCardById)
		cardRouter.PUT("/details/:id", middleware.RequireAuth, cardHandler.UpdateCardByID)
	}

	tagRouter := r.Group("/tag")
	{
		tagRouter.POST("/create", middleware.RequireAuth, tagHandler.CreateTag)
		tagRouter.POST("/add-to-card", middleware.RequireAuth, tagHandler.AddTag)
		tagRouter.GET("/", middleware.RequireAuth, tagHandler.GetAllTags)
		tagRouter.DELETE("/:id", middleware.RequireAuth, tagHandler.DeleteTag)
		tagRouter.PUT("/", middleware.RequireAuth, tagHandler.EditTag)
		tagRouter.POST("/remove-from-card", middleware.RequireAuth, tagHandler.RemoveTagAssociation)
	}

	keyRouter := r.Group("/key")
	{
		keyRouter.GET("/api", middleware.RequireAuth, keyHandler.CreateAPI)
		keyRouter.GET("/", middleware.RequireAuth, keyHandler.GetAPI)
	}

	APIRouter := r.Group("/api")
	{
		APIRouter.POST("/bulk-prospect", middleware.RequireAPIKey, cardHandler.BulkCreate)
	}

	fieldRouter := r.Group("/field")
	{
		fieldRouter.POST("/field-definitions", middleware.RequireAuth, fieldHandler.CreateField)
		fieldRouter.POST("/field-value", middleware.RequireAuth, fieldHandler.InsertFieldVal)
		fieldRouter.GET("/", middleware.RequireAuth, fieldHandler.GetFields)
	}

	activityRouter := r.Group("/activity")
	{
		activityRouter.POST("/create", middleware.RequireAuth, activityHandler.CreateActivity)
		activityRouter.DELETE("/:id", middleware.RequireAuth, activityHandler.DeleteActivity)
		activityRouter.PUT("/:id", middleware.RequireAuth, activityHandler.UpdateActivity)
	}
}

func Start(addr string) error {
	return r.Run(addr)
}
