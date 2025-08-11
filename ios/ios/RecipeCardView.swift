//
//  RecipeCardView.swift
//  ios
//
//  Created by 菊池裕夢 on 2025/08/11.
//

import SwiftUI

struct RecipeCardView: View {
    let recipe: RecipeDetail

    var body: some View {
        VStack(alignment: .leading, spacing: 12) {
            // Recipe Title
            Text(recipe.title)
                .font(.title2)
                .fontWeight(.bold)

            // Summary
            if !recipe.summary.isEmpty {
                Text(recipe.summary)
                    .font(.subheadline)
                    .foregroundColor(.secondary)
            }
            
            Divider()

            // Ingredients
            VStack(alignment: .leading, spacing: 8) {
                Text("材料")
                    .font(.headline)
                    .fontWeight(.semibold)
                ForEach(recipe.ingredients, id: \.self) { ingredient in
                    HStack {
                        Image(systemName: "checkmark.circle.fill")
                            .foregroundColor(.accentColor)
                        Text(ingredient)
                    }
                }
            }

            Divider()

            // Instructions
            VStack(alignment: .leading, spacing: 8) {
                Text("作り方")
                    .font(.headline)
                    .fontWeight(.semibold)
                ForEach(recipe.instructions.indices, id: \.self) { index in
                    HStack(alignment: .top) {
                        Text("\(index + 1).")
                            .fontWeight(.bold)
                            .frame(width: 25, alignment: .leading)
                        Text(recipe.instructions[index])
                    }
                }
            }
        }
        .padding()
        .background(Color.white)
        .cornerRadius(15)
        .shadow(color: .black.opacity(0.1), radius: 10, x: 0, y: 5)
    }
}

#Preview {
    let sampleRecipe = RecipeDetail(
        title: "豚の生姜焼き",
        ingredients: ["豚ロース肉", "生姜", "醤油", "みりん", "酒", "砂糖"],
        instructions: [
            "豚肉に軽く塩コショウを振る。",
            "フライパンに油を熱し、豚肉を両面焼く。",
            "生姜、醤油、みりん、酒、砂糖を混ぜ合わせたタレを加え、全体に絡める。"
        ],
        summary: "ご飯が進む定番おかず！高タンパクで元気が出ます。"
    )
    
    return ScrollView {
        RecipeCardView(recipe: sampleRecipe)
            .padding()
    }
    .background(Color.gray.opacity(0.1))
}
