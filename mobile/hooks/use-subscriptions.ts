import { useFocusEffect } from "expo-router"
import { useCallback, useState } from "react"

import { Api, SubscriptionsResponse } from "@/api"

export function useSubscriptions() {
  const [data, setData] = useState<SubscriptionsResponse | null>(null)
  const [isLoading, setIsLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  const fetchSubs = useCallback(async () => {
    try {
      setIsLoading(true)
      setError(null)
      const res = await Api.subscriptions.getAll()
      setData(res)
    } catch (e: any) {
      setError(e.message)
    } finally {
      setIsLoading(false)
    }
  }, [])

  useFocusEffect(
    useCallback(() => {
      fetchSubs()
    }, [fetchSubs])
  )

  return { data, isLoading, error, refetch: fetchSubs }
}
