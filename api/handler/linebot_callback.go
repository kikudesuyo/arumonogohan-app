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
	if !strings.Contains(userAgent, "LineBotWebhook") {
		return
	}
	// LINE Botクライアントの場合
	lineBot, err := service.NewLineBotClient()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	events, err := lineBot.Bot.ParseRequest(c.Request)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	lineMsg, err := lineBot.GetLineMsg(events)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	userID := lineMsg.UserID
	msg := lineMsg.Msg
	session, err := store.Get(userID)
	if err != nil {
		fmt.Println("session not found. creating new session")
		state := entity.StateMenuCategorySelect
		store.Save(userID, state)
		session = &entity.ChatHistory{State: state, Timestamp: time.Now()}
	}

	var replyMsg string
	switch session.State {
	case "menu_category_select":
		if entity.IsMenuCategorySelected(msg) {
			// メニュー選択時の処理
			newState := entity.StateMenuCategorySelect
			store.Save(userID, newState)
			session.State = newState
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
			newState := entity.StateMenuCategorySelect
			store.Save(userID, newState)
			session.State = newState
			replyMsg = fmt.Sprintf("「%s」ですね✨️ 使う食材を教えて下さい!!", msg)
		} else {
			menuCategory := session.Msg
			recipeInput := service.RecipeInput{
				MenuCategory: menuCategory,
				Ingredients:  msg,
			}
			m, err := service.SuggestRecipe(recipeInput)
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
