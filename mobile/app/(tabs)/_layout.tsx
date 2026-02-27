import { BlurView } from "expo-blur"
import { Tabs } from "expo-router"
import React from "react"
import { StyleSheet } from "react-native"

import { HapticTab } from "@/components/haptic-tab"
import { IconSymbol } from "@/components/ui/icon-symbol"
import { useColorScheme } from "@/hooks/use-color-scheme"

export default function TabLayout() {
  const colorScheme = useColorScheme()
  const isDark = colorScheme === "dark"

  return (
    <Tabs
      screenOptions={{
        tabBarActiveTintColor: "#00E5FF",
        tabBarInactiveTintColor: "#8E8E93",
        headerShown: false,
        tabBarButton: HapticTab,
        tabBarShowLabel: false,

        tabBarStyle: styles.tabBar,
        tabBarItemStyle: styles.tabBarItem,

        tabBarBackground: () => (
          <BlurView
            intensity={80}
            tint={isDark ? "dark" : "light"}
            style={StyleSheet.absoluteFill}
          />
        ),
      }}>
      <Tabs.Screen
        name="index"
        options={{
          tabBarIcon: ({ color }) => <IconSymbol size={24} name="house.fill" color={color} />,
        }}
      />
      <Tabs.Screen
        name="stats"
        options={{
          tabBarIcon: ({ color }) => <IconSymbol size={24} name="chart.bar.fill" color={color} />,
        }}
      />
      <Tabs.Screen
        name="profile"
        options={{
          tabBarIcon: ({ color }) => <IconSymbol size={24} name="person.fill" color={color} />,
        }}
      />
    </Tabs>
  )
}

const styles = StyleSheet.create({
  tabBar: {
    position: "absolute",
    bottom: 25,
    left: 20,
    right: 20,
    height: 65,
    borderRadius: 35,
    backgroundColor: "rgba(255, 255, 255, 0.1)",
    display: "flex",
    alignItems: "center",
    justifyContent: "center",
    
    borderWidth: 1,
    borderColor: "rgba(255, 255, 255, 0.1)",
    
    elevation: 0,
    shadowColor: "#000",
    shadowOffset: { width: 0, height: 10 },
    shadowOpacity: 0.3,
    shadowRadius: 20,
    
    overflow: "hidden",
  },
  tabBarItem: {
    alignItems: "center",
    justifyContent: "center",
  },
})
