export const Config = {
  apiBaseUrl: process.env.EXPO_PUBLIC_API_BASE_URL || "https://api.subtrack.study",
  google: {
    webClientId: process.env.EXPO_PUBLIC_GOOGLE_WEB_CLIENT_ID || "",
  },
} as const
