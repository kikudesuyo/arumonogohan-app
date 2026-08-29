package route

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/kikudesuyo/arumonogohan-app/api/handler"
	v1 "github.com/kikudesuyo/arumonogohan-app/api/handler/v1"
)

func NewRouter() http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Post("/callback", v1.HandleLinebotCallback)
	r.Post("/suggest", func(w http.ResponseWriter, r *http.Request) {
		handler.HandleRequestAndResponse(r, w, v1.HandleSuggestRecipe)
	})
	return r
}
