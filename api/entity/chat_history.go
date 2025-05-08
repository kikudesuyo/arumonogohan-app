package entity

import (
	"time"
)

type ChatState interface {
	GetMessage() string
	GetState() string
}

type MenuCategorySelect struct {
	Category string
}

func (m MenuCategorySelect) GetMessage() string { return m.Category }
func (m MenuCategorySelect) GetState() string   { return "menu_category_select" }

type IngredientInput struct {
	Message string
}

func (i IngredientInput) GetMessage() string { return i.Message }
func (i IngredientInput) GetState() string   { return "ingredient_input" }

type ChatHistory struct {
	StateData ChatState // 状態ごとのデータ
	Timestamp time.Time // 最終更新時刻
}
