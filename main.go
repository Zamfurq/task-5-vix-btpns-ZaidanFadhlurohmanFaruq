package main

import (
	"fmt"
	"github/golang-api/controller"
	"github/golang-api/database"
	"github/golang-api/middlewares"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	DbHost := os.Getenv("DB_HOST")
	DbUser := os.Getenv("DB_USER")
	DbPassword := os.Getenv("DB_PASSWORD")
	DbName := os.Getenv("DB_NAME")
	DbPort := os.Getenv("DB_PORT")

	database.Connect(fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", DbUser, DbPassword, DbHost, DbPort, DbName))
	database.Migrate()

	router := initRouter()
	router.Run(":8080")
}

func initRouter() *gin.Engine {
	router := gin.Default()
	user := router.Group("/user")
	{
		user.POST("/login", controller.GenerateToken)
		user.POST("/register", controller.RegisterUser)
		user.PUT("/:id", controller.UpdateUser).Use(middlewares.Auth())
		user.DELETE("/:id", controller.DeleteUser).Use(middlewares.Auth())
		secured := user.Group("/admin").Use(middlewares.Auth())
		{
			secured.GET("/photo", controller.GetPhotos)
			secured.POST("/photo", controller.PostPhoto)
			secured.PUT("/photo/:id", controller.UpdatePhoto)
			secured.DELETE("/photo/:id", controller.DeletePhoto)
		}
	}
	return router
}
