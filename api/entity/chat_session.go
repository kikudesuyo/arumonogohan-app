package entity

import "time"

// ChatSession はユーザーとの対話セッションを表す
type ChatSession struct {
	UserID       string
	MenuCategory string
	State        ChatState
	Timestamp    time.Time
}

// NewChatSession は新しいチャットセッションを作成する
func NewChatSession(userID string) *ChatSession {
	return &ChatSession{
		UserID:    userID,
		State:     StateMenuCategorySelect,
		Timestamp: time.Now(),
	}
}

// IsExpired はセッションが期限切れかどうかを判定する
func (cs *ChatSession) IsExpired() bool {
	// 5分以上経過したら期限切れ
	return time.Since(cs.Timestamp) > 5*time.Minute
}

// UpdateMenuCategory はメニューカテゴリを更新し、状態を食材入力中に変更する
func (cs *ChatSession) UpdateMenuCategory(category string) {
	cs.MenuCategory = category
	cs.State = StateIngredientInput
	cs.Timestamp = time.Now()
}

// Reset はセッションを初期状態に戻す
func (cs *ChatSession) Reset() {
	cs.MenuCategory = ""
	cs.State = StateMenuCategorySelect
	cs.Timestamp = time.Now()
}
