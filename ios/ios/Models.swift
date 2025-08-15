//
//  Models.swift
//  ios
//
//  Created by 菊池裕夢 on 2025/08/11.
//

import Foundation

// チャットメッセージのデータ構造
// レシピ情報も保持できるようにする
struct Message: Identifiable, Equatable {
    static func == (lhs: Message, rhs: Message) -> Bool {
        lhs.id == rhs.id
    }
    
    let id = UUID()
    let text: String
    let isFromUser: Bool
    let recipe: RecipeDetail?
    
    // テキストメッセージ用のイニシャライザ
    init(text: String, isFromUser: Bool) {
        self.text = text
        self.isFromUser = isFromUser
        self.recipe = nil
    }
    
    // レシピメッセージ用のイニシャライザ
    init(recipe: RecipeDetail, isFromUser: Bool) {
        self.text = recipe.title // メッセージテキストとしてはタイトルを流用
        self.isFromUser = isFromUser
        self.recipe = recipe
    }
}

// レシピ提案APIへの入力データ構造
struct RecipeInput: Codable {
    let menu_category: String
    let ingredients: String
}

// レシピ提案APIからの応答データ構造
struct RecipeResponse: Codable {
    let recipe: RecipeDetail
}

// レシピ詳細のデータ構造
struct RecipeDetail: Codable, Equatable {
    let title: String
    let ingredients: [String]
    let instructions: [String]
    let summary: String
}