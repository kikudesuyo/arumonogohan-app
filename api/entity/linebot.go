package entity

import (
	"fmt"
	"os"

	"github.com/line/line-bot-sdk-go/linebot"
)

type LineBotClient struct {
	Bot   *linebot.Client
	Store *LineSessionStore
}

func NewLineBotClient(store *LineSessionStore) (*LineBotClient, error) {
	bot, err := linebot.New(
		os.Getenv("LINE_BOT_CHANNEL_SECRET"),
		os.Getenv("LINE_BOT_CHANNEL_TOKEN"),
	)
	if err != nil {
		return nil, fmt.Errorf("error creating LINE bot client: %v", err)
	}
	return &LineBotClient{Bot: bot, Store: store}, nil
}

func (c *LineBotClient) ReplyMessage(events []*linebot.Event, resMessage string) error {
	for _, event := range events {
		if event.Type == linebot.EventTypeMessage {
			_, err := c.Bot.ReplyMessage(
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

func (c *LineBotClient) SaveMessageToStore(events []*linebot.Event) error {
	for _, event := range events {
		if event.Type == linebot.EventTypeMessage {
			switch message := event.Message.(type) {
			case *linebot.TextMessage:
				userID := event.Source.UserID
				c.Store.Save(userID, message.Text, "")
				return nil
			}
		}
	}
	return fmt.Errorf("no message found in events")
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
