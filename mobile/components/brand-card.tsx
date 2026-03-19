import { BlurView } from "expo-blur"
import React from "react"
import { Image, StyleSheet, Text, View } from "react-native"

import { useColorScheme } from "@/hooks/use-color-scheme"

import { ThemedText } from "./themed-text"

// Типы для пропсов компонента
interface BrandCardProps {
  title: string;
  price: number;
  currency: string;
  paymentDate: string; // Например, "25 окт."
  brandColor: string; // HEX-код цвета бренда
  iconUrl: string; // URL иконки
}

export const BrandCard: React.FC<BrandCardProps> = ({
  title,
  price,
  currency,
  paymentDate,
  brandColor,
  iconUrl,
}) => {
  const colorScheme = useColorScheme()
  const isDark = colorScheme === "dark"

  // Создаем полупрозрачные версии brandColor для фона и тени
  const backgroundColor = `${brandColor}15` // 15 - это примерно 8% прозрачности в HEX
  const shadowColor = `${brandColor}40` // 40 - это примерно 25% прозрачности в HEX

  return (
    <View style={[styles.shadowWrapper, { shadowColor }]}>
      <BlurView
        intensity={isDark ? 40 : 60} // Интенсивность размытия
        tint={isDark ? "dark" : "light"} // Тон размытия
        style={[
          styles.card,
          {
            borderColor: `${brandColor}80`, // Полупрозрачная обводка цвета бренда
            backgroundColor: backgroundColor, // Легкий оттенок бренда на фоне
          },
        ]}
      >
        <View style={styles.header}>
          <Image source={{ uri: iconUrl }} style={styles.icon} />
          <ThemedText type="defaultSemiBold" style={styles.title} numberOfLines={1}>
            {title}
          </ThemedText>
        </View>

        <View style={styles.footer}>
          <ThemedText type="subtitle" style={styles.price}>
            {price} {currency}
          </ThemedText>
          <Text style={[styles.date, { color: isDark ? "#ccc" : "#666" }]}>
            {paymentDate}
          </Text>
        </View>
      </BlurView>
    </View>
  )
}

const styles = StyleSheet.create({
  shadowWrapper: {
    // Настройки тени для iOS
    shadowOffset: { width: 0, height: 8 },
    shadowOpacity: 1,
    shadowRadius: 12,
    // Настройки для Android (elevation не поддерживает цветные тени, поэтому эффект будет слабее)
    elevation: 8,
    marginHorizontal: 8, // Отступ между карточками в карусели
  },
  card: {
    width: 160, // Фиксированная ширина для карусели
    height: 110, // Фиксированная высота
    borderRadius: 24,
    padding: 16,
    overflow: "hidden", // Обязательно для BlurView и borderRadius
    borderWidth: 1, // Тонкая грань "стекла"
    justifyContent: "space-between",
  },
  header: {
    flexDirection: "row",
    alignItems: "center",
    gap: 8,
  },
  icon: {
    width: 24,
    height: 24,
    borderRadius: 6,
  },
  title: {
    flex: 1,
    fontSize: 14,
  },
  footer: {
    gap: 2,
  },
  price: {
    fontSize: 18,
    fontWeight: "700",
  },
  date: {
    fontSize: 11,
  },
})
