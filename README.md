<div style="display: flex; align-items: center;">
  <img src="https://github.com/user-attachments/assets/036dec07-e68a-4309-b626-38c238bc8e4c" width="60" alt="あるものごはんのロゴ" style="margin-right: 10px;"/>
  <h1>あるものごはん</h1>
</div>

#### 概要

家にある残り物で、料理のレシピを提案してくれるアプリです。

#### [あるものごはん](https://lin.ee/pLzG7zn)

### QR コード

<img src="./assets/line-bot-qr.png" width=160, alt="LINE公式アカウントのQRコード">

## 機能

- Google Gemini API を用いて、食材の情報を取得
- LINE Bot を用いて、ユーザーとの対話形式でレシピを提案

## 開発環境

- Go
- GCP Cloud Run
- Google Gemini API
- LINE Messaging API

## 開発手順

### ローカル開発

#### .env ファイル作成

```
cp -p .env.example .env
```

`GEMINI_` については開発者の Google アカウントで発効してください
`LINE_BOT_` については、個別で渡します。

#### localhost を外部公開

- ngrok のインストール
- .env に記載した PORT を同じポート番号を用いて ngrok 起動

```env
//.env
PORT=8081
```

```bash
ngrok http 8081
```

- ngrok を用いてローカルサーバーから`WebhookURL`を作成
- 作成した`WebhookURL`を`https://developers.line.biz`にて作成された`WebhookURL`を設定

### 実行

```bash
cd arumonogohan-app/
go mod tidy
go run cmd/main.go
```

## 本番環境

GCP Cloud Run にデプロイ

```bash
bash deploy_cloud_functions.sh
```

### 備考

新しい環境変数を追加する際には、以下の操作を行ってください。

```text
1. .env に新しい変数を追加
1. deploy_cloud_functions.sh にて登録した変数の追加
```

## 技術スタック

[![My Skills](https://skillicons.dev/icons?i=go,gcp)](https://skilldev)
