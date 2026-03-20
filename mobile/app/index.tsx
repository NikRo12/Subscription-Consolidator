import { LinearGradient } from "expo-linear-gradient"
import React, { useMemo } from "react"
import {
  ActivityIndicator,
  Pressable,
  SafeAreaView,
  ScrollView,
  StyleSheet,
  Text,
  View,
} from "react-native"
import Animated, { FadeIn, FadeInDown } from "react-native-reanimated"

import { BrandCard } from "@/components/brand-card"
import { LoginScreen } from "@/components/login-screen"
import { SubscriptionRow } from "@/components/subscription-row"
import { ThemedText } from "@/components/themed-text"
import { ThemedView } from "@/components/themed-view"
import { useGoogleAuth } from "@/hooks/use-google-auth"
import { useSubscriptions } from "@/hooks/use-subscriptions"

const formatPaymentDate = (isoDate: string) =>
  new Intl.DateTimeFormat("ru-RU", {
    day: "numeric",
    month: "short",
  }).format(new Date(isoDate))

export default function HomeScreen() {
  const { jwtToken, isLoading: authLoading, login, logout, user } = useGoogleAuth()
  const { data, isLoading: subsLoading, error } = useSubscriptions(jwtToken)

  const { sortedSubscriptions, upcomingSubscriptions } = useMemo(() => {
    if (!data) return { sortedSubscriptions: [], upcomingSubscriptions: [] }
    const sorted = [...data.items].sort(
      (left, right) =>
        new Date(left.next_payment_date).getTime() - new Date(right.next_payment_date).getTime()
    )
    return {
      sortedSubscriptions: sorted,
      upcomingSubscriptions: sorted.slice(0, 3),
    }
  }, [data])

  if (!jwtToken) {
    return <LoginScreen onLogin={login} isLoading={authLoading} />
  }

  if (subsLoading && !data) {
    return (
      <ThemedView style={[styles.container, styles.centered]}>
        <SafeAreaView style={styles.centered}>
          <Animated.View entering={FadeIn.duration(600)} style={styles.centered}>
            <ActivityIndicator size="large" color="#7C3AED" />
            <ThemedText style={styles.loadingText}>Загрузка подписок…</ThemedText>
          </Animated.View>
        </SafeAreaView>
      </ThemedView>
    )
  }

  if (error) {
    return (
      <ThemedView style={[styles.container, styles.centered]}>
        <SafeAreaView style={styles.centered}>
          <Animated.View entering={FadeIn.duration(500)} style={styles.centered}>
            <Text style={styles.errorEmoji}>😔</Text>
            <ThemedText style={styles.errorText}>
              {error === "Unauthorized" ? "Сессия истекла" : `Ошибка: ${error}`}
            </ThemedText>
            {error === "Unauthorized" ? (
              <Pressable
                style={({ pressed }) => [styles.retryButton, pressed && { opacity: 0.7 }]}
                onPress={logout}
              >
                <ThemedText style={styles.retryButtonText}>Выйти и войти снова</ThemedText>
              </Pressable>
            ) : (
              <ThemedText style={styles.errorHint}>Попробуйте позже</ThemedText>
            )}
          </Animated.View>
        </SafeAreaView>
      </ThemedView>
    )
  }

  return (
    <ThemedView style={styles.container}>
      <SafeAreaView style={{ flex: 1 }}>
        <Animated.View entering={FadeInDown.duration(500).springify()} style={styles.header}>
          <View style={styles.headerLeft}>
            <View style={styles.avatar}>
              <Text style={styles.avatarText}>
                {user?.name ? user.name[0].toUpperCase() : "?"}
              </Text>
            </View>
            <View>
              <ThemedText style={styles.greeting}>
                {user?.name ? `${user.name.split(" ")[0]}` : "Подписки"}
              </ThemedText>
              {user?.email && (
                <ThemedText style={styles.email} numberOfLines={1}>{user.email}</ThemedText>
              )}
            </View>
          </View>
          <Pressable
            style={({ pressed }) => [
              styles.logoutButton,
              pressed && styles.logoutButtonPressed,
            ]}
            onPress={logout}
          >
            <ThemedText style={styles.logoutText}>Выйти</ThemedText>
          </Pressable>
        </Animated.View>

        <ScrollView
          contentContainerStyle={styles.scrollContent}
          showsVerticalScrollIndicator={false}
        >
          {data?.monthly_spend && data.monthly_spend.length > 0 && (
            <Animated.View entering={FadeInDown.delay(100).duration(600).springify()}>
              <LinearGradient
                colors={["rgba(124, 58, 237, 0.15)", "rgba(99, 102, 241, 0.05)"]}
                start={{ x: 0, y: 0 }}
                end={{ x: 1, y: 1 }}
                style={styles.spendCard}
              >
                <ThemedText style={styles.spendLabel}>Траты в месяц</ThemedText>
                {data.monthly_spend.map((spend, idx) => (
                  <ThemedText key={idx} style={styles.spendAmount}>
                    {spend.amount.toFixed(2)}
                    <ThemedText style={styles.spendCurrency}> {spend.currency}</ThemedText>
                  </ThemedText>
                ))}
              </LinearGradient>
            </Animated.View>
          )}

          {upcomingSubscriptions.length > 0 && (
            <>
              <Animated.View entering={FadeInDown.delay(200).duration(500)}>
                <ThemedText type="subtitle" style={styles.sectionTitle}>
                  Ближайшие оплаты
                </ThemedText>
              </Animated.View>

              <ScrollView
                horizontal
                showsHorizontalScrollIndicator={false}
                contentContainerStyle={styles.horizontalScroll}
              >
                {upcomingSubscriptions.map((subscription, idx) => (
                  <BrandCard
                    key={subscription.id}
                    title={subscription.title}
                    price={subscription.price}
                    currency={subscription.currency}
                    paymentDate={formatPaymentDate(subscription.next_payment_date)}
                    brandColor={subscription.brand_color}
                    iconUrl={subscription.icon_url}
                    index={idx}
                  />
                ))}
              </ScrollView>
            </>
          )}

          {sortedSubscriptions.length > 0 && (
            <>
              <Animated.View entering={FadeInDown.delay(400).duration(500)}>
                <ThemedText type="subtitle" style={[styles.sectionTitle, { marginTop: 32 }]}>
                  Все подписки
                </ThemedText>
              </Animated.View>

              <View style={styles.listContainer}>
                {sortedSubscriptions.map((subscription, idx) => (
                  <SubscriptionRow
                    key={subscription.id}
                    title={subscription.title}
                    price={subscription.price}
                    currency={subscription.currency}
                    period={subscription.period}
                    paymentDate={formatPaymentDate(subscription.next_payment_date)}
                    brandColor={subscription.brand_color}
                    iconUrl={subscription.icon_url}
                    index={idx}
                  />
                ))}
              </View>
            </>
          )}

          {sortedSubscriptions.length === 0 && !subsLoading && (
            <Animated.View entering={FadeIn.delay(300).duration(800)} style={styles.emptyState}>
              <Text style={styles.emptyEmoji}>📭</Text>
              <ThemedText style={styles.emptyTitle}>Подписок пока нет</ThemedText>
              <ThemedText style={styles.emptySubtitle}>
                Синхронизируйте почту, чтобы найти подписки
              </ThemedText>
            </Animated.View>
          )}
        </ScrollView>
      </SafeAreaView>
    </ThemedView>
  )
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
  },
  centered: {
    flex: 1,
    justifyContent: "center",
    alignItems: "center",
  },
  header: {
    flexDirection: "row",
    justifyContent: "space-between",
    alignItems: "center",
    paddingHorizontal: 20,
    paddingTop: 12,
    paddingBottom: 12,
  },
  headerLeft: {
    flexDirection: "row",
    alignItems: "center",
    gap: 12,
    flex: 1,
  },
  avatar: {
    width: 42,
    height: 42,
    borderRadius: 14,
    backgroundColor: "rgba(124, 58, 237, 0.2)",
    justifyContent: "center",
    alignItems: "center",
  },
  avatarText: {
    fontSize: 18,
    fontWeight: "700",
    color: "#7C3AED",
  },
  greeting: {
    fontSize: 20,
    fontWeight: "700",
    letterSpacing: -0.3,
  },
  email: {
    fontSize: 12,
    opacity: 0.4,
    marginTop: 1,
    maxWidth: 180,
  },
  logoutButton: {
    paddingHorizontal: 16,
    paddingVertical: 8,
    borderRadius: 12,
    backgroundColor: "rgba(255, 59, 48, 0.1)",
    borderWidth: 1,
    borderColor: "rgba(255, 59, 48, 0.15)",
  },
  logoutButtonPressed: {
    opacity: 0.7,
    transform: [{ scale: 0.96 }],
  },
  logoutText: {
    color: "#FF3B30",
    fontSize: 14,
    fontWeight: "600",
  },
  scrollContent: {
    paddingTop: 8,
    paddingBottom: 40,
  },
  spendCard: {
    marginHorizontal: 20,
    marginBottom: 28,
    padding: 20,
    borderRadius: 20,
    borderWidth: 1,
    borderColor: "rgba(124, 58, 237, 0.15)",
  },
  spendLabel: {
    fontSize: 14,
    fontWeight: "600",
    opacity: 0.5,
    marginBottom: 6,
    textTransform: "uppercase",
    letterSpacing: 0.8,
  },
  spendAmount: {
    fontSize: 36,
    fontWeight: "800",
    letterSpacing: -1,
  },
  spendCurrency: {
    fontSize: 20,
    fontWeight: "600",
    opacity: 0.6,
  },
  sectionTitle: {
    paddingHorizontal: 20,
    marginBottom: 16,
    fontSize: 18,
    fontWeight: "700",
    letterSpacing: -0.2,
  },
  horizontalScroll: {
    paddingHorizontal: 14,
  },
  listContainer: {
    paddingHorizontal: 20,
    gap: 8,
  },
  loadingText: {
    marginTop: 16,
    opacity: 0.5,
    fontSize: 15,
  },
  errorEmoji: {
    fontSize: 48,
    marginBottom: 16,
  },
  errorText: {
    fontSize: 17,
    fontWeight: "600",
    marginBottom: 16,
    textAlign: "center",
    paddingHorizontal: 32,
  },
  errorHint: {
    opacity: 0.5,
  },
  retryButton: {
    paddingHorizontal: 24,
    paddingVertical: 12,
    borderRadius: 14,
    backgroundColor: "rgba(255, 59, 48, 0.1)",
    borderWidth: 1,
    borderColor: "rgba(255, 59, 48, 0.15)",
  },
  retryButtonText: {
    color: "#FF3B30",
    fontWeight: "600",
  },
  emptyState: {
    alignItems: "center",
    paddingTop: 60,
    paddingHorizontal: 32,
    gap: 8,
  },
  emptyEmoji: {
    fontSize: 56,
    marginBottom: 12,
  },
  emptyTitle: {
    fontSize: 22,
    fontWeight: "700",
    letterSpacing: -0.3,
  },
  emptySubtitle: {
    fontSize: 15,
    opacity: 0.45,
    textAlign: "center",
    lineHeight: 22,
  },
})
