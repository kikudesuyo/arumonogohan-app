package v1

import (
	"encoding/json"
	"net/http"

	"github.com/kikudesuyo/arumonogohan-app/api/entity"
	"github.com/kikudesuyo/arumonogohan-app/api/service"
	"github.com/kikudesuyo/arumonogohan-app/api/xerror"
)

func HandleSuggestRecipe(r *http.Request, requestData map[string]interface{}) (http.Handler, error) {
	var req entity.RecipeInputReq
	requestJSON, err := json.Marshal(requestData)
	if err != nil {
		return nil, xerror.ClientValidationErr(err)
	}
	if err := json.Unmarshal(requestJSON, &req); err != nil {
		return nil, xerror.ClientValidationErr(err)
	}
	recipe, err := service.SuggestRecipe(r.Context(), req)
	if err != nil {
		return nil, err
	}
	return entity.NewObjectResponse(map[string]any{"recipe": recipe}), nil
}
