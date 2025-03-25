package handlers

import (
	"context"
	"fmt"

	"github.com/kikudesuyo/arumonogohan-app/api/entity"
)

func HandleSuggestRecipe(clientMessage string) (string, error) {
	geminiAI, err := entity.NewGeminiAI()
	if err != nil {
		fmt.Println(err.Error())
		return "", err
	}
	ctx := context.Background()
	mealRecipe, err := geminiAI.GenerateRecipe(ctx, clientMessage)
	if err != nil {
		fmt.Println(err.Error())
		return "", err
	}
	return mealRecipe, nil
}
