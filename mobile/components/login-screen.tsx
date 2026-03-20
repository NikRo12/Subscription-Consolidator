import { LinearGradient } from "expo-linear-gradient"
import React from "react"
import {
  ActivityIndicator,
  Pressable,
  StyleSheet,
  Text,
  View,
} from "react-native"
import Animated, {
  FadeIn,
  FadeInDown,
  FadeInUp,
} from "react-native-reanimated"

interface LoginScreenProps {
  onLogin: () => void;
  isLoading: boolean;
}

export function LoginScreen({ onLogin, isLoading }: LoginScreenProps) {
  return (
    <LinearGradient
      colors={["#0F0C29", "#302B63", "#24243E"]}
      start={{ x: 0, y: 0 }}
      end={{ x: 1, y: 1 }}
      style={styles.container}
    >
      <Animated.View
        entering={FadeIn.delay(200).duration(1000)}
        style={styles.circle1}
      />
      <Animated.View
        entering={FadeIn.delay(400).duration(1000)}
        style={styles.circle2}
      />
      <Animated.View
        entering={FadeIn.delay(600).duration(1000)}
        style={styles.circle3}
      />

      <View style={styles.content}>
        <Animated.View
          entering={FadeInDown.delay(300).duration(800).springify()}
          style={styles.logoContainer}
        >
          <View style={styles.logoCircle}>
            <Text style={styles.logoEmoji}>💳</Text>
          </View>
        </Animated.View>

        <Animated.View
          entering={FadeInDown.delay(500).duration(800).springify()}
        >
          <Text style={styles.title}>SubTrack</Text>
          <Text style={styles.subtitle}>
            Все подписки в одном месте
          </Text>
        </Animated.View>

        <Animated.View
          entering={FadeInDown.delay(700).duration(800).springify()}
          style={styles.featuresContainer}
        >
          <View style={styles.featureRow}>
            <Text style={styles.featureIcon}>📊</Text>
            <Text style={styles.featureText}>Отслеживай траты</Text>
          </View>
          <View style={styles.featureRow}>
            <Text style={styles.featureIcon}>🔔</Text>
            <Text style={styles.featureText}>Не пропускай оплаты</Text>
          </View>
          <View style={styles.featureRow}>
            <Text style={styles.featureIcon}>📱</Text>
            <Text style={styles.featureText}>Синхронизация с почтой</Text>
          </View>
        </Animated.View>

        <Animated.View
          entering={FadeInUp.delay(900).duration(800).springify()}
          style={styles.buttonContainer}
        >
          <Pressable
            style={({ pressed }) => [
              styles.googleButton,
              pressed && styles.googleButtonPressed,
              isLoading && styles.googleButtonDisabled,
            ]}
            onPress={onLogin}
            disabled={isLoading}
          >
            {isLoading ? (
              <ActivityIndicator color="#000" size="small" />
            ) : (
              <>
                <Text style={styles.googleIcon}>G</Text>
                <Text style={styles.googleButtonText}>
                  Войти через Google
                </Text>
              </>
            )}
          </Pressable>
        </Animated.View>
      </View>
    </LinearGradient>
  )
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
    overflow: "hidden",
  },
  circle1: {
    position: "absolute",
    width: 300,
    height: 300,
    borderRadius: 150,
    backgroundColor: "rgba(99, 102, 241, 0.15)",
    top: -80,
    right: -60,
  },
  circle2: {
    position: "absolute",
    width: 200,
    height: 200,
    borderRadius: 100,
    backgroundColor: "rgba(139, 92, 246, 0.1)",
    bottom: 120,
    left: -80,
  },
  circle3: {
    position: "absolute",
    width: 150,
    height: 150,
    borderRadius: 75,
    backgroundColor: "rgba(59, 130, 246, 0.08)",
    top: "45%",
    right: -40,
  },
  content: {
    flex: 1,
    justifyContent: "center",
    alignItems: "center",
    paddingHorizontal: 32,
    paddingBottom: 60,
  },
  logoContainer: {
    marginBottom: 32,
  },
  logoCircle: {
    width: 96,
    height: 96,
    borderRadius: 32,
    backgroundColor: "rgba(255, 255, 255, 0.1)",
    borderWidth: 1,
    borderColor: "rgba(255, 255, 255, 0.2)",
    justifyContent: "center",
    alignItems: "center",
  },
  logoEmoji: {
    fontSize: 44,
  },
  title: {
    fontSize: 40,
    fontWeight: "800",
    color: "#FFFFFF",
    textAlign: "center",
    letterSpacing: -1,
  },
  subtitle: {
    fontSize: 17,
    color: "rgba(255, 255, 255, 0.6)",
    textAlign: "center",
    marginTop: 8,
    letterSpacing: 0.3,
  },
  featuresContainer: {
    marginTop: 48,
    gap: 16,
    width: "100%",
    paddingHorizontal: 16,
  },
  featureRow: {
    flexDirection: "row",
    alignItems: "center",
    gap: 14,
    backgroundColor: "rgba(255, 255, 255, 0.06)",
    paddingVertical: 14,
    paddingHorizontal: 20,
    borderRadius: 16,
    borderWidth: 1,
    borderColor: "rgba(255, 255, 255, 0.08)",
  },
  featureIcon: {
    fontSize: 22,
  },
  featureText: {
    color: "rgba(255, 255, 255, 0.85)",
    fontSize: 16,
    fontWeight: "500",
  },
  buttonContainer: {
    width: "100%",
    marginTop: 48,
    paddingHorizontal: 16,
  },
  googleButton: {
    flexDirection: "row",
    height: 56,
    borderRadius: 16,
    backgroundColor: "#FFFFFF",
    justifyContent: "center",
    alignItems: "center",
    gap: 12,
    shadowColor: "#000",
    shadowOffset: { width: 0, height: 4 },
    shadowOpacity: 0.3,
    shadowRadius: 12,
    elevation: 8,
  },
  googleButtonPressed: {
    opacity: 0.9,
    transform: [{ scale: 0.98 }],
  },
  googleButtonDisabled: {
    opacity: 0.7,
  },
  googleIcon: {
    fontSize: 22,
    fontWeight: "700",
    color: "#4285F4",
  },
  googleButtonText: {
    fontSize: 17,
    fontWeight: "600",
    color: "#1A1A1A",
  },
})
