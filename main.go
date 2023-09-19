package main

import (
	"github/be/common"
	"github/be/controller"
	"github/be/database"
	docs "github/be/docs"
	"github/be/middleware"
	"github/be/model"
	"log"

	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func init() {
	common.Env()
}

//	@title			Simple BE
//	@version		1.0
//	@description	Simple example of BE

func main() {
	database.Connect()
	if err := database.Database.AutoMigrate(&model.User{}, &model.Entry{}); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	router := gin.Default()

	v1 := router.Group("/api/v1")
	{
		publicRoutes := v1.Group("/auth")
		publicRoutes.POST("/register", controller.Register)
		publicRoutes.POST("/login", controller.Login)

		protectedRoutes := v1.Group("/")
		protectedRoutes.Use(middleware.JWTAuthMiddleware())
		protectedRoutes.POST("/entry", controller.AddEntry)
		protectedRoutes.GET("/entry", controller.GetAllEntries)
	}

	docs.SwaggerInfo.BasePath = "/api/v1"
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
