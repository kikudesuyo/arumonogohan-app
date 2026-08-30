package service

import (
	"context"
	"fmt"

	"github.com/kikudesuyo/arumonogohan-app/api/entity"
	"github.com/kikudesuyo/arumonogohan-app/api/repository"
)

func SuggestRecipe(ctx context.Context, req entity.RecipeInputReq) (entity.RecipeInputResp, error) {
	geminiAI, err := repository.NewGeminiAI(ctx)
	if err != nil {
		return entity.RecipeInputResp{}, fmt.Errorf("failed to create GeminiAI client: %v", err)
	}
	isTempered, err := isPromptTempered(geminiAI, ctx, req)
	if err != nil {
		return entity.RecipeInputResp{}, fmt.Errorf("prompt tampering detected: %v", err)
	}
	if isTempered {
		return entity.RecipeInputResp{
			IsTempered: true,
		}, nil
	}

	tools := []repository.Tool{
		repository.MakeToolFromStruct[entity.RecipeInputResp]("submit_recipe", "ユーザーに提案するレシピを送信する"),
	}

	prompt := fmt.Sprintf(
		`ユーザーが提供した食材を使って、創造的で美味しいレシピを一つ提案してください。
		食材が不十分でレシピが作れない場合は、その旨をsummaryフィールドで伝えてください。
		料理カテゴリ: %s
		食材: %s`,
		req.MenuCategory,
		req.Ingredients,
	)

	aiRecipe, err := repository.GetGeminiJSONResp[entity.RecipeInputResp](geminiAI, ctx, prompt, tools)
	if err != nil {
		return entity.RecipeInputResp{}, fmt.Errorf("AI generate failed: %v", err)
	}
	return aiRecipe, nil
}
