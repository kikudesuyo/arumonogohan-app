package usecase

import (
	"context"
	"fmt"

	"github.com/google/generative-ai-go/genai"
	"github.com/kikudesuyo/arumonogohan-app/api/entity"
)

type RecipeInput struct {
	MenuCategory string `json:"menu_category"`
	Ingredients  string `json:"ingredients"`
}

// SuggestRecipe は、材料を基にレシピを提案するUsecaseです。
func SuggestRecipe(ctx context.Context, input RecipeInput) (entity.Recipe, error) {
	geminiAI, err := NewGeminiAI(ctx)
	if err != nil {
		return entity.Recipe{}, fmt.Errorf("failed to create GeminiAI client: %v", err)
	}
	mealRecipe, err := geminiAI.GenerateRecipe(ctx, input)
	if err != nil {
		return entity.Recipe{}, fmt.Errorf("failed to generate recipe: %v", err)
	}
	return mealRecipe, nil
}

// GenerateRecipe は、Function Callingを使用してレシピを生成します。
func (g *GeminiAI) GenerateRecipe(ctx context.Context, input RecipeInput) (entity.Recipe, error) {
	// 1. Geminiに渡す「道具（関数）」を定義
	funcTool := &genai.Tool{
		FunctionDeclarations: []*genai.FunctionDeclaration{
			{
				Name:        "submit_recipe",
				Description: "ユーザーに提案するレシピを送信する",
				Parameters: &genai.Schema{
					Type: genai.TypeObject,
					Properties: map[string]*genai.Schema{
						"title":        {Type: genai.TypeString, Description: "料理名"},
						"ingredients":  {Type: genai.TypeArray, Items: &genai.Schema{Type: genai.TypeString}, Description: "材料のリスト"},
						"instructions": {Type: genai.TypeArray, Items: &genai.Schema{Type: genai.TypeString}, Description: "作り方の手順リスト"},
						"summary":      {Type: genai.TypeString, Description: "レシピの簡単な要約、栄養情報、コツなど"},
					},
					Required: []string{"title", "ingredients", "instructions", "summary"},
				},
			},
		},
	}

	// 2. モデルに道具をセット
	g.model.Tools = []*genai.Tool{funcTool}

	// 3. シンプルなプロンプトを作成
	prompt := fmt.Sprintf(`ユーザーが提供した食材を使って、創造的で美味しいレシピを一つ提案してください。
	食材が不十分でレシピが作れない場合は、その旨をsummaryフィールドで伝えてください。
	
	料理カテゴリ: %s
	食材: %s`, input.MenuCategory, input.Ingredients)

	// 4. コンテンツを生成
	resp, err := g.model.GenerateContent(ctx, genai.Text(prompt))
	if err != nil {
		return entity.Recipe{}, fmt.Errorf("error generating content: %v", err)
	}

	// 5. 応答からFunction Callを抽出し、Recipe構造体にマッピング
	part := resp.Candidates[0].Content.Parts[0]
	if fc, ok := part.(genai.FunctionCall); ok {
		return mapFunctionCallToRecipe(fc)
	} else {
		return entity.Recipe{}, fmt.Errorf("unexpected response from AI, no function call found")
	}
}

// mapFunctionCallToRecipe は、genai.FunctionCallをentity.Recipeに変換します。
func mapFunctionCallToRecipe(fc genai.FunctionCall) (entity.Recipe, error) {
	recipe := entity.Recipe{}

	if title, ok := fc.Args["title"].(string); ok {
		recipe.Title = title
	}

	if summary, ok := fc.Args["summary"].(string); ok {
		recipe.Summary = summary
	}

	if ingredients, ok := fc.Args["ingredients"].([]any); ok {
		for _, item := range ingredients {
			if ing, ok := item.(string); ok {
				recipe.Ingredients = append(recipe.Ingredients, ing)
			}
		}
	}

	if instructions, ok := fc.Args["instructions"].([]any); ok {
		for _, item := range instructions {
			if inst, ok := item.(string); ok {
				recipe.Instructions = append(recipe.Instructions, inst)
			}
		}
	}

	return recipe, nil
}