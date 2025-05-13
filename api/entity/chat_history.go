package entity

import (
	"time"
)

type ChatState string

const (
	StateMenuCategorySelect ChatState = "menu_category_select"
	StateIngredientInput    ChatState = "ingredient_input"
)

type ChatHistory struct {
	State     ChatState `json:"state_data"` // 状態データ
	Msg       string    `json:"msg"`        // メッセージ
	Timestamp time.Time `json:"timestamp"`  // 更新日時
}
