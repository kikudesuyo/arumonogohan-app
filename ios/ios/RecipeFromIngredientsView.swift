//
//  RecipeFromIngredientsView.swift
//  ios
//
//  Created by 菊池裕夢 on 2025/08/11.
//

import SwiftUI
import UIKit

struct RecipeFromIngredientsView: View {
    @State private var messages: [Message] = [
        Message(text: "お持ちの食材を教えてください。", isFromUser: false)
    ]
    @State private var ingredients: [String] = []
    @State private var newIngredientText: String = ""
    @State private var isWaitingForResponse = false

    var body: some View {
        VStack(spacing: 0) {
            // Message Area
            messageScrollView
            
            // Typing Indicator
            if isWaitingForResponse {
                typingIndicatorView
            }
            
            // Input Area
            VStack(spacing: 0) {
                ingredientTagsView
                inputBarView
            }
            .background(Color(UIColor.systemGray5))
        }
        .background(Color(UIColor.systemGray6))
        .navigationTitle("食材からレシピ提案")
        .navigationBarTitleDisplayMode(.inline)
    }

    // MARK: - Subviews

    private var messageScrollView: some View {
        ScrollViewReader { scrollViewProxy in
            ScrollView {
                VStack(spacing: 12) {
                    ForEach(messages) { message in
                        MessageView(message: message)
                            .id(message.id)
                    }
                }
                .padding(.horizontal)
                .padding(.top, 20)
            }
            .onChange(of: messages) { oldValue, newValue in
                if let lastMessage = newValue.last {
                    withAnimation {
                        scrollViewProxy.scrollTo(lastMessage.id, anchor: .bottom)
                    }
                }
            }
        }
    }
    
    private var typingIndicatorView: some View {
        HStack(spacing: 4) {
            Text("入力中")
                .font(.caption)
                .foregroundColor(.gray)
            DotView(delay: 0)
            DotView(delay: 0.2)
            DotView(delay: 0.4)
            Spacer()
        }
        .padding(.horizontal)
        .padding(.vertical, 8)
        .background(Color(UIColor.systemGray6))
    }
    
    private var ingredientTagsView: some View {
        ScrollView(.horizontal, showsIndicators: false) {
            HStack {
                ForEach(ingredients, id: \.self) { ingredient in
                    TagView(name: ingredient) {
                        ingredients.removeAll { $0 == ingredient }
                    }
                }
            }
            .padding(.horizontal)
            .padding(.vertical, 10)
        }
        .frame(height: ingredients.isEmpty ? 0 : 50)
        .background(Color(UIColor.systemGray5))
        .animation(.easeInOut, value: ingredients.isEmpty)
    }

    private var inputBarView: some View {
        HStack(spacing: 16) {
            TextField("食材を追加して改行", text: $newIngredientText)
                .onSubmit(addIngredient)
                .padding(12)
                .background(Color.white)
                .cornerRadius(20)
                .overlay(
                    RoundedRectangle(cornerRadius: 20)
                        .stroke(Color.gray.opacity(0.3), lineWidth: 1)
                )

            Button(action: sendMessage) {
                Image(systemName: "arrow.up")
                    .font(.headline.weight(.bold))
                    .foregroundColor(.white)
                    .padding(12)
                    .background(ingredients.isEmpty ? Color.gray : Color.blue)
                    .clipShape(Circle())
            }
            .disabled(ingredients.isEmpty)
        }
        .padding()
    }

    // MARK: - Helper Views

    struct MessageView: View {
        let message: Message
        
