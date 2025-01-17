package api

import (
	"bot/api/handler"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	files "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	_ "bot/docs" // Import generated Swagger docs
)

// @title Password Management API
// @version 1.0
// @description This is an API for managing user passwords
// @BasePath /
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func NewGin(h *handler.HTTPHandler) *gin.Engine {
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"https://password-manager.eslab.uz/"}, // Frontend URL
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
	}))

	password := r.Group("/password")
	{
		password.POST("", h.CreatePassword)
		password.GET("/:userID", h.GetAllPasswordsByUserID)
		password.GET("", h.GetByName)
	}

	url := ginSwagger.URL("/api/swagger/doc.json")
	r.GET("/api/swagger/*any", ginSwagger.WrapHandler(files.Handler, url))

	return r
}
