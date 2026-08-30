package v1

import (
	"log"
	"net/http"
	"strings"

	"github.com/kikudesuyo/arumonogohan-app/api/repository"
	"github.com/kikudesuyo/arumonogohan-app/api/service"
)

var chatSessions = repository.NewInMemoryChatSessionRepository()

func HandleLinebotCallback(w http.ResponseWriter, r *http.Request) {
	if !strings.Contains(r.Header.Get("User-Agent"), "LineBotWebhook") {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	linebot, err := repository.NewLineBotClient()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	events, err := repository.ParseLinebotRequest(r, linebot)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	lineUserMsg, err := repository.GetLineUserMsg(events)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := service.ProcessLinebotMessage(r.Context(), linebot, events, lineUserMsg, chatSessions); err != nil {
		log.Printf("linebot callback failed: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
