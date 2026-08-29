package usecase

import (
	"context"
	"testing"
	"time"

	"github.com/kikudesuyo/arumonogohan-app/api/entity"
	"github.com/kikudesuyo/arumonogohan-app/api/infrastructure"
	"github.com/kikudesuyo/arumonogohan-app/api/repository"
)

func TestProcessLinebotMessageReturnsErrorForUnknownState(t *testing.T) {
	store := &repository.ChatSessionStore{}
	store.UpsertChatSession(repository.ChatSession{
		SessionID: "user-1",
		State:     entity.ChatState("unknown"),
		Timestamp: time.Now(),
	})

	err := ProcessLinebotMessage(
		context.Background(),
		nil,
		nil,
		infrastructure.LineUserMsg{UserID: "user-1", Msg: "材料"},
		store,
	)
	if err == nil {
		t.Fatal("ProcessLinebotMessage() error = nil, want an error")
	}
}
