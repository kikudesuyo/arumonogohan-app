package main

import (
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/kikudesuyo/arumonogohan-app/api/entity"
	"github.com/kikudesuyo/arumonogohan-app/api/handler"
)

var store = &entity.LineSessionStore{}

func main() {
	if err := godotenv.Load(".env"); err != nil {
		fmt.Println("No .env file found")
	}

	engine := gin.Default()
	engine.POST("/callback", postCallback)
	engine.Run(fmt.Sprintf(":%s", os.Getenv("PORT")))
}

func postCallback(c *gin.Context) {
	lineBot, err := entity.NewLineBotClient(store)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	events, err := lineBot.Bot.ParseRequest(c.Request)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	lineBot.SaveMessageToStore(events)
	userID, err := lineBot.GetUserID(events)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	d, err := store.Get(userID)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	messages := d.Messages
	fmt.Println(messages)
	var menus = []string{
		"時短メニュー⏱️",
		"家庭の味🥢",
		"さっぱりヘルシー🥗",
		"ガッツリメニュー🍖",
	}
	if len(messages) == 1 {
		for _, menu := range menus {
			if menu == messages[len(messages)-1] {
				replyMessage := fmt.Sprintf("「%s」ですね✨️ 使う食材を教えて下さい!!", menu)
				err = lineBot.ReplyMessage(events, replyMessage)
				if err != nil {
					fmt.Println(err.Error())
					return
				}
			}
		}
		replyMessage := "メニューを選んでください🍽️"
		err = lineBot.ReplyMessage(events, replyMessage)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		return
	}
	var replyMessage string
	for _, menu := range menus {
		if menu == messages[len(messages)-1] {
			replyMessage := fmt.Sprintf("「%s」ですね✨️ 使う食材を教えて下さい!!", menu)
			err = lineBot.ReplyMessage(events, replyMessage)
			if err != nil {
				fmt.Println(err.Error())
				return
			}
		}
	}
	for _, menu := range menus {
		if menu == messages[len(messages)-2] {
			replyMessage, err = handler.HandleSuggestRecipe(messages[len(messages)-1])
			if err != nil {
				fmt.Println(err.Error())
				return
			}
		}
	}
	if replyMessage == "" {
		replyMessage = "メニューを選んでください🍽️"
	}
	err = lineBot.ReplyMessage(events, replyMessage)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
}
