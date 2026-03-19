import { useRouter } from "expo-router"
import React, { useMemo } from "react"
import { ActivityIndicator, Pressable, SafeAreaView, ScrollView, StyleSheet, View } from "react-native"

import { BrandCard } from "@/components/brand-card"
import { SubscriptionRow } from "@/components/subscription-row"
import { ThemedText } from "@/components/themed-text"
import { ThemedView } from "@/components/themed-view"
import { useSubscriptions } from "@/hooks/use-subscriptions"

const formatPaymentDate = (isoDate: string) =>
  new Intl.DateTimeFormat("ru-RU", {
    day: "numeric",
    month: "short",
  }).format(new Date(isoDate))

export default function HomeScreen() {
  const { data, isLoading, error } = useSubscriptions()

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

  const router = useRouter()

  if (error) {
    return (
      <ThemedView style={[styles.container, { justifyContent: "center", alignItems: "center" }]}>
        <ThemedText style={{ marginBottom: 16 }}>
          {error === "Unauthorized" ? "Войдите в аккаунт, чтобы увидеть подписки" : `Ошибка: ${error}`}
        </ThemedText>
        <Pressable 
          style={{ backgroundColor: "#0A84FF", padding: 12, borderRadius: 8 }}
          onPress={() => router.push("/profile")}
        >
          <ThemedText style={{ color: "#fff", fontWeight: "600" }}>Войти в профиль</ThemedText>
        </Pressable>
      </ThemedView>
    )
  }

  if (isLoading) {
    return (
      <ThemedView style={[styles.container, { justifyContent: "center", alignItems: "center" }]}>
        <ActivityIndicator size="large" color="#E50914" />
        <ThemedText style={{ marginTop: 16 }}>Обновление данных...</ThemedText>
      </ThemedView>
    )
  }

  return (
    <ThemedView style={styles.container}>
      <SafeAreaView style={{ flex: 1 }}>
        <ScrollView 
          contentContainerStyle={styles.scrollContent}
          showsVerticalScrollIndicator={false}
        >
          {data?.monthly_spend && data.monthly_spend.length > 0 && (
            <View style={styles.spendContainer}>
              <ThemedText type="subtitle" style={[styles.sectionTitle, { marginBottom: 8 }]}>
                Траты в месяц
              </ThemedText>
              {data.monthly_spend.map((spend, idx) => (
                <ThemedText key={idx} style={styles.spendText}>
                  {spend.amount.toFixed(2)} {spend.currency}
                </ThemedText>
              ))}
            </View>
          )}

          <ThemedText type="subtitle" style={styles.sectionTitle}>
            Ближайшие оплаты
          </ThemedText>
          
          <ScrollView 
            horizontal 
            showsHorizontalScrollIndicator={false}
            contentContainerStyle={styles.horizontalScroll}
          >
            {upcomingSubscriptions.map((subscription) => (
              <BrandCard
                key={subscription.id}
                title={subscription.title}
                price={subscription.price}
                currency={subscription.currency}
                paymentDate={formatPaymentDate(subscription.next_payment_date)}
                brandColor={subscription.brand_color}
                iconUrl={subscription.icon_url}
              />
            ))}
          </ScrollView>
          
          <ThemedText type="subtitle" style={[styles.sectionTitle, { marginTop: 32 }]}>
            Все подписки
          </ThemedText>

          <View style={styles.listContainer}>
            {sortedSubscriptions.map((subscription) => (
              <SubscriptionRow
                key={subscription.id}
                title={subscription.title}
                price={subscription.price}
                currency={subscription.currency}
                period={subscription.period}
                paymentDate={formatPaymentDate(subscription.next_payment_date)}
                brandColor={subscription.brand_color}
                iconUrl={subscription.icon_url}
              />
            ))}
          </View>

        </ScrollView>
      </SafeAreaView>
    </ThemedView>
  )
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
  },
  scrollContent: {
    paddingTop: 40,
    paddingBottom: 120,
  },
  sectionTitle: {
    paddingHorizontal: 20,
    marginBottom: 16,
    fontSize: 20,
    fontWeight: "700",
  },
  horizontalScroll: {
    paddingHorizontal: 12,
  },
  listContainer: {
    paddingHorizontal: 20,
    gap: 12,
  },
  spendContainer: {
    paddingHorizontal: 20,
    marginBottom: 24,
  },
  spendText: {
    fontSize: 24,
    fontWeight: "600",
  },
})
