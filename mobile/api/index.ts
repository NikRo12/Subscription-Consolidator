import * as SecureStore from "expo-secure-store"
import { Platform } from "react-native"

import { Config } from "@/constants/config"

export const JWT_STORE_KEY = "subtrack_jwt"

const getBaseUrl = () => {
  let url = Config.apiBaseUrl || ""
  if (url && !url.startsWith("http")) {
    url = `http://${url}`
  }
  return url
}

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
      const url = `${getBaseUrl()}/auth/google`
      console.log(`[DEBUG] POST ${url}`)
      const response = await fetch(url, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ serverAuthCode }),
      })
      const data = await response.json().catch(() => ({}))
      if (!response.ok) throw new Error(data.error || "Server error")
      return data
    }
  },
  subscriptions: {
    getAll: async (category?: string): Promise<SubscriptionsResponse> => {
      const url = `${getBaseUrl()}/subscriptions${category ? `?category=${category}` : ""}`
      console.log(`[DEBUG] GET ${url}`)
      const token = await getToken()
      
      if (!token) throw new Error("Unauthorized")

      const response = await fetch(url, {
        method: "GET",
        headers: { 
          "Content-Type": "application/json",
          "Authorization": `Bearer ${token}`
        },
      })
      const data = await response.json().catch(() => ({}))
      if (!response.ok) throw new Error(data.error || "Server error")
      return data
    }
  }
}
