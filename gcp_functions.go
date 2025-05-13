package gcp

import (
	"net/http"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"github.com/gin-gonic/gin"
	"github.com/kikudesuyo/arumonogohan-app/api/handler"
)

func init() {
	functions.HTTP("LinbotCallback", func(w http.ResponseWriter, req *http.Request) {
		r := gin.Default()
		r.POST("/callback", handler.HandleLinebotCallback)
		r.ServeHTTP(w, req)
	})
}
