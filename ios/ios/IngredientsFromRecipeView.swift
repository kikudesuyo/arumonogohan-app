//
//  IngredientsFromRecipeView.swift
//  ios
//
//  Created by 菊池裕夢 on 2025/08/11.
//

import SwiftUI

struct IngredientsFromRecipeView: View {
    var body: some View {
        VStack {
            Text("料理名から作り方提案")
                .font(.title)
                .padding()
            Text("この機能は現在開発中です。")
                .foregroundColor(.gray)
        }
        .navigationTitle("料理名から作り方提案")
        .navigationBarTitleDisplayMode(.inline)
    }
}

#Preview {
    NavigationView {
        IngredientsFromRecipeView()
    }
}
