package service

import (
	"context"

	"github.com/kikudesuyo/arumonogohan-app/api/entity"
	"github.com/kikudesuyo/arumonogohan-app/api/infrastructure"
	"github.com/kikudesuyo/arumonogohan-app/api/repository"
	"github.com/kikudesuyo/arumonogohan-app/api/xerror"
	"github.com/line/line-bot-sdk-go/linebot"
)

// ProcessLinebotMessage は、ユーザー入力に応じた対話フローを処理する。
func ProcessLinebotMessage(
	ctx context.Context,
	bot *linebot.Client,
	events []*linebot.Event,
	lineUserMsg infrastructure.LineUserMsg,
	store repository.ChatSessionRepository,
) error {
	chatSession, err := store.Find(lineUserMsg.UserID)
	if err != nil {
		return err
	}
	if chatSession == nil {
		chatSession = entity.NewChatSession(lineUserMsg.UserID)
		if entity.IsMenuCategorySelected(lineUserMsg.Msg) {
			return ProcessSelectMenuCategory(bot, events, lineUserMsg, store)
		}
		if err := store.Save(chatSession); err != nil {
			return err
		}
		return infrastructure.ReplyMsgToLine(bot, events, "あるものごはんへようこそ！\nメニューから料理するジャンルを選択してください🍽️")
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
		return xerror.InvalidChatState(chatSession.State)
	}
}
