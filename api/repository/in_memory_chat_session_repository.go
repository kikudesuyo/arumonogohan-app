package repository

import (
	"sync"

	"github.com/kikudesuyo/arumonogohan-app/api/entity"
)

// InMemoryChatSessionRepository は sync.Map を使用したインメモリの ChatSessionRepository
type InMemoryChatSessionRepository struct {
	store sync.Map
}

// NewInMemoryChatSessionRepository は新しい InMemoryChatSessionRepository を生成する
func NewInMemoryChatSessionRepository() ChatSessionRepository {
	return &InMemoryChatSessionRepository{}
}

// Find はユーザーIDでセッションを検索する。見つからない、または期限切れの場合は nil を返す
func (r *InMemoryChatSessionRepository) Find(userID string) (*entity.ChatSession, error) {
	value, found := r.store.Load(userID)
	if !found {
		return nil, nil // 見つからない場合はエラーではなく nil
	}

	session := value.(*entity.ChatSession)
	if session.IsExpired() {
		r.store.Delete(userID)
		return nil, nil
	}

	return session, nil
}

// Save はセッションを保存する
func (r *InMemoryChatSessionRepository) Save(session *entity.ChatSession) error {
	r.store.Store(session.UserID, session)
	return nil
}

// Delete はセッションを削除する
func (r *InMemoryChatSessionRepository) Delete(userID string) error {
	r.store.Delete(userID)
	return nil
}
