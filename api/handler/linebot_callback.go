package handler

import (
	"fmt"
	"strings"

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
	lineUserMsg, err := infrastructure.GetLineUserMsg(events)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	chatSession, found := store.GetChatSession(lineUserMsg.UserID)

	if !found {
		fmt.Println("session not found. creating new session")
		store.InsertInitChatSession(lineUserMsg.UserID)
		err := infrastructure.ReplyMsgToLine(linebot, events, "メニューから料理するジャンルを選択ください🍽️")
		if err != nil {
			fmt.Println(err.Error())
			return
		}
	}
	// メニューカテゴリが選択された場合
	if entity.IsMenuCategorySelected(lineUserMsg.Msg) {
		err := usecase.ProcessSelectMenuCategory(linebot, events, lineUserMsg, store)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
	}
	//メニューカテゴリが選択フェーズの場合
	switch chatSession.State {
	case entity.StateMenuCategorySelect:
		err := usecase.ProcessSelectMenuCategory(linebot, events, lineUserMsg, store)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
	// メニュー入力フェーズの場合
	case entity.StateIngredientInput:
		err := usecase.ProcessInputIngredient(linebot, events, lineUserMsg, chatSession, store)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
	}
}
