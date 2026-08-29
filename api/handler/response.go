package handler

import "net/http"

// Renderer は API レスポンスの HTTP 出力を抽象化します。
type Renderer interface {
	ServeHTTP(http.ResponseWriter, *http.Request)
}
