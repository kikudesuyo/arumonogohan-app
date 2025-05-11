package service

import (
	"fmt"
	"os"

	"github.com/kikudesuyo/arumonogohan-app/api/repository"
	"github.com/line/line-bot-sdk-go/linebot"
)

type LineBotClient struct {
	Bot *linebot.Client
}

type LineUserMsg struct {
	UserID string
	Msg    string
}

func NewLineBotClient(store *repository.ChatSessionStore) (*LineBotClient, error) {
	bot, err := linebot.New(
		os.Getenv("LINE_BOT_CHANNEL_SECRET"),
		os.Getenv("LINE_BOT_CHANNEL_TOKEN"),
	)
	if err != nil {
		return nil, fmt.Errorf("error creating LINE bot client: %v", err)
	}
	return &LineBotClient{Bot: bot}, nil
}

func (c *LineBotClient) GetLineEvent(events []*linebot.Event) (*LineUserMsg, error) {
	for _, event := range events {
		if event.Type == linebot.EventTypeMessage {
			// メッセージがテキスト型の場合
			if msg, ok := event.Message.(*linebot.TextMessage); ok {
				return &LineUserMsg{
					UserID: event.Source.UserID,
					Msg:    msg.Text,
				}, nil
			}
		}
	}
	return nil, fmt.Errorf("no message found in events")
}
