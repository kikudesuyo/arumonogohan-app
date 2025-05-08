package main

import (
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/kikudesuyo/arumonogohan-app/api/handler"
)

func main() {
	if err := godotenv.Load(".env"); err != nil {
		fmt.Println("No .env file found", err)
	}
	engine := gin.Default()
	engine.POST("/callback", handler.HandleLinebotCallback)
	engine.Run(fmt.Sprintf(":%s", os.Getenv("PORT")))
}
