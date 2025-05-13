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
	r := gin.Default()
	r.POST("/callback", handler.HandleLinebotCallback)
	r.Run(fmt.Sprintf(":%s", os.Getenv("PORT")))
}