        var body: some View {
            HStack(alignment: .bottom, spacing: 8) {
                if message.isFromUser {
                    Spacer()
                    Text(message.text)
                        .padding(12)
                        .background(Color.blue)
                        .foregroundColor(.white)
                        .cornerRadius(20, corners: [.topLeft, .topRight, .bottomLeft])
                    Image(systemName: "person.fill")
                        .font(.system(size: 28))
                        .foregroundColor(Color(UIColor.systemGray3))
                        .clipShape(Circle())
                } else {
                    // Bot message can be a recipe card or simple text
                    Image(systemName: "sparkles")
                        .font(.system(size: 28))
                        .foregroundColor(.white)
                        .padding(4)
                        .background(Color.orange)
                        .clipShape(Circle())
                    
                    if let recipe = message.recipe {
                        RecipeCardView(recipe: recipe)
                    } else {
                        Text(message.text)
                            .padding(12)
                            .background(Color.white)
                            .foregroundColor(.black)
                            .cornerRadius(20, corners: [.topLeft, .topRight, .bottomRight])
                    }
                    Spacer()
                }
            }
        }
    }
    
    struct DotView: View {
        let delay: Double
        @State private var scale: CGFloat = 0.5

        var body: some View {
            Circle()
                .frame(width: 5, height: 5)
                .scaleEffect(scale)
                .foregroundColor(.gray)
                .onAppear {
                    withAnimation(Animation.easeInOut(duration: 0.6).repeatForever().delay(delay)) {
                        self.scale = 1
                    }
                }
        }
    }
    
    struct TagView: View {
        let name: String
        let onDelete: () -> Void
        
        var body: some View {
            HStack(spacing: 4) {
                Text(name)
                    .font(.footnote)
                    .padding(.leading, 8)
                Button(action: onDelete) {
                    Image(systemName: "xmark.circle.fill")
                }
            }
            .padding(.vertical, 5)
            .padding(.trailing, 8)
            .background(Color.blue.opacity(0.2))
            .cornerRadius(12)
        }
    }

    // MARK: - Functions
    
    func addIngredient() {
        let trimmed = newIngredientText.trimmingCharacters(in: .whitespacesAndNewlines)
        if !trimmed.isEmpty {
            ingredients.append(trimmed)
            newIngredientText = ""
        }
    }

    func sendMessage() {
        guard !ingredients.isEmpty else { return }
        
        let ingredientsText = ingredients.joined(separator: ", ")
        let userMessage = Message(text: ingredientsText, isFromUser: true)
        withAnimation(.spring()) {
            messages.append(userMessage)
        }
        
        let recipeInput = RecipeInput(menu_category: "指定なし", ingredients: ingredientsText)
        ingredients = []
        isWaitingForResponse = true

        guard let url = URL(string: "http://localhost:8081/suggest") else {
            isWaitingForResponse = false
            return
        }
        var request = URLRequest(url: url)
        request.httpMethod = "POST"
        request.addValue("application/json", forHTTPHeaderField: "Content-Type")
        
        do {
            request.httpBody = try JSONEncoder().encode(recipeInput)
        } catch {
            print("Failed to encode recipe input: \(error)")
            isWaitingForResponse = false
            return
        }

        URLSession.shared.dataTask(with: request) { data, response, error in
            DispatchQueue.main.async {
                isWaitingForResponse = false
                if let error = error {
                    print("Failed to send request: \(error)")
                    let botMessage = Message(text: "エラーが発生しました。もう一度お試しください。", isFromUser: false)
                    withAnimation(.spring()) {
                        messages.append(botMessage)
                    }
                    return
                }
                
                guard let data = data else { return }
                
                // Print the raw server response for debugging
                print("--- Server Response ---")
                print(String(data: data, encoding: .utf8) ?? "Unable to print data as UTF-8 string")
                print("-----------------------")

                do {
                    let recipeResponse = try JSONDecoder().decode(RecipeResponse.self, from: data)
                    let botMessage = Message(recipe: recipeResponse.recipe, isFromUser: false)
                    withAnimation(.spring()) {
                        messages.append(botMessage)
                    }
                } catch {
                    print("Failed to decode response: \(error)")
                    let botMessage = Message(text: "レシピの解析に失敗しました。", isFromUser: false)
                    withAnimation(.spring()) {
                        messages.append(botMessage)
                    }
                }
            }
        }.resume()
    }
}

#Preview {
    NavigationView {
        RecipeFromIngredientsView()
    }
}