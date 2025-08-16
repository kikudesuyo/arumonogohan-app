package usecase

import (
	"fmt"
	"time"

	"github.com/kikudesuyo/arumonogohan-app/api/entity"
	"github.com/kikudesuyo/arumonogohan-app/api/infrastructure"
	"github.com/kikudesuyo/arumonogohan-app/api/repository"
	"github.com/line/line-bot-sdk-go/linebot"
)

func ProcessSelectMenuCategory(bot *linebot.Client, events []*linebot.Event, lineUserMsg infrastructure.LineUserMsg, store *repository.ChatSessionStore) error {
	replyMsg := fmt.Sprintf("「%s」ですね✨️ \n次に使う食材を教えて下さい!👨‍🍳", lineUserMsg.Msg)
	err := infrastructure.ReplyMsgToLine(bot, events, replyMsg)
	if err != nil {
		return err
	}

	//状態更新
	chatSession := &repository.ChatSession{
		SessionID:    lineUserMsg.UserID,
		MenuCategory: lineUserMsg.Msg,
		State:        entity.StateIngredientInput,
		Timestamp:    time.Now(),
	}
	store.UpsertChatSession(*chatSession)
	return nil

}
