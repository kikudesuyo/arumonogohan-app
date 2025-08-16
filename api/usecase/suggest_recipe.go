package usecase

import (
	"context"
	"fmt"

	"github.com/kikudesuyo/arumonogohan-app/api/entity"
	"github.com/kikudesuyo/arumonogohan-app/api/infrastructure"
)

func SuggestRecipe(ctx context.Context, req entity.RecipeInputReq) (entity.RecipeInputResp, error) {
	geminiAI, err := infrastructure.NewGeminiAI(ctx)
	if err != nil {
		return entity.RecipeInputResp{}, fmt.Errorf("failed to create GeminiAI client: %v", err)
	}

	// AIRecipe の構造から自動でツールを生成
	tools := []infrastructure.Tool{
		infrastructure.MakeToolFromStruct[entity.RecipeInputResp]("submit_recipe", "ユーザーに提案するレシピを送信する"),
	}

	// プロンプトを作成
	prompt := fmt.Sprintf(
		`ユーザーが提供した食材を使って、創造的で美味しいレシピを一つ提案してください。
		食材が不十分でレシピが作れない場合は、その旨をsummaryフィールドで伝えてください。
		料理カテゴリ: %s
		食材: %s`,
		req.MenuCategory,
		req.Ingredients,
	)

	// Infra 層に投げる（ビジネスロジックはここでは扱わない）
	aiRecipe, err := infrastructure.GetJSONResp[entity.RecipeInputResp](geminiAI, ctx, prompt, tools)
	if err != nil {
		return entity.RecipeInputResp{}, fmt.Errorf("AI generate failed: %v", err)
	}

	// AI のレスポンスを Usecase の構造体に変換
	return entity.RecipeInputResp{
		Title:        aiRecipe.Title,
		Summary:      aiRecipe.Summary,
		Ingredients:  aiRecipe.Ingredients,
		Instructions: aiRecipe.Instructions,
	}, nil
}
