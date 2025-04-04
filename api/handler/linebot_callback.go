package handler

import (
	"fmt"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/kikudesuyo/arumonogohan-app/api/entity"
	"github.com/kikudesuyo/arumonogohan-app/api/repository"
	"github.com/kikudesuyo/arumonogohan-app/api/service"
)

var store = &repository.ChatSessionStore{}

func HandleLinebotCallback(c *gin.Context) {
	userAgent := c.GetHeader("User-Agent")
	// LINE Botクライアントの場合
	if strings.Contains(userAgent, "LineBotWebhook") {
		lineBot, err := service.NewLineBotClient(store)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		events, err := lineBot.Bot.ParseRequest(c.Request)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		message, err := lineBot.GetMessage(events)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		userID, err := lineBot.GetUserID(events)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		session, err := store.Get(userID)
		if err != nil {
			fmt.Println("session not found. creating new session")
			session = &entity.ChatHistory{Messages: []string{}, State: "menu_select", Timestamp: time.Now()}
			store.Save(userID, message, session.State)
		}

		replayParams := &service.LineReplyParams{
			LineBot: lineBot,
			UserID:  userID,
			Message: message,
			Session: session,
			Store:   store,
			Events:  events,
		}

		switch {
		case session.State == "menu_select":
			err = service.ReplyMenuSelect(replayParams)
			if err != nil {
				fmt.Println(err.Error())
				return
			}
		case session.State == "ingredient_input":
			err = service.ReplyIngredientInput(replayParams)
			if err != nil {
				fmt.Println(err.Error())
				return
			}
		}
		return
	}
}
