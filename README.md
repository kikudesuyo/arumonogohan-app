<div style="display: flex; align-items: center;">
  <img src="https://github.com/user-attachments/assets/036dec07-e68a-4309-b626-38c238bc8e4c" width="60" alt="あるものごはんのロゴ" style="margin-right: 10px;"/>
  <h1>あるものごはん</h1>
</div>

家にある残り物で、料理を提案してくれるアプリです。

## 機能

ライン Bot を用いて、対話形式で料理の提案を支援します (**2025/4/10** リリース予定です)

## 開発手順

# .env ファイルの作成

```
cp -p .env.example .env
```

`GEMINI_` については開発者のアカウントで取得してください
`LINE_BOT_` については、チャットで送信します

### ローカル開発

- ngrok のインストール
- .env に記載した PORT を同じポート番号を用いて ngrok 起動

```env
//.env
PORT=8081
```

```bash
ngrok http 8081
```

- ngrok を用いてローカルサーバーから WebhookURL を作成
- 作成した WebhookURL をhttps://developers.line.bizにて作成されたWebhookURLを設定

```bash
cd arumonogohan-app/
go mod tidy
go run cmd/main.go
```

## 技術スタック

[![My Skills](https://skillicons.dev/icons?i=go,gcp)](https://skilldev)
