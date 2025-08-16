package entity

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
}
