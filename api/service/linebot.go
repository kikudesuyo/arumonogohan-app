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

func (c *LineBotClient) GetMessage(events []*linebot.Event) (string, error) {
	for _, event := range events {
		if event.Type == linebot.EventTypeMessage {
			switch message := event.Message.(type) {
			case *linebot.TextMessage:
				return message.Text, nil
			}
		}
	}
	return "", fmt.Errorf("no message found in events")
}

func (c *LineBotClient) GetUserID(events []*linebot.Event) (string, error) {
	for _, event := range events {
		if event.Type == linebot.EventTypeMessage {
			userID := event.Source.UserID
			return userID, nil
		}
	}
	return "", fmt.Errorf("no user ID found in events")
}
