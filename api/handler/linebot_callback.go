package handler

import (
	"context"
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

// formatRecipeForLine は、Recipe構造体をLINEメッセージ用の整形済み文字列に変換します。
func formatRecipeForLine(recipe entity.Recipe) string {
	// レシピが生成できなかった場合やエラーの場合は、サマリーメッセージのみを返す
	if recipe.Title == "提案できません" || recipe.Title == "無効な入力です" || recipe.Title == "エラー" {
		return recipe.Summary
	}

	var builder strings.Builder
	builder.WriteString(fmt.Sprintf("今日のレシピは「%s」で決まり！\n\n", recipe.Title))

	if len(recipe.Ingredients) > 0 {
		builder.WriteString("【材料】\n")
		for _, ingredient := range recipe.Ingredients {
			builder.WriteString(fmt.Sprintf("- %s\n", ingredient))
		}
		builder.WriteString("\n")
	}

	if len(recipe.Instructions) > 0 {
		builder.WriteString("【作り方】\n")
		for i, instruction := range recipe.Instructions {
			builder.WriteString(fmt.Sprintf("%d. %s\n", i+1, instruction))
		}
		builder.WriteString("\n")
	}

	if recipe.Summary != "" {
		builder.WriteString(fmt.Sprintf("【ポイント】\n%s\n", recipe.Summary))
	}

	return builder.String()
}

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
			ctx := context.Background()
			recipe, err := usecase.SuggestRecipe(ctx, recipeInput)
			if err != nil {
				fmt.Println(err.Error())
				return
			}

			// Recipe構造体をLINE用の文字列にフォーマットする
			replyMsg := formatRecipeForLine(recipe)

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
	},
	nil
}