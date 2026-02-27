import React from "react"
import { SafeAreaView, ScrollView, StyleSheet, View } from "react-native"

import { BrandCard } from "@/components/brand-card"
import { SubscriptionRow } from "@/components/subscription-row"
import { ThemedText } from "@/components/themed-text"
import { ThemedView } from "@/components/themed-view"

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
            <BrandCard
              title="Netflix"
              price={12.99}
              currency="$"
              paymentDate="3 мар."
              brandColor="#E50914"
              iconUrl="https://logo.clearbit.com/netflix.com"
            />
            <BrandCard
              title="Spotify"
              price={10.99}
              currency="$"
              paymentDate="7 мар."
              brandColor="#1DB954"
              iconUrl="https://logo.clearbit.com/spotify.com"
            />
            <BrandCard
              title="Яндекс Плюс"
              price={299}
              currency="₽"
              paymentDate="12 мар."
              brandColor="#FC3F1D"
              iconUrl="https://logo.clearbit.com/yandex.ru"
            />
          </ScrollView>
          
          <ThemedText type="subtitle" style={[styles.sectionTitle, { marginTop: 32 }]}>
            Все подписки
          </ThemedText>

          <View style={styles.listContainer}>
            <SubscriptionRow
              title="Apple iCloud"
              price={0.99}
              currency="$"
              period="мес"
              paymentDate="15 мар."
              brandColor="#0A84FF"
              iconUrl="https://logo.clearbit.com/apple.com"
            />
            <SubscriptionRow
              title="Local Coffee"
              price={1490}
              currency="₽"
              period="мес"
              paymentDate="18 мар."
              brandColor="#8B5E3C"
              iconUrl="https://logo.clearbit.com/starbucks.com"
            />
            <SubscriptionRow
              title="LinkedIn Premium"
              price={29.99}
              currency="$"
              period="мес"
              paymentDate="20 мар."
              brandColor="#0A66C2"
              iconUrl="https://logo.clearbit.com/linkedin.com"
            />
            <SubscriptionRow
              title="YouTube Premium"
              price={12.99}
              currency="$"
              period="мес"
              paymentDate="22 мар."
              brandColor="#FF0000"
              iconUrl="https://logo.clearbit.com/youtube.com"
            />
            <SubscriptionRow
              title="Adobe CC"
              price={52.99}
              currency="$"
              period="мес"
              paymentDate="25 мар."
              brandColor="#FA0F00"
              iconUrl="https://logo.clearbit.com/adobe.com"
            />
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
