import { LinearGradient } from "expo-linear-gradient"
import React from "react"
import { Image, Pressable, StyleSheet, View } from "react-native"
import Animated, { FadeInDown } from "react-native-reanimated"

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
  index?: number;
}

export const SubscriptionRow: React.FC<SubscriptionRowProps> = ({
  title,
  price,
  currency,
  period,
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
      entering={FadeInDown.delay(80 + index * 60).duration(500).springify()}
    >
      <Pressable
        style={({ pressed }) => [
          styles.container,
          pressed && styles.pressed,
        ]}
      >
        <LinearGradient
          colors={
            isDark
              ? ["rgba(255,255,255,0.06)", "rgba(255,255,255,0.02)"]
              : ["rgba(0,0,0,0.03)", "rgba(0,0,0,0.01)"]
          }
          start={{ x: 0, y: 0 }}
          end={{ x: 1, y: 0 }}
          style={styles.gradient}
        >
          <View style={[styles.brandStripe, { backgroundColor: safeColor }]} />

          <View style={styles.content}>
            <View style={styles.leftSection}>
              <View style={[styles.iconWrapper, { backgroundColor: `${safeColor}15` }]}>
                <Image source={{ uri: iconUrl }} style={styles.icon} />
              </View>
              <ThemedText type="defaultSemiBold" style={styles.title} numberOfLines={1}>
                {title}
              </ThemedText>
            </View>

            <View style={styles.rightSection}>
              <View style={styles.priceRow}>
                <ThemedText style={styles.price}>
                  {price} {currency}
                </ThemedText>
                <ThemedText style={[styles.period, { color: isDark ? "rgba(255,255,255,0.35)" : "rgba(0,0,0,0.35)" }]}>
                  /{period}
                </ThemedText>
              </View>
              <ThemedText style={[styles.date, { color: isDark ? "rgba(255,255,255,0.3)" : "rgba(0,0,0,0.3)" }]}>
                {paymentDate}
              </ThemedText>
            </View>
          </View>
        </LinearGradient>
      </Pressable>
    </Animated.View>
  )
}

const styles = StyleSheet.create({
  container: {
    borderRadius: 18,
    overflow: "hidden",
    borderWidth: 1,
    borderColor: "rgba(255, 255, 255, 0.07)",
  },
  pressed: {
    opacity: 0.8,
    transform: [{ scale: 0.98 }],
  },
  gradient: {
    flexDirection: "row",
    height: 68,
    alignItems: "center",
  },
  brandStripe: {
    width: 3.5,
    height: "55%",
    borderRadius: 2,
    marginLeft: 5,
  },
  content: {
    flex: 1,
    flexDirection: "row",
    justifyContent: "space-between",
    alignItems: "center",
    paddingHorizontal: 14,
  },
  leftSection: {
    flexDirection: "row",
    alignItems: "center",
    gap: 12,
    flex: 1,
  },
  iconWrapper: {
    width: 40,
    height: 40,
    borderRadius: 12,
    justifyContent: "center",
    alignItems: "center",
    overflow: "hidden",
  },
  icon: {
    width: 28,
    height: 28,
    borderRadius: 6,
  },
  title: {
    fontSize: 16,
    letterSpacing: -0.2,
    flex: 1,
  },
  rightSection: {
    alignItems: "flex-end",
    marginLeft: 8,
  },
  priceRow: {
    flexDirection: "row",
    alignItems: "baseline",
  },
  price: {
    fontSize: 16,
    fontWeight: "700",
    letterSpacing: -0.3,
  },
  period: {
    fontSize: 12,
    fontWeight: "500",
  },
  date: {
    fontSize: 11,
    fontWeight: "500",
    marginTop: 2,
  },
})
