import { GoogleSignin, isSuccessResponse } from "@react-native-google-signin/google-signin"
import * as Google from "expo-auth-session/providers/google"
import * as WebBrowser from "expo-web-browser"
import { useCallback, useEffect, useState } from "react"
import { Alert, Platform } from "react-native"

import { Config } from "@/constants/config"

WebBrowser.maybeCompleteAuthSession()

export type AuthUser = {
  email?: string
  name?: string
}

export function useGoogleAuth() {
  const [isLoading, setIsLoading] = useState(false)
  const [jwtToken, setJwtToken] = useState<string | null>(null)
  const [user, setUser] = useState<AuthUser | null>(null)

  const [request, authResponse, promptAsync] = Google.useAuthRequest({
    webClientId: Config.google.webClientId,
    iosClientId: Config.google.iosClientId,
    responseType: Platform.OS === "web" ? "id_token" : "code",
    scopes: ["openid", "profile", "email"],
    shouldAutoExchangeCode: false,
  })

  useEffect(() => {
    if (authResponse?.type === "success") {
      const { code, id_token } = authResponse.params
      const tokenOrCode = id_token || code
      
      if (tokenOrCode) {
        setIsLoading(true)
        exchangeAuth(tokenOrCode, id_token ? "id_token" : "code")
          .catch(err => {
            console.warn("Backend exchange failed:", err.message)
          })
          .finally(() => setIsLoading(false))
      }
    }
  }, [authResponse])

  useEffect(() => {
    if (Platform.OS === "android") {
      GoogleSignin.configure({
        webClientId: Config.google.webClientId,
        offlineAccess: true,
        forceCodeForRefreshToken: true,
      })
    }
  }, [])

  const exchangeAuth = useCallback(async (value: string, type: "code" | "id_token") => {
    const url = `${Config.apiBaseUrl}/auth/google`
    console.log(`[DEBUG] Exchanging ${type} with backend...`)

    try {
      const response = await fetch(url, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({
          serverAuthCode: value,
        }),
      })

      const data = await response.json()
      if (!response.ok) throw new Error(data.error || "Server error")

      setJwtToken(data.token)
      if (data.user) setUser(data.user)
      return data
    } catch (e) {
      console.error("[DEBUG] Exchange failed:", e)
      throw e
    }
  }, [])

  const login = useCallback(async () => {
    try {
      setIsLoading(true)

      if (Platform.OS === "web") {
        await promptAsync()
        return
      }

      if (Platform.OS === "android") {
        await GoogleSignin.hasPlayServices()
        const signInResponse = await GoogleSignin.signIn()
        if (isSuccessResponse(signInResponse)) {
          const code = signInResponse.data.serverAuthCode
          if (!code) throw new Error("No serverAuthCode from Google")
          
          await exchangeAuth(code, "code")
          setUser({
            email: signInResponse.data.user.email,
            name: signInResponse.data.user.name ?? undefined,
          })
        }
      }
    } catch (error) {
      Alert.alert("Ошибка", error instanceof Error ? error.message : "Неизвестная ошибка")
    } finally {
      setIsLoading(false)
    }
  }, [promptAsync, exchangeAuth])

  const logout = useCallback(async () => {
    try {
      setIsLoading(true)
      if (Platform.OS !== "web") await GoogleSignin.signOut()
      setJwtToken(null)
      setUser(null)
    } finally {
      setIsLoading(false)
    }
  }, [])

  return { user, jwtToken, isLoading, login, logout }
}
