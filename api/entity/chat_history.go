package entity

import (
	"time"
)

type ChatHistory struct {
	Messages  []string  // メッセージの履歴
	State     string    // チャット状態(menu_select, ingredient_input)
	Timestamp time.Time // 最新の時間
}
