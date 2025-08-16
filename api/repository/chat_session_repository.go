package repository

import "github.com/kikudesuyo/arumonogohan-app/api/entity"

// ChatSessionRepository はチャットセッションの永続化を責務に持つ
type ChatSessionRepository interface {
	Find(userID string) (*entity.ChatSession, error)
	Save(session *entity.ChatSession) error
	Delete(userID string) error
}
