package gcp

import (
	"net/http"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"github.com/gin-gonic/gin"
	"github.com/kikudesuyo/arumonogohan-app/api/handler"
)

func init() {
	functions.HTTP("LinbotCallback", func(w http.ResponseWriter, req *http.Request) {
		// Ginのルーターを作成
		r := gin.Default()
		// Ginのハンドラを定義
		r.POST("/callback", handler.HandleLinebotCallback)

		// リクエストをGinに渡して処理させる
		r.ServeHTTP(w, req)
	})
}
