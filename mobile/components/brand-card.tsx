import { LinearGradient } from "expo-linear-gradient"
import React from "react"
import { Image, StyleSheet, Text, View } from "react-native"
import Animated, { FadeInRight } from "react-native-reanimated"

import { useColorScheme } from "@/hooks/use-color-scheme"

import { ThemedText } from "./themed-text"

interface BrandCardProps {
  title: string;
  price: number;
  currency: string;
  paymentDate: string;
  brandColor: string;
  iconUrl: string;
  index?: number;
}

export const BrandCard: React.FC<BrandCardProps> = ({
  title,
  price,
  currency,
  paymentDate,
  brandColor,
  iconUrl,
  index = 0,
}) => {
  const colorScheme = useColorScheme()
  const isDark = colorScheme === "dark"
  const safeColor = brandColor || "#B4B4B4"

  return (
    <Animated.View
      entering={FadeInRight.delay(100 + index * 120).duration(600).springify()}
      style={[styles.shadowWrapper, { shadowColor: safeColor }]}
    >
      <LinearGradient
        colors={
          isDark
            ? [`${safeColor}25`, `${safeColor}08`]
            : [`${safeColor}18`, `${safeColor}06`]
        }
        start={{ x: 0, y: 0 }}
        end={{ x: 1, y: 1 }}
        style={[
          styles.card,
          {
            borderColor: `${safeColor}40`,
          },
        ]}
      >
        <View style={styles.topRow}>
          <View style={[styles.iconWrapper, { backgroundColor: `${safeColor}20` }]}>
            <Image source={{ uri: iconUrl }} style={styles.icon} />
          </View>
          <View style={styles.titleContainer}>
            <ThemedText type="defaultSemiBold" style={styles.title} numberOfLines={1}>
              {title}
            </ThemedText>
          </View>
        </View>

        <View style={styles.bottomRow}>
          <View>
            <ThemedText style={styles.price}>
              {price} {currency}
            </ThemedText>
            <Text style={[styles.date, { color: isDark ? "rgba(255,255,255,0.45)" : "rgba(0,0,0,0.4)" }]}>
              {paymentDate}
            </Text>
          </View>
          <View style={[styles.dot, { backgroundColor: safeColor }]} />
        </View>
      </LinearGradient>
    </Animated.View>
  )
}

const styles = StyleSheet.create({
  shadowWrapper: {
    shadowOffset: { width: 0, height: 8 },
    shadowOpacity: 0.4,
    shadowRadius: 16,
    elevation: 10,
    marginHorizontal: 6,
  },
  card: {
    width: 170,
    height: 120,
    borderRadius: 22,
    padding: 16,
    overflow: "hidden",
    borderWidth: 1,
    justifyContent: "space-between",
  },
  topRow: {
    flexDirection: "row",
    alignItems: "center",
    gap: 10,
  },
  iconWrapper: {
    width: 32,
    height: 32,
    borderRadius: 10,
    justifyContent: "center",
    alignItems: "center",
    overflow: "hidden",
  },
  icon: {
    width: 24,
    height: 24,
    borderRadius: 6,
  },
  titleContainer: {
    flex: 1,
  },
  title: {
    fontSize: 14,
    letterSpacing: -0.2,
  },
  bottomRow: {
    flexDirection: "row",
    justifyContent: "space-between",
    alignItems: "flex-end",
  },
  price: {
    fontSize: 20,
    fontWeight: "800",
    letterSpacing: -0.5,
  },
  date: {
    fontSize: 11,
    marginTop: 1,
    fontWeight: "500",
  },
  dot: {
    width: 8,
    height: 8,
    borderRadius: 4,
    marginBottom: 4,
  },
})
