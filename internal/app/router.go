package app

import (
	"wxcloudrun-golang/internal/handler"
	"wxcloudrun-golang/internal/middleware"

	"github.com/gin-gonic/gin"
)

func NewRouter() *gin.Engine {
	router := gin.Default()

	router.GET("/", handler.Index)
	router.GET("/healthz", func(c *gin.Context) {
		c.JSON(200, gin.H{"ok": true})
	})

	api := router.Group("/api")
	api.Use(middleware.RequireOpenID())
	{
		api.POST("/login", handler.Login)
		api.GET("/families", handler.ListFamilies)
		api.POST("/families", handler.CreateFamily)
		api.POST("/families/join", handler.JoinFamily)
		api.GET("/families/:familyID", handler.GetFamilyDetail)
		api.POST("/families/:familyID/invite-code", handler.CreateFamilyInviteCode)
		api.GET("/babies", handler.ListBabies)
		api.POST("/babies", handler.CreateBaby)
		api.GET("/actions", handler.ListActions)
		api.POST("/actions", handler.CreateAction)
		api.POST("/actions/batch", handler.BatchCreateActions)
		api.GET("/summary/today", handler.TodaySummary)
		api.POST("/ai/parse", handler.ParseRecord)
		api.POST("/ai/chat", handler.Chat)
	}

	return router
}
