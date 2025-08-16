package repository

import (
	"sync"
	"time"

	"github.com/kikudesuyo/arumonogohan-app/api/entity"
)

type ChatSessionStore struct {
	store sync.Map
}

type ChatSession struct {
	SessionID    string
	MenuCategory string
	State        entity.ChatState
	Timestamp    time.Time
}

func (l *ChatSessionStore) GetChatSession(sessionID string) (*ChatSession, bool) {
	value, found := l.store.Load(sessionID)
	if !found {
		return nil, false
	}
	session := value.(ChatSession)
	// 5分以上経過したら履歴を削除
	if time.Since(session.Timestamp) > 5*time.Minute {
		l.store.Delete(sessionID)
		return nil, false
	}
	return &session, true
}

func (l *ChatSessionStore) UpsertChatSession(session ChatSession) {
	l.store.Store(session.SessionID, session)
}

func (l *ChatSessionStore) DeleteChatSession(sessionID string) {
	l.store.Delete(sessionID)
}

func (l *ChatSessionStore) InsertInitChatSession(sessionID string) {
	session := ChatSession{
		SessionID:    sessionID,
		MenuCategory: "",
		State:        entity.StateMenuCategorySelect,
		Timestamp:    time.Now(),
	}
	l.UpsertChatSession(session)
}
