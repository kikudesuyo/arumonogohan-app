package handler

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/kikudesuyo/arumonogohan-app/api/entity"
	"github.com/kikudesuyo/arumonogohan-app/api/repository"
	"github.com/kikudesuyo/arumonogohan-app/api/usecase"
)

var store = &repository.ChatSessionStore{}

func HandleLinebotCallback(c *gin.Context) {
	userAgent := c.GetHeader("User-Agent")
	if !strings.Contains(userAgent, "LineBotWebhook") {
		return
	}

	lineMsgCtx, err := parseLineRequest(c.Request)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	lineMsg := lineMsgCtx.Msg
	lineBot := lineMsgCtx.Bot
	events := lineMsgCtx.Events
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
	case entity.StateMenuCategorySelect:
		if entity.IsMenuCategorySelected(msg) {
			// メニュー選択時の処理
			newState := entity.StateMenuCategorySelect
			store.Save(userID, newState)
			session.State = newState
			replyMsg = fmt.Sprintf("「%s」ですね✨️ 使う食材を教えて下さい!!", msg)
		} else {
			replyMsg = "メニューから料理するジャンルを選択ください🍽️"
		}
		err := usecase.ReplyMsgToLine(lineBot, events, replyMsg)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
	case entity.StateIngredientInput:
		// メニュー再選択の場合
		if entity.IsMenuCategorySelected(msg) {
			newState := entity.StateMenuCategorySelect
			store.Save(userID, newState)
			session.State = newState
			replyMsg = fmt.Sprintf("「%s」ですね✨️ 使う食材を教えて下さい!!", msg)
		} else {
			menuCategory := session.Msg
			recipeInput := usecase.RecipeInput{
				MenuCategory: menuCategory,
				Ingredients:  msg,
			}
			m, err := usecase.SuggestRecipe(recipeInput)
			if err != nil {
				fmt.Println(err.Error())
				return
			}
			replyMsg = m
		}
		err := usecase.ReplyMsgToLine(lineBot, events, replyMsg)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
	}
}

func parseLineRequest(r *http.Request) (*usecase.LineMsgContext, error) {
	bot, err := usecase.NewLineBotClient()
	if err != nil {
		return nil, fmt.Errorf("failed to create LINE bot client: %v", err)
	}
	events, err := bot.ParseRequest(r)
	if err != nil {
		return nil, fmt.Errorf("failed to parse request: %v", err)
	}
	msg, err := usecase.GetLineMsg(events)
	if err != nil {
		return nil, fmt.Errorf("failed to get line message: %v", err)
	}
	return &usecase.LineMsgContext{
		Bot:    bot,
		Events: events,
		Msg:    msg,
	}, nil
}
