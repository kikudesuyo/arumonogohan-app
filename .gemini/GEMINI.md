# Arumonogohan App

## プロジェクト概要

このプロジェクトは、手元にある食材から作れるレシピを提案するアプリケーションです。LINE Bot と連携し、GCP Functions 上で Go 言語で実装されています。

## 技術スタック

- **バックエンド:** Go
- **フレームワーク:** (フレームワーク名)
- **AI:** Gemini
- **プラットフォーム:** GCP Cloud Functions, LINE Messaging API

## コーディング規約

### api/

- レイヤードアーキテクチャに基づく設計

- entity/ に配置
- DB やセッション管理は repository/ に配置
- ビジネスロジックは service/ に配置
- プレゼンテーション層は handler/ に配置

## ディレクトリ構成

- `api/`: API 関連のロジック
- `gcp_functions.go`: GCP Cloud Functions のエントリーポイント
- `ios/`: iOS アプリのソースコード
