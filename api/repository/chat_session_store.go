package repository

import (
	"fmt"
	"sync"
	"time"

	"github.com/kikudesuyo/arumonogohan-app/api/entity"
)

type ChatSessionStore struct {
	store sync.Map
}

func (l *ChatSessionStore) Save(sessionID string, state entity.ChatState) {
	var session entity.ChatHistory
	value, found := l.store.Load(sessionID)
	if found {
		session = value.(entity.ChatHistory)
	}
	session.StateData = state
	session.Timestamp = time.Now()

	l.store.Store(sessionID, session)
}

func (l *ChatSessionStore) Get(sessionID string) (*entity.ChatHistory, error) {
	value, ok := l.store.Load(sessionID)
	if !ok {
		return nil, fmt.Errorf("no session found for user %s", sessionID)
	}
	session := value.(entity.ChatHistory)
	// 5分以上経過したら履歴を削除
	if time.Since(session.Timestamp) > 5*time.Minute {
		l.store.Delete(sessionID)
		return nil, fmt.Errorf("session expired for user %s", sessionID)
	}

	return &session, nil
}

// Session削除
func (l *ChatSessionStore) Delete(userID string) {
	l.store.Delete(userID)
}
