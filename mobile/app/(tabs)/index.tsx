import React from "react"
import { SafeAreaView, ScrollView, StyleSheet, View } from "react-native"

import { BrandCard } from "@/components/brand-card"
import { SubscriptionRow } from "@/components/subscription-row"
import { ThemedText } from "@/components/themed-text"
import { ThemedView } from "@/components/themed-view"

type ServerSubscription = {
  id: string
  title: string
  price: number
  currency: string
  period: string
  nextPaymentAt: string
  brandColor: string
  iconUrl: string
}

// Хардкод ответа сервера.
// TODO: заменить на реальный API когда он будет готов.
const serverResponse: { subscriptions: ServerSubscription[] } = {
  subscriptions: [
    {
      id: "netflix",
      title: "Netflix",
      price: 12.99,
      currency: "$",
      period: "мес",
      nextPaymentAt: "2026-03-03",
      brandColor: "#E50914",
      iconUrl: "https://logo.clearbit.com/netflix.com",
    },
    {
      id: "spotify",
      title: "Spotify",
      price: 10.99,
      currency: "$",
      period: "мес",
      nextPaymentAt: "2026-03-07",
      brandColor: "#1DB954",
      iconUrl: "https://logo.clearbit.com/spotify.com",
    },
    {
      id: "yandex-plus",
      title: "Яндекс Плюс",
      price: 299,
      currency: "₽",
      period: "мес",
      nextPaymentAt: "2026-03-12",
      brandColor: "#FC3F1D",
      iconUrl: "https://logo.clearbit.com/yandex.ru",
    },
    {
      id: "icloud",
      title: "Apple iCloud",
      price: 0.99,
      currency: "$",
      period: "мес",
      nextPaymentAt: "2026-03-15",
      brandColor: "#0A84FF",
      iconUrl: "https://logo.clearbit.com/apple.com",
    },
    {
      id: "local-coffee",
      title: "Local Coffee",
      price: 1490,
      currency: "₽",
      period: "мес",
      nextPaymentAt: "2026-03-18",
      brandColor: "#8B5E3C",
      iconUrl: "https://logo.clearbit.com/starbucks.com",
    },
    {
      id: "linkedin-premium",
      title: "LinkedIn Premium",
      price: 29.99,
      currency: "$",
      period: "мес",
      nextPaymentAt: "2026-03-20",
      brandColor: "#0A66C2",
      iconUrl: "https://logo.clearbit.com/linkedin.com",
    },
    {
      id: "youtube-premium",
      title: "YouTube Premium",
      price: 12.99,
      currency: "$",
      period: "мес",
      nextPaymentAt: "2026-03-22",
      brandColor: "#FF0000",
      iconUrl: "https://logo.clearbit.com/youtube.com",
    },
    {
      id: "adobe-cc",
      title: "Adobe CC",
      price: 52.99,
      currency: "$",
      period: "мес",
      nextPaymentAt: "2026-03-25",
      brandColor: "#FA0F00",
      iconUrl: "https://logo.clearbit.com/adobe.com",
    },
  ],
}

const sortedSubscriptions = [...serverResponse.subscriptions].sort(
  (left, right) => new Date(left.nextPaymentAt).getTime() - new Date(right.nextPaymentAt).getTime(),
)

const upcomingSubscriptions = sortedSubscriptions.slice(0, 3)

const formatPaymentDate = (isoDate: string) =>
  new Intl.DateTimeFormat("ru-RU", {
    day: "numeric",
    month: "short",
  }).format(new Date(isoDate))

export default function HomeScreen() {
  return (
    <ThemedView style={styles.container}>
      <SafeAreaView style={{ flex: 1 }}>
        <ScrollView 
          contentContainerStyle={styles.scrollContent}
          showsVerticalScrollIndicator={false}
        >
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
                paymentDate={formatPaymentDate(subscription.nextPaymentAt)}
                brandColor={subscription.brandColor}
                iconUrl={subscription.iconUrl}
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
                paymentDate={formatPaymentDate(subscription.nextPaymentAt)}
                brandColor={subscription.brandColor}
                iconUrl={subscription.iconUrl}
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
})
