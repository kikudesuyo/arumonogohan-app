package usecase

import (
	"fmt"

	"github.com/line/line-bot-sdk-go/linebot"
)

func (l *LineBotClient) ReplyMsg(events []*linebot.Event, msg string) error {
	for _, event := range events {
		if event.Type == linebot.EventTypeMessage {
			_, err := l.Bot.ReplyMessage(
				event.ReplyToken,
				linebot.NewTextMessage(msg),
			).Do()
			if err != nil {
				return fmt.Errorf("error sending reply message: %v", err)
			}
		}
	}
	return nil
}
