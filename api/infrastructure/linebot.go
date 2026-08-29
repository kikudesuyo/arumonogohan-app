package infrastructure

import (
	"fmt"
	"net/http"
	"os"

	"github.com/line/line-bot-sdk-go/linebot"
)

type LinebotContext struct {
	Bot    *linebot.Client
	Events []*linebot.Event
}

type LineUserMsg struct {
	UserID string
	Msg    string
}

func NewLineBotClient() (*linebot.Client, error) {
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
	return bot, nil
}

func ParseLinebotRequest(r *http.Request, bot *linebot.Client) ([]*linebot.Event, error) {
	events, err := bot.ParseRequest(r)
	if err != nil {
		return nil, fmt.Errorf("failed to parse request: %v", err)
	}
	return events, nil
}

func GetLineUserMsg(events []*linebot.Event) (LineUserMsg, error) {
	for _, event := range events {
		if event.Type == linebot.EventTypeMessage {
			// メッセージがテキスト型の場合
			if msg, ok := event.Message.(*linebot.TextMessage); ok {
				return LineUserMsg{
					UserID: event.Source.UserID,
					Msg:    msg.Text,
				}, nil
			}
		}
	}
	return LineUserMsg{}, fmt.Errorf("no text message found in events")
}

func ReplyMsgToLine(bot *linebot.Client, events []*linebot.Event, msg string) error {
	for _, event := range events {
		if event.Type == linebot.EventTypeMessage {
			_, err := bot.ReplyMessage(
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
