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
		lineEvent, err := lineBot.GetLineEvent(events)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		userID := lineEvent.UserID
		msg := lineEvent.Msg
		session, err := store.Get(userID)
		if err != nil {
			fmt.Println("session not found. creating new session")
			state := &entity.MenuCategorySelect{Category: msg}
			store.Save(userID, state)
			session = &entity.ChatHistory{StateData: state, Timestamp: time.Now()}
		}

		var replyMsg string
		switch session.StateData.GetState() {
		case "menu_category_select":
			if entity.IsMenuCategorySelected(msg) {
				// メニュー選択時の処理
				newState := &entity.IngredientInput{Msg: msg}
				store.Save(userID, newState)
				session.StateData = newState
				replyMsg = fmt.Sprintf("「%s」ですね✨️ 使う食材を教えて下さい!!", msg)
			} else {
				replyMsg = "メニューから料理するジャンルを選択ください🍽️"
			}
			err := lineBot.ReplyMsg(events, replyMsg)
			if err != nil {
				fmt.Println(err.Error())
				return
			}
		case "ingredient_input":
			// メニュー再選択の場合
			if entity.IsMenuCategorySelected(msg) {
				newState := &entity.IngredientInput{Msg: msg}
				store.Save(userID, newState)
				session.StateData = newState
				replyMsg = fmt.Sprintf("「%s」ですね✨️ 使う食材を教えて下さい!!", msg)
			} else {
				menuCategory := session.StateData.GetMsg()
				m, err := service.SuggestRecipe(menuCategory, msg)
				if err != nil {
					fmt.Println(err.Error())
					return
				}
				replyMsg = m
			}
			err := lineBot.ReplyMsg(events, replyMsg)
			if err != nil {
				fmt.Println(err.Error())
				return
			}
		}
	}
}
