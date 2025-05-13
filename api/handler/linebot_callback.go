package handler

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/kikudesuyo/arumonogohan-app/api/entity"
	"github.com/kikudesuyo/arumonogohan-app/api/repository"
	"github.com/kikudesuyo/arumonogohan-app/api/usecase"
)

var store = &repository.ChatSessionStore{}

// HandleLinebotCallback processes LINE webhook callbacks
func HandleLinebotCallback(c *gin.Context) {
	userAgent := c.GetHeader("User-Agent")
	if !strings.Contains(userAgent, "LineBotWebhook") {
		return
	}

	lineMsgCtx, err := parseLineRequest(c.Request)
	if err != nil {
		logError("parse line request", err)
		return
	}

	userID := lineMsgCtx.UserMsg.UserID
	msg := lineMsgCtx.UserMsg.Msg

	chatSession, found := store.Get(userID)
	if !found {
		chatSession = createNewSession(userID)
	}

	// Process the message based on current state
	switch chatSession.State {
	case entity.StateMenuCategorySelect:
		handleMenuCategorySelect(lineMsgCtx, chatSession, msg)
	case entity.StateIngredientInput:
		handleIngredientInput(lineMsgCtx, chatSession, msg)
	}
}

// createNewSession creates a new chat session for a user
func createNewSession(userID string) *repository.ChatSession {
	fmt.Println("session not found. creating new session")
	chatSession := &repository.ChatSession{
		SessionID:    userID,
		MenuCategory: "",
		State:        entity.StateMenuCategorySelect,
		Timestamp:    time.Now(),
	}
	store.Save(*chatSession)
	return chatSession
}

// handleMenuCategorySelect processes messages in the menu category selection state
func handleMenuCategorySelect(lineMsgCtx *usecase.LineMsgContext, chatSession *repository.ChatSession, msg string) {
	if !entity.IsMenuCategorySelected(msg) {
		replyToUser(lineMsgCtx, "メニューから料理するジャンルを選択ください🍽️")
		return
	}

	// Update session with selected category
	chatSession.MenuCategory = msg
	chatSession.State = entity.StateIngredientInput
	chatSession.Timestamp = time.Now()
	store.Save(*chatSession)

	// Ask for ingredients
	replyMsg := fmt.Sprintf("「%s」ですね✨️ 使う食材を教えて下さい!!", msg)
	replyToUser(lineMsgCtx, replyMsg)
}

// handleIngredientInput processes messages in the ingredient input state
func handleIngredientInput(lineMsgCtx *usecase.LineMsgContext, chatSession *repository.ChatSession, msg string) {
	// Check if user is selecting a different menu category
	if entity.IsMenuCategorySelected(msg) {
		handleMenuCategoryReselection(lineMsgCtx, chatSession, msg)
		return
	}

	// Process ingredients and suggest recipe
	recipeInput := usecase.RecipeInput{
		MenuCategory: chatSession.MenuCategory,
		Ingredients:  msg,
	}
	
	replyMsg, err := usecase.SuggestRecipe(recipeInput)
	if err != nil {
		logError("suggest recipe", err)
		return
	}
	
	// Reset session state
	chatSession.State = entity.StateMenuCategorySelect
	chatSession.MenuCategory = ""
	chatSession.Timestamp = time.Now()
	store.Save(*chatSession)

	replyToUser(lineMsgCtx, replyMsg)
}

// handleMenuCategoryReselection handles when a user selects a different menu category
func handleMenuCategoryReselection(lineMsgCtx *usecase.LineMsgContext, chatSession *repository.ChatSession, msg string) {
	chatSession.MenuCategory = msg
	chatSession.State = entity.StateIngredientInput
	chatSession.Timestamp = time.Now()
	store.Save(*chatSession)

	replyMsg := fmt.Sprintf("「%s」ですね✨️ 使う食材を教えて下さい!!", msg)
	replyToUser(lineMsgCtx, replyMsg)
}

// replyToUser sends a reply message to the user
func replyToUser(lineMsgCtx *usecase.LineMsgContext, message string) {
	err := usecase.ReplyMsgToLine(lineMsgCtx.Bot, lineMsgCtx.Events, message)
	if err != nil {
		logError("reply message", err)
	}
}

// logError logs an error with context
func logError(context string, err error) {
	fmt.Printf("Error in %s: %v\n", context, err)
}

// parseLineRequest parses the LINE webhook request
func parseLineRequest(r *http.Request) (*usecase.LineMsgContext, error) {
	bot, err := usecase.NewLineBotClient()
	if err != nil {
		return nil, fmt.Errorf("failed to create LINE bot client: %v", err)
	}
	events, err := bot.ParseRequest(r)
	if err != nil {
		return nil, fmt.Errorf("failed to parse request: %v", err)
	}
	msg, err := usecase.GetLineMsg(events)
	if err != nil {
		return nil, fmt.Errorf("failed to get line message: %v", err)
	}
	return &usecase.LineMsgContext{
		Bot:     bot,
		Events:  events,
		UserMsg: msg,
	}, nil
}
