import * as SecureStore from "expo-secure-store"
import { Platform } from "react-native"

import { Config } from "@/constants/config"

export const JWT_STORE_KEY = "subtrack_jwt"

export async function getToken() {
  if (Platform.OS === "web") {
    if (typeof window !== "undefined") {
      return localStorage.getItem(JWT_STORE_KEY)
    }
    return null
  }
  return await SecureStore.getItemAsync(JWT_STORE_KEY)
}

export async function setToken(token: string | null) {
  if (Platform.OS === "web") {
    if (token) {
      localStorage.setItem(JWT_STORE_KEY, token)
    } else {
      localStorage.removeItem(JWT_STORE_KEY)
    }
    return
  }
  if (token) {
    await SecureStore.setItemAsync(JWT_STORE_KEY, token)
  } else {
    await SecureStore.deleteItemAsync(JWT_STORE_KEY)
  }
}

export type Subscription = {
  id: string
  title: string
  price: number
  currency: string
  period: string
  category: string
  next_payment_date: string
  icon_url: string
  brand_color: string
  is_active: boolean
  description: string
}

export type SubscriptionsResponse = {
  monthly_spend: { amount: number; currency: string }[]
  items: Subscription[]
}

export const Api = {
  auth: {
    google: async (serverAuthCode: string) => {
      const url = `${Config.apiBaseUrl}/auth/google`
      console.log(`[DEBUG] POST ${url}`)
      const response = await fetch(url, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ serverAuthCode }),
      })
      const data = await response.json()
      if (!response.ok) throw new Error(data.error || "Server error")
      return data
    }
  },
  subscriptions: {
    getAll: async (category?: string): Promise<SubscriptionsResponse> => {
      
      return new Promise((resolve) => {
        setTimeout(() => {
          resolve({
            monthly_spend: [
              { amount: 1789, currency: "₽" },
              { amount: 109.95, currency: "$" }
            ],
            items: [
              {
                id: "netflix",
                title: "Netflix",
                price: 12.99,
                currency: "$",
                period: "monthly",
                category: "entertainment",
                next_payment_date: "2026-03-03",
                brand_color: "#E50914",
                icon_url: "https://logo.clearbit.com/netflix.com",
                is_active: true,
                description: "Тариф 'Standard', списание через 3 дня"
              },
              {
                id: "spotify",
                title: "Spotify",
                price: 10.99,
                currency: "$",
                period: "monthly",
                category: "entertainment",
                next_payment_date: "2026-03-07",
                brand_color: "#1DB954",
                icon_url: "https://logo.clearbit.com/spotify.com",
                is_active: true,
                description: "Premium Individual"
              },
              {
                id: "yandex-plus",
                title: "Яндекс Плюс",
                price: 299,
                currency: "₽",
                period: "monthly",
                category: "entertainment",
                next_payment_date: "2026-03-12",
                brand_color: "#FC3F1D",
                icon_url: "https://logo.clearbit.com/yandex.ru",
                is_active: true,
                description: "Плюс Мульти"
              },
              {
                id: "icloud",
                title: "Apple iCloud",
                price: 0.99,
                currency: "$",
                period: "monthly",
                category: "clouds",
                next_payment_date: "2026-03-15",
                brand_color: "#0A84FF",
                icon_url: "https://logo.clearbit.com/apple.com",
                is_active: true,
                description: "50 GB Storage"
              },
              {
                id: "local-coffee",
                title: "Local Coffee",
                price: 1490,
                currency: "₽",
                period: "monthly",
                category: "food",
                next_payment_date: "2026-03-18",
                brand_color: "#8B5E3C",
                icon_url: "https://logo.clearbit.com/starbucks.com",
                is_active: true,
                description: "Кофейный абонемент"
              },
              {
                id: "linkedin-premium",
                title: "LinkedIn Premium",
                price: 29.99,
                currency: "$",
                period: "monthly",
                category: "work",
                next_payment_date: "2026-03-20",
                brand_color: "#0A66C2",
                icon_url: "https://logo.clearbit.com/linkedin.com",
                is_active: true,
                description: "Career Plan"
              },
              {
                id: "youtube-premium",
                title: "YouTube Premium",
                price: 12.99,
                currency: "$",
                period: "monthly",
                category: "entertainment",
                next_payment_date: "2026-03-22",
                brand_color: "#FF0000",
                icon_url: "https://logo.clearbit.com/youtube.com",
                is_active: true,
                description: "Individual Plan"
              },
              {
                id: "adobe-cc",
                title: "Adobe CC",
                price: 52.99,
                currency: "$",
                period: "monthly",
                category: "work",
                next_payment_date: "2026-03-25",
                brand_color: "#FA0F00",
                icon_url: "https://logo.clearbit.com/adobe.com",
                is_active: true,
                description: "All Apps"
              }
            ]
          })
        }, 1000)
      })
    }
  }
}
