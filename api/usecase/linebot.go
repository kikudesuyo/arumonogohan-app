package usecase

import (
	"fmt"
	"os"

	"github.com/line/line-bot-sdk-go/linebot"
)

type LineBotClient struct {
	Bot *linebot.Client
}

type LineUserMsg struct {
	UserID string
	Msg    string
}

func NewLineBotClient() (*LineBotClient, error) {
	channelSecret := os.Getenv("LINE_BOT_CHANNEL_SECRET")
	channelToken := os.Getenv("LINE_BOT_CHANNEL_TOKEN")
	if channelSecret == "" || channelToken == "" {
		return nil, fmt.Errorf("LINE_BOT_CHANNEL_SECRET or LINE_BOT_CHANNEL_TOKEN is not set")
	}
	bot, err := linebot.New(
		channelSecret,
		channelToken,
	)
	if err != nil {
		return nil, fmt.Errorf("error creating LINE bot client: %v", err)
	}
	return &LineBotClient{Bot: bot}, nil
}

func (c *LineBotClient) GetLineMsg(events []*linebot.Event) (*LineUserMsg, error) {
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
