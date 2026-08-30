package service

import (
	"context"
	"errors"
	"testing"

	"github.com/kikudesuyo/arumonogohan-app/api/entity"
	"github.com/kikudesuyo/arumonogohan-app/api/repository"
)

type trackingChatSessionRepository struct {
	findCalled bool
}

func (r *trackingChatSessionRepository) Find(string) (*entity.ChatSession, error) {
	r.findCalled = true
	return nil, errors.New("Find should not be called for a menu category")
}

func (*trackingChatSessionRepository) Save(*entity.ChatSession) error { return nil }

func (*trackingChatSessionRepository) Delete(string) error { return nil }

func TestProcessLinebotMessageReturnsEarlyForMenuCategory(t *testing.T) {
	store := &trackingChatSessionRepository{}

	err := ProcessLinebotMessage(
		context.Background(),
		nil,
		nil,
		repository.LineUserMsg{UserID: "user-1", Msg: "家庭の味🥢"},
		store,
	)
	if err != nil {
		t.Fatalf("ProcessLinebotMessage() error = %v, want nil", err)
	}
	if store.findCalled {
		t.Fatal("ProcessLinebotMessage() called Find() for a menu category")
	}
}

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
		repository.LineUserMsg{UserID: "user-1", Msg: "材料"},
		store,
	)
	if err == nil {
		t.Fatal("ProcessLinebotMessage() error = nil, want an error")
	}
}
