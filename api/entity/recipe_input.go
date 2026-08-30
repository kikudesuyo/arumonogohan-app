package entity

import (
	"bytes"
	"encoding/json"
	"strings"
)

type RecipeInputReq struct {
	MenuCategory string `json:"menu_category"`
	Ingredients  string `json:"ingredients"`
}

// RecipeInputResp はレシピの詳細な構造を定義します。
type RecipeInputResp struct {
	Title        string   `json:"title"`
	Ingredients  []string `json:"ingredients"`
	Instructions []string `json:"instructions"`
	Summary      string   `json:"summary,omitempty"`
	IsTempered   bool     `json:"is_tempered"`
}

// UnmarshalJSON accepts both JSON booleans and the string values sometimes
// returned by an LLM for boolean fields.
func (r *RecipeInputResp) UnmarshalJSON(data []byte) error {
	type recipeInputResp RecipeInputResp
	var raw struct {
		Title        string          `json:"title"`
		Ingredients  []string        `json:"ingredients"`
		Instructions []string        `json:"instructions"`
		Summary      string          `json:"summary"`
		IsTempered   json.RawMessage `json:"is_tempered"`
	}
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	var tempered bool
	if len(bytes.TrimSpace(raw.IsTempered)) > 0 {
		if err := json.Unmarshal(raw.IsTempered, &tempered); err != nil {
			var value string
			if json.Unmarshal(raw.IsTempered, &value) == nil {
				tempered = strings.EqualFold(strings.TrimSpace(value), "true")
			}
		}
	}
	*r = RecipeInputResp(recipeInputResp{Title: raw.Title, Ingredients: raw.Ingredients, Instructions: raw.Instructions, Summary: raw.Summary, IsTempered: tempered})
	return nil
}
