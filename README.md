<div style="display: flex; align-items: center;">
  <img src="https://github.com/user-attachments/assets/036dec07-e68a-4309-b626-38c238bc8e4c" width="60" alt="あるものごはんのロゴ" style="margin-right: 10px;"/>
  <h1>あるものごはん</h1>
</div>

#### 概要

食べたい料理のカテゴリと、使用したい食材を LINE で Bot に送信すると、それに応じたレシピを提案してくれるサービスです。

<div style="display: flex; flex-wrap: wrap; gap: 10px;">
  <img src="https://qiita-image-store.s3.ap-northeast-1.amazonaws.com/0/3767231/77d770c9-ad66-4f49-9be0-d1e438076e50.jpeg" width="300" />
  <img src="https://qiita-image-store.s3.ap-northeast-1.amazonaws.com/0/3767231/ff9d9c9a-5d5f-4169-9331-68ccc09a294e.png" width="300" />
  <img src="https://qiita-image-store.s3.ap-northeast-1.amazonaws.com/0/3767231/ef622424-cbbc-443b-8626-109b28059a00.png" width="300" />
  <img src="https://qiita-image-store.s3.ap-northeast-1.amazonaws.com/0/3767231/933e34bc-6412-4bf0-8122-922c126868e4.jpeg" width="300" />
  <img src="https://qiita-image-store.s3.ap-northeast-1.amazonaws.com/0/3767231/5539d521-7083-44b8-87ba-ee8f0e363a25.png" width="300" />
  <img src="https://qiita-image-store.s3.ap-northeast-1.amazonaws.com/0/3767231/c766122e-e8f4-4dd1-8ca3-37327af25cba.png" width="300" />
  <img src="https://qiita-image-store.s3.ap-northeast-1.amazonaws.com/0/3767231/45d7a6a9-c78b-47be-ab02-964370c3c8ec.png" width="300" />
</div>

#### [LINE で「あるものごはん」を使ってみる](https://lin.ee/pLzG7zn)

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
