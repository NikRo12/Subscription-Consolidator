import { useCallback, useEffect, useState } from "react"

import { Api, SubscriptionsResponse } from "@/api"

export function useSubscriptions(jwtToken: string | null) {
  const [data, setData] = useState<SubscriptionsResponse | null>(null)
  const [isLoading, setIsLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)

  const fetchSubs = useCallback(async () => {
    if (!jwtToken) return
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
  }, [jwtToken])

  useEffect(() => {
    if (jwtToken) {
      fetchSubs()
    } else {
      setData(null)
      setError(null)
    }
  }, [jwtToken, fetchSubs])

  return { data, isLoading, error, refetch: fetchSubs }
}
