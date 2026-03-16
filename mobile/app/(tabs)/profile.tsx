import { useMemo } from "react"
import { ActivityIndicator, Pressable, StyleSheet, View } from "react-native"

import ParallaxScrollView from "@/components/parallax-scroll-view"
import { ThemedText } from "@/components/themed-text"
import { IconSymbol } from "@/components/ui/icon-symbol"
import { useGoogleAuth } from "@/hooks/use-google-auth"

export default function ProfileScreen() {
  const { user, jwtToken, isLoading, login, logout } = useGoogleAuth()

  const userLabel = useMemo(() => {
    if (!user) {
      return "Не авторизован"
    }

    if (user.name && user.email) {
      return `${user.name} (${user.email})`
    }

    return user.email ?? user.name ?? "Авторизован"
  }, [user])

  return (
    <ParallaxScrollView
      headerBackgroundColor={{ light: "#D0D0D0", dark: "#353636" }}
      headerImage={
        <IconSymbol
          size={310}
          color="#808080"
          name="chevron.left.forwardslash.chevron.right"
          style={styles.headerImage}
        />
      }>
      <View style={styles.content}>
        <ThemedText type="subtitle">Google OAuth</ThemedText>
        <ThemedText>{userLabel}</ThemedText>
        {jwtToken ? (
          <ThemedText numberOfLines={1}>JWT получен</ThemedText>
        ) : (
          <ThemedText>JWT отсутствует</ThemedText>
        )}

        <Pressable
          style={[styles.button, isLoading ? styles.buttonDisabled : undefined]}
          onPress={login}
          disabled={isLoading}>
          {isLoading ? (
            <ActivityIndicator color="#FFFFFF" />
          ) : (
            <ThemedText style={styles.buttonText}>Войти через Google</ThemedText>
          )}
        </Pressable>

        {user && (
          <Pressable
            style={[styles.button, styles.secondaryButton, isLoading ? styles.buttonDisabled : undefined]}
            onPress={logout}
            disabled={isLoading}>
            <ThemedText style={styles.secondaryButtonText}>Выйти</ThemedText>
          </Pressable>
        )}
      </View>
    </ParallaxScrollView>
  )
}

const styles = StyleSheet.create({
  headerImage: {
    color: "#808080",
    bottom: -90,
    left: -35,
    position: "absolute",
  },
  content: {
    paddingHorizontal: 20,
    paddingTop: 8,
    gap: 12,
  },
  button: {
    height: 48,
    borderRadius: 12,
    alignItems: "center",
    justifyContent: "center",
    backgroundColor: "#0A84FF",
  },
  buttonDisabled: {
    opacity: 0.6,
  },
  buttonText: {
    color: "#FFFFFF",
    fontWeight: "700",
  },
  secondaryButton: {
    backgroundColor: "transparent",
    borderWidth: 1,
    borderColor: "#0A84FF",
  },
  secondaryButtonText: {
    color: "#0A84FF",
    fontWeight: "700",
  },
})
