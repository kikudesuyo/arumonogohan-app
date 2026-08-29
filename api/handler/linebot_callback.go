package handler

import (
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/kikudesuyo/arumonogohan-app/api/infrastructure"
	"github.com/kikudesuyo/arumonogohan-app/api/repository"
	"github.com/kikudesuyo/arumonogohan-app/api/usecase"
)

var store = &repository.ChatSessionStore{}

func HandleLinebotCallback(c *gin.Context) {
	userAgent := c.GetHeader("User-Agent")
	if !strings.Contains(userAgent, "LineBotWebhook") {
		return
	}

	linebot, err := infrastructure.NewLineBotClient()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	events, err := infrastructure.ParseLinebotRequest(c.Request, linebot)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	lineUserMsg, err := infrastructure.GetLineUserMsg(events)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	if err := usecase.ProcessLinebotMessage(c.Request.Context(), linebot, events, lineUserMsg, store); err != nil {
		fmt.Println(err.Error())
	}
}
