package handler

import (
	"fmt"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/kikudesuyo/arumonogohan-app/api/entity"
	"github.com/kikudesuyo/arumonogohan-app/api/infrastructure"
	"github.com/kikudesuyo/arumonogohan-app/api/repository"
	"github.com/kikudesuyo/arumonogohan-app/api/usecase"
)

var store = &repository.ChatSessionStore{}

func HandleLinebotCallback(c *gin.Context) {
	userAgent := c.GetHeader("User-Agent")
	if !strings.Contains(userAgent, "LineBotWebhook") {
		return
	}

	linebot, err := infrastructure.NewLineBotClient()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	events, err := infrastructure.ParseLinebotRequest(c.Request, linebot)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	lineUserMsg, err := usecase.GetLineUserMsg(events)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	userID := lineUserMsg.UserID
	msg := lineUserMsg.Msg

	chatSession, found := store.Get(userID)
	if !found {
		fmt.Println("session not found. creating new session")
		chatSession = &repository.ChatSession{
			SessionID:    userID,
			MenuCategory: "",
			State:        entity.StateMenuCategorySelect,
			Timestamp:    time.Now(),
		}
		store.Save(*chatSession)
		err := infrastructure.ReplyMsgToLine(linebot, events, "メニューから料理するジャンルを選択ください🍽️")
		if err != nil {
			fmt.Println(err.Error())
			return
		}
	}
	// メニューカテゴリが選択された場合
	if entity.IsMenuCategorySelected(msg) {
		err := usecase.ProcessSelectMenuCategory(linebot, events, msg, chatSession)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
	}
	//メニューカテゴリが選択フェーズの場合
	switch chatSession.State {
	case entity.StateMenuCategorySelect:
		err := usecase.ProcessSelectMenuCategory(linebot, events, msg, chatSession)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
	// メニュー入力フェーズの場合
	case entity.StateIngredientInput:
		err := usecase.ProcessInputIngredient(linebot, events, msg, chatSession)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
	}
}
