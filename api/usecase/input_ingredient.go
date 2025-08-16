package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/kikudesuyo/arumonogohan-app/api/entity"
	"github.com/kikudesuyo/arumonogohan-app/api/infrastructure"
	"github.com/kikudesuyo/arumonogohan-app/api/repository"
	"github.com/line/line-bot-sdk-go/linebot"
)

func ProcessInputIngredient(bot *linebot.Client, events []*linebot.Event, lineUserMsg infrastructure.LineUserMsg, chatSession *repository.ChatSession, store *repository.ChatSessionStore) error {
	// メニューカテゴリ再選択の場合
	if entity.IsMenuCategorySelected(lineUserMsg.Msg) {
		chatSession.MenuCategory = lineUserMsg.Msg
		chatSession.State = entity.StateIngredientInput
		chatSession.Timestamp = time.Now()
		store.UpsertChatSession(*chatSession)

		replyMsg := fmt.Sprintf("「%s」ですね✨️ 使う食材を教えて下さい!!", lineUserMsg.Msg)
		err := infrastructure.ReplyMsgToLine(bot, events, replyMsg)
		if err != nil {
			return err
		}
		return nil
	}

	recipeInput := entity.RecipeInputReq{
		MenuCategory: chatSession.MenuCategory,
		Ingredients:  lineUserMsg.Msg,
	}
	ctx := context.Background()
	recipe, err := SuggestRecipe(ctx, recipeInput)
	if err != nil {
		return err
	}

	// Recipe構造体をLINE用の文字列にフォーマットする
	replyMsg := formatRecipeForLine(recipe)
	err = infrastructure.ReplyMsgToLine(bot, events, replyMsg)
	if err != nil {
		return err
	}
	store.InsertInitChatSession(lineUserMsg.UserID)
	return nil
}
