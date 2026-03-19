import { BlurView } from "expo-blur"
import React from "react"
import { Image, StyleSheet, View } from "react-native"

import { useColorScheme } from "@/hooks/use-color-scheme"

import { ThemedText } from "./themed-text"

interface SubscriptionRowProps {
  title: string;
  price: number;
  currency: string;
  period: string;
  paymentDate: string;
  brandColor: string;
  iconUrl: string;
  category?: string;
}

export const SubscriptionRow: React.FC<SubscriptionRowProps> = ({
  title,
  price,
  currency,
  period,
  paymentDate,
  brandColor,
  iconUrl,
}) => {
  const colorScheme = useColorScheme()
  const isDark = colorScheme === "dark"

  return (
    <View style={styles.container}>
      <BlurView
        intensity={isDark ? 25 : 40}
        tint={isDark ? "dark" : "light"}
        style={styles.blurWrapper}
      >
        {/* Тот самый вертикальный индикатор бренда */}
        <View style={[styles.brandIndicator, { backgroundColor: brandColor }]} />

        <View style={styles.content}>
          <View style={styles.leftSection}>
            <View style={styles.iconWrapper}>
              <Image source={{ uri: iconUrl }} style={styles.icon} />
            </View>
            <ThemedText type="defaultSemiBold" style={styles.title}>
              {title}
            </ThemedText>
          </View>

          <View style={styles.rightSection}>
            <View style={styles.priceContainer}>
              <ThemedText type="defaultSemiBold">
                {price} {currency}
              </ThemedText>
              <ThemedText style={styles.period}>/{period}</ThemedText>
            </View>
            <ThemedText style={styles.date}>
              {paymentDate}
            </ThemedText>
          </View>
        </View>
      </BlurView>
    </View>
  )
}

const styles = StyleSheet.create({
  container: {
    marginVertical: 4,
    borderRadius: 18,
    overflow: "hidden",
    borderWidth: 1,
    borderColor: "rgba(255, 255, 255, 0.08)",
  },
  blurWrapper: {
    flexDirection: "row",
    height: 64,
    alignItems: "center",
  },
  brandIndicator: {
    width: 4,
    height: "60%", // Не на всю высоту, чтобы выглядело аккуратнее
    borderRadius: 2,
    marginLeft: 4,
  },
  content: {
    flex: 1,
    flexDirection: "row",
    justifyContent: "space-between",
    alignItems: "center",
    paddingHorizontal: 12,
  },
  leftSection: {
    flexDirection: "row",
    alignItems: "center",
    gap: 12,
  },
  iconWrapper: {
    width: 36,
    height: 36,
    borderRadius: 10,
    backgroundColor: "rgba(255, 255, 255, 0.1)",
    justifyContent: "center",
    alignItems: "center",
    overflow: "hidden",
  },
  icon: {
    width: "100%",
    height: "100%",
  },
  title: {
    fontSize: 16,
  },
  rightSection: {
    alignItems: "flex-end",
  },
  priceContainer: {
    flexDirection: "row",
    alignItems: "baseline",
  },
  period: {
    fontSize: 12,
    opacity: 0.6,
  },
  date: {
    fontSize: 11,
    opacity: 0.5,
    marginTop: 2,
  },
})
