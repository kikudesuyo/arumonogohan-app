//
//  ContentView.swift
//  ios
//
//  Created by 菊池裕夢 on 2025/08/11.
//

import SwiftUI

struct ContentView: View {
    var body: some View {
        NavigationView {
            VStack(spacing: 20) {
                Spacer()
                
                Text("どちらの機能を使いますか？")
                    .font(.title2)
                    .fontWeight(.bold)
                
                // Card for "Suggest recipe from ingredients"
                NavigationLink(destination: RecipeFromIngredientsView()) {
                    FeatureCardView(
                        iconName: "frying.pan.fill",
                        title: "食材からレシピを提案",
                        description: "冷蔵庫にある食材から作れる料理を提案します。"
                    )
                }
                
                // Card for "Suggest instructions from recipe name"
                NavigationLink(destination: IngredientsFromRecipeView()) {
                    FeatureCardView(
                        iconName: "book.fill",
                        title: "料理名から作り方を提案",
                        description: "料理名を指定して、必要な材料と作り方を調べます。"
                    )
                }
                
                Spacer()
                Spacer()
            }
            .padding()
            .navigationTitle("あるものごはん")
            .background(Color(UIColor.systemGroupedBackground).edgesIgnoringSafeArea(.all))
        }
    }
}

struct FeatureCardView: View {
    let iconName: String
    let title: String
    let description: String
    
    var body: some View {
        HStack {
            Image(systemName: iconName)
                .font(.largeTitle)
                .foregroundColor(.accentColor)
                .frame(width: 60)
            
            VStack(alignment: .leading) {
                Text(title)
                    .font(.headline)
                    .fontWeight(.bold)
                Text(description)
                    .font(.subheadline)
                    .foregroundColor(.secondary)
            }
            Spacer()
        }
        .padding()
        .background(Color(UIColor.secondarySystemGroupedBackground))
        .cornerRadius(12)
        .shadow(color: Color.black.opacity(0.05), radius: 5, x: 0, y: 2)
    }
}

#Preview {
    ContentView()
}