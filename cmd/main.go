package main

import (
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/kikudesuyo/arumonogohan-app/api/entity"
	"github.com/kikudesuyo/arumonogohan-app/api/handlers"
)

func main() {
	if err := godotenv.Load(".env"); err != nil {
		fmt.Println("No .env file found")
	}

	engine := gin.Default()
	engine.POST("/callback", postCallback)
	engine.Run(fmt.Sprintf(":%s", os.Getenv("PORT")))
}

func postCallback(c *gin.Context) {
	lineBot, err := entity.NewLineBotClient()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	events, err := lineBot.Bot.ParseRequest(c.Request)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	clientMessage, err := lineBot.GetMessage(events)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	//clientMessageの内容に応じて該当するhandlerを呼び出す
	var replyMessage string
	if clientMessage == "" {
		return
	} else if clientMessage != "" { // 現在は文字列があればレシピを提案
		replyMessage, err = handlers.HandleSuggestRecipe(clientMessage)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
	}
	err = lineBot.ReplyMessage(events, replyMessage)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
}
