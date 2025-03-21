package main

import (
	"bufio"
	"context"
	"fmt"
	"os"

	"github.com/google/generative-ai-go/genai"
	"github.com/joho/godotenv"
	"google.golang.org/api/option"
)

func main() {
	if err := godotenv.Load(".env"); err != nil {
		fmt.Println("No .env file found")
	}
	geminiAPIKey := os.Getenv("GEMINI_API_KEY")
	if geminiAPIKey == "" {
		fmt.Println("GEMINI_API_KEY is not set")
		return
	}

	fmt.Println("家にある食材をカンマ区切りで入力してください (例: 卵, キャベツ, ツナ缶):")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	ingredients := scanner.Text()

	ctx := context.Background()
	client, err := genai.NewClient(ctx, option.WithAPIKey(geminiAPIKey))
	if err != nil {
		fmt.Println("Failed to create Gemini client:", err)
		return
	}
	defer client.Close()

	model := client.GenerativeModel("gemini-1.5-flash")
	prompt := fmt.Sprintf(`あなたはプロの料理人です。
ユーザーが家にある食材を入力すると、その食材を活用した美味しくて簡単なレシピを提案してください。
【要件】
1. レシピはシンプルで調理が簡単なものを提案してください。
2. 最低3つの食材を活用し、なるべく少ない材料で作れるように工夫してください。
3. すべての食材を使わなくても構いませんが、できるだけ多くの食材を活用するようにしてください。
4. 具体的な作り方（手順）をステップ形式で説明してください。
5. 可能なら追加のアレンジも提案してください（例: 「〇〇を加えるとさらに美味しくなります！」）。
6. 日本の家庭でよく使われる調味料（醤油・味噌・塩・砂糖など）を前提としてレシピを考えてください。
7. カロリーや栄養面のポイントも簡単に述べてください（例：「高たんぱくでヘルシー」）。
入力された食材: %s`, ingredients)

	resp, err := model.GenerateContent(ctx, genai.Text(prompt))
	if err != nil {
		fmt.Println("Error generating content:", err)
		return
	}

	for _, cand := range resp.Candidates {
		if cand.Content == nil {
			continue
		}
		for _, part := range cand.Content.Parts {
			fmt.Println(part)
		}
	}
}
