package repository

import (
	"fmt"
	"sync"
	"time"

	"github.com/kikudesuyo/arumonogohan-app/api/entity"
)

const MaxLogs = 5

type ChatSessionStore struct {
	store sync.Map
}

func (l *ChatSessionStore) Save(userID, message, state string) {
	var session entity.ChatHistory
	value, found := l.store.Load(userID)
	if found {
		session = value.(entity.ChatHistory)
	} else {
		session = entity.ChatHistory{Messages: []string{}, State: state}
	}
	if len(session.Messages) >= MaxLogs {
		session.Messages = session.Messages[1:] // 最も古いメッセージを削除
	}
	session.Messages = append(session.Messages, message)
	session.State = state
	session.Timestamp = time.Now()

	l.store.Store(userID, session)
}

// Session取得
func (l *ChatSessionStore) Get(userID string) (*entity.ChatHistory, error) {
	value, ok := l.store.Load(userID)
	if !ok {
		return nil, fmt.Errorf("no session found for user %s", userID)
	}
	session := value.(entity.ChatHistory)
	// 5分以上経過したら履歴を削除
	if time.Since(session.Timestamp) > 5*time.Minute {
		l.store.Delete(userID)
		return nil, fmt.Errorf("session expired for user %s", userID)
	}

	return &session, nil
}

// Session削除
func (l *ChatSessionStore) Delete(userID string) {
	l.store.Delete(userID)
}
