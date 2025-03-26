package service

import (
	"fmt"

	"github.com/kikudesuyo/arumonogohan-app/api/entity"
	"github.com/kikudesuyo/arumonogohan-app/api/repository"
	"github.com/line/line-bot-sdk-go/linebot"
)

type LineReplyParams struct {
	UserID  string
	Message string
	Events  []*linebot.Event
	Session *entity.ChatHistory
	Store   *repository.ChatSessionStore
	LineBot *LineBotClient
}

func ReplyMenuSelect(params *LineReplyParams) error {
	var replyMessage string
	switch {
	case isMenuSelected(params.Message):
		params.Session.State = "ingredient_input"
		params.Store.Save(params.UserID, params.Message, params.Session.State)
		replyMessage = fmt.Sprintf("「%s」ですね✨️ 使う食材を教えて下さい!!", params.Message)
	default:
		replyMessage = "メニューから料理するジャンルを選択ください🍽️"
	}

	// LineBotがnilかどうか確認
	if params.LineBot == nil {
		return fmt.Errorf("LineBot is nil")
	}
	err := params.LineBot.ReplyMessage(params.Events, replyMessage)
	if err != nil {
		return fmt.Errorf("error sending reply message: %v", err)
	}
	return nil
}

func ReplyIngredientInput(params *LineReplyParams) error {
	var replyMessage string

	switch {
	case isMenuSelected(params.Message):
		replyMessage = fmt.Sprintf("「%s」ですね✨️ 使う食材を教えて下さい!!", params.Message)
	default:
		m, err := SuggestRecipe(params.Message)
		if err != nil {
			return fmt.Errorf("error handling suggest recipe: %v", err)
		}
		replyMessage = m
		params.Session.State = "menu_select"
		params.Store.Save(params.UserID, params.Message, params.Session.State)
	}
	if err := params.LineBot.ReplyMessage(params.Events, replyMessage); err != nil {
		return fmt.Errorf("error sending reply message: %v", err)
	}
	return nil
}

func isMenuSelected(message string) bool {
	menus := map[string]struct{}{
		"時短メニュー⏱️":  {},
		"家庭の味🥢":     {},
		"さっぱりヘルシー🥗": {},
		"ガッツリメニュー🍖": {},
	}
	_, exists := menus[message]
	return exists
}

func (l *LineBotClient) ReplyMessage(events []*linebot.Event, resMessage string) error {
	for _, event := range events {
		if event.Type == linebot.EventTypeMessage {
			_, err := l.Bot.ReplyMessage(
				event.ReplyToken,
				linebot.NewTextMessage(resMessage),
			).Do()
			if err != nil {
				return fmt.Errorf("error sending reply message: %v", err)
			}
		}
	}
	return nil
}
