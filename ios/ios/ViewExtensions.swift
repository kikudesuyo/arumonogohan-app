//
//  ViewExtensions.swift
//  ios
//
//  Created by 菊池裕夢 on 2025/08/11.
//

import SwiftUI
import UIKit

// Viewを拡張して、特定の角だけを丸める機能を追加
extension View {
    func cornerRadius(_ radius: CGFloat, corners: UIRectCorner) -> some View {
        clipShape(RoundedCorner(radius: radius, corners: corners))
    }
}

// 特定の角を丸めるためのShape
struct RoundedCorner: Shape {
    var radius: CGFloat = .infinity
    var corners: UIRectCorner = .allCorners

    func path(in rect: CGRect) -> Path {
        let path = UIBezierPath(roundedRect: rect, byRoundingCorners: corners, cornerRadii: CGSize(width: radius, height: radius))
        return Path(path.cgPath)
    }
}
