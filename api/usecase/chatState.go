package usecase

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/kikudesuyo/arumonogohan-app/api/entity"
	"github.com/kikudesuyo/arumonogohan-app/api/infrastructure"
	"github.com/kikudesuyo/arumonogohan-app/api/repository"
	"github.com/line/line-bot-sdk-go/linebot"
)

func ProcessSelectMenuCategory(bot *linebot.Client, events []*linebot.Event, lineUserMsg infrastructure.LineUserMsg, store *repository.ChatSessionStore) error {
	replyMsg := fmt.Sprintf("「%s」ですね✨️ \n次に使う食材を教えて下さい!👨‍🍳", lineUserMsg.Msg)
	err := infrastructure.ReplyMsgToLine(bot, events, replyMsg)
	if err != nil {
		return err
	}

	//状態更新
	chatSession := &repository.ChatSession{
		SessionID:    lineUserMsg.UserID,
		MenuCategory: lineUserMsg.Msg,
		State:        entity.StateIngredientInput,
		Timestamp:    time.Now(),
	}
	store.UpsertChatSession(*chatSession)
	return nil

}

func ProcessInputIngredient(bot *linebot.Client, events []*linebot.Event, lineUserMsg infrastructure.LineUserMsg, chatSession *repository.ChatSession, store *repository.ChatSessionStore) error {
	// メニューカテゴリ再選択の場合
	if entity.IsMenuCategorySelected(lineUserMsg.Msg) {
		chatSession.MenuCategory = lineUserMsg.Msg
		chatSession.State = entity.StateIngredientInput
		chatSession.Timestamp = time.Now()
		store.UpsertChatSession(*chatSession)

		replyMsg := fmt.Sprintf("「%s」ですね✨️ 使う食材を教えて下さい!!", lineUserMsg.Msg)
		err := infrastructure.ReplyMsgToLine(bot, events, replyMsg)
		if err != nil {
			return err
		}
		return nil
	}

	recipeInput := RecipeInput{
		MenuCategory: chatSession.MenuCategory,
		Ingredients:  lineUserMsg.Msg,
	}
	ctx := context.Background()
	recipe, err := SuggestRecipe(ctx, recipeInput)
	if err != nil {
		return err
	}

	// Recipe構造体をLINE用の文字列にフォーマットする
	replyMsg := formatRecipeForLine(recipe)
	err = infrastructure.ReplyMsgToLine(bot, events, replyMsg)
	if err != nil {
		return err
	}
	store.InsertInitChatSession(lineUserMsg.UserID)
	return nil
}

// formatRecipeForLine は、Recipe構造体をLINEメッセージ用の整形済み文字列に変換します。
func formatRecipeForLine(recipe entity.Recipe) string {
	// レシピが生成できなかった場合やエラーの場合は、サマリーメッセージのみを返す
	if recipe.Title == "提案できません" || recipe.Title == "無効な入力です" || recipe.Title == "エラー" {
		return recipe.Summary
	}

	var builder strings.Builder
	builder.WriteString(fmt.Sprintf("今日のレシピは「%s」で決まり！\n\n", recipe.Title))

	if len(recipe.Ingredients) > 0 {
		builder.WriteString("【材料】\n")
		for _, ingredient := range recipe.Ingredients {
			builder.WriteString(fmt.Sprintf("- %s\n", ingredient))
		}
		builder.WriteString("\n")
	}

	if len(recipe.Instructions) > 0 {
		builder.WriteString("【作り方】\n")
		for i, instruction := range recipe.Instructions {
			builder.WriteString(fmt.Sprintf("%d. %s\n", i+1, instruction))
		}
		builder.WriteString("\n")
	}

	if recipe.Summary != "" {
		builder.WriteString(fmt.Sprintf("【ポイント】\n%s\n", recipe.Summary))
	}

	return builder.String()
}
