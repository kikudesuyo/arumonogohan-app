package entity

var MenuCategories = map[string]struct{}{
	"時短メニュー⏱️":  {},
	"家庭の味🥢":     {},
	"さっぱりヘルシー🥗": {},
	"ガッツリメニュー🍖": {},
}

func IsMenuCategorySelected(msg string) bool {
	_, exists := MenuCategories[msg]
	return exists
}
