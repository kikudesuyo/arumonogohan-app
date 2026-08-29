package usecase

import (
	"context"
	"fmt"

	"github.com/kikudesuyo/arumonogohan-app/api/entity"
	"github.com/kikudesuyo/arumonogohan-app/api/infrastructure"
	"github.com/kikudesuyo/arumonogohan-app/api/repository"
	"github.com/line/line-bot-sdk-go/linebot"
)

// ProcessLinebotMessage は、ユーザー入力に応じた対話フローを処理する。
func ProcessLinebotMessage(
	ctx context.Context,
	bot *linebot.Client,
	events []*linebot.Event,
	lineUserMsg infrastructure.LineUserMsg,
	store *repository.ChatSessionStore,
) error {
	chatSession, found := store.GetChatSession(lineUserMsg.UserID)
	if !found {
		store.InsertInitChatSession(lineUserMsg.UserID)
		if entity.IsMenuCategorySelected(lineUserMsg.Msg) {
			return ProcessSelectMenuCategory(bot, events, lineUserMsg, store)
		}
		return infrastructure.ReplyMsgToLine(bot, events, "メニューから料理するジャンルを選択ください🍽️")
	}

	if entity.IsMenuCategorySelected(lineUserMsg.Msg) {
		return ProcessSelectMenuCategory(bot, events, lineUserMsg, store)
	}

	switch chatSession.State {
	case entity.StateMenuCategorySelect:
		return ProcessSelectMenuCategory(bot, events, lineUserMsg, store)
	case entity.StateIngredientInput:
		return ProcessInputIngredient(ctx, bot, events, lineUserMsg, chatSession, store)
	default:
		return fmt.Errorf("unknown chat session state: %s", chatSession.State)
	}
}
