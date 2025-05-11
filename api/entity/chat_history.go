package entity

import (
	"time"
)

type ChatState interface {
	GetMsg() string
	GetState() string
}

type MenuCategorySelect struct {
	Category string
}

func (m MenuCategorySelect) GetMsg() string   { return m.Category }
func (m MenuCategorySelect) GetState() string { return "menu_category_select" }

type IngredientInput struct {
	Msg string
}

func (i IngredientInput) GetMsg() string   { return i.Msg }
func (i IngredientInput) GetState() string { return "ingredient_input" }

type ChatHistory struct {
	StateData ChatState // 状態ごとのデータ
	Timestamp time.Time // 最終更新時刻
}
