// /entity/linesession.go
package entity

import (
	"fmt"
	"sync"
	"time"
)

const MaxLogs = 5

type LineSession struct {
	Messages  []string  // 最新5件のメッセージ
	State     string    // チャット状態(1:メニュー選択, 2:食材入力)
	Timestamp time.Time // 最新の時間
}

type LineSessionStore struct {
	store sync.Map
}

func (l *LineSessionStore) Save(userID, message, state string) {
	value, found := l.store.Load(userID)
	var session LineSession

	if found {
		session = value.(LineSession)
	} else {
		session = LineSession{Messages: []string{}, State: state}
	}
	if len(session.Messages) >= MaxLogs {
		session.Messages = session.Messages[1:] // 最も古いメッセージを削除
	}

	session.Messages = append(session.Messages, message)
	session.State = state
	session.Timestamp = time.Now()

	l.store.Store(userID, session)
}

// LineSessionを取得
func (l *LineSessionStore) Get(userID string) (*LineSession, error) {
	value, ok := l.store.Load(userID)
	if !ok {
		return nil, fmt.Errorf("no session found for user %s", userID)
	}
	session := value.(LineSession)
	// 5分以上経過したら履歴を削除
	if time.Since(session.Timestamp) > 5*time.Minute {
		l.store.Delete(userID)
		return nil, fmt.Errorf("session expired for user %s", userID)
	}

	return &session, nil
}

// LineSession削除
func (l *LineSessionStore) Delete(userID string) {
	l.store.Delete(userID)
}
