import { apiClient, handleAxiosError } from "./base"

export interface Subscription {
  tier: Tier
  since: string
  till: string
}

export type Tier = "student" | "team" | "community"

export interface SubscriptionResponse {
  isActive: boolean
  subscription?: Subscription
}

interface SubscriptionService {
  getSubscription: () => Promise<SubscriptionResponse>
  cancelSubscription: () => Promise<void>
}

export const subscriptionService: SubscriptionService = {
  getSubscription: async () => {
    try {
      const response = await apiClient.get("/subscription")

      return response.data
    } catch (error) {
      return handleAxiosError(error)
    }
  },
  cancelSubscription: async () => {
    try {
      await apiClient.delete("/subscription")
    } catch (error) {
      return handleAxiosError(error)
    }
  },
}
