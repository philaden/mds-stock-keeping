package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	controller "github.com/philaden/mds-stock-keeping/controllers"
	"github.com/philaden/mds-stock-keeping/infrastructure"
)

// @title mds-stock-keeping docs
// @version 1.0
// @description This is the api documentation for my solution.
func main() {

	dbConfig, err := infrastructure.LoadConfiguration(".")

	if err != nil {
		fmt.Println(err)
		panic("failed to load application configuration settings")
	}

	port := fmt.Sprintf(":%d", dbConfig.AppPort)

	if port == ":" {
		port = ":8000"
	}

	infrastructure.SetUpDatabaseServices(dbConfig)

	router := gin.Default()

	router.Use(corsMiddleware())

	RegisterController(router)

	if err := router.Run(port); err != nil {
		fmt.Print(err)

	}
}

func RegisterController(router *gin.Engine) {
	controller.SetupContollerRoutes(router)
}

func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.Abort()
			return
		}
		c.Next()
	}
}
