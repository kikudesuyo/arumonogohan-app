package service

import (
	"context"
	"testing"

	"github.com/kikudesuyo/arumonogohan-app/api/entity"
	"github.com/kikudesuyo/arumonogohan-app/api/infrastructure"
	"github.com/kikudesuyo/arumonogohan-app/api/repository"
)

func TestProcessLinebotMessageReturnsErrorForUnknownState(t *testing.T) {
	store := repository.NewInMemoryChatSessionRepository()
	session := entity.NewChatSession("user-1")
	session.State = entity.ChatState("unknown")
	if err := store.Save(session); err != nil {
		t.Fatalf("Save() error = %v", err)
	}

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
