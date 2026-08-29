package service

import (
	"testing"

	"github.com/kikudesuyo/arumonogohan-app/api/entity"
)

func TestFormatRecipeForLine(t *testing.T) {
	tests := []struct {
		name   string
		recipe entity.RecipeInputResp
		want   string
	}{
		{
			name: "formats a recipe",
			recipe: entity.RecipeInputResp{
				Title:        "野菜炒め",
				Ingredients:  []string{"キャベツ", "卵"},
				Instructions: []string{"切る", "炒める"},
				Summary:      "強火で仕上げる",
			},
			want: "今日のレシピは「野菜炒め」で決まり！\n\n【材料】\n- キャベツ\n- 卵\n\n【作り方】\n1. 切る\n2. 炒める\n\n【ポイント】\n強火で仕上げる\n",
		},
		{
			name: "returns summary for an unavailable recipe",
			recipe: entity.RecipeInputResp{
				Title:   "提案できません",
				Summary: "食材が不足しています",
			},
			want: "食材が不足しています",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := formatRecipeForLine(tt.recipe); got != tt.want {
				t.Fatalf("formatRecipeForLine() = %q, want %q", got, tt.want)
			}
		})
	}
}
