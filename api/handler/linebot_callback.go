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
	lineBot := lineMsgCtx.Bot
	events := lineMsgCtx.Events
	userID := lineMsgCtx.UserMsg.UserID
	msg := lineMsgCtx.UserMsg.Msg

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
	}

	switch chatSession.State {
	case entity.StateMenuCategorySelect:
		if !entity.IsMenuCategorySelected(msg) {
			replyMsg := "メニューから料理するジャンルを選択ください🍽️"
			err := usecase.ReplyMsgToLine(lineBot, events, replyMsg)
			if err != nil {
				fmt.Println(err.Error())
				return
			}
			return
		}
		// メニューカテゴリ選択時の処理

		chatSession.MenuCategory = msg
		chatSession.State = entity.StateIngredientInput
		chatSession.Timestamp = time.Now()

		store.Save(*chatSession)
		replyMsg := fmt.Sprintf("「%s」ですね✨️ 使う食材を教えて下さい!!", msg)
		err := usecase.ReplyMsgToLine(lineBot, events, replyMsg)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
	case entity.StateIngredientInput:
		// メニューカテゴリ再選択の場合
		if entity.IsMenuCategorySelected(msg) {
			chatSession.MenuCategory = msg
			chatSession.State = entity.StateIngredientInput
			chatSession.Timestamp = time.Now()
			store.Save(*chatSession)

			replyMsg := fmt.Sprintf("「%s」ですね✨️ 使う食材を教えて下さい!!", msg)
			err := usecase.ReplyMsgToLine(lineBot, events, replyMsg)
			if err != nil {
				fmt.Println(err.Error())
				return
			}
			return
		} else {
			recipeInput := usecase.RecipeInput{
				MenuCategory: chatSession.MenuCategory,
				Ingredients:  msg,
			}
			replyMsg, err := usecase.SuggestRecipe(recipeInput)
			if err != nil {
				fmt.Println(err.Error())
				return
			}
			chatSession.State = entity.StateMenuCategorySelect
			chatSession.MenuCategory = ""
			chatSession.Timestamp = time.Now()
			store.Save(*chatSession)

			err = usecase.ReplyMsgToLine(lineBot, events, replyMsg)
			if err != nil {
				fmt.Println(err.Error())
				return
			}
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
		Bot:     bot,
		Events:  events,
		UserMsg: msg,
	}, nil
}
