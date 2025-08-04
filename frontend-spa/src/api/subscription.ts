import { apiClient, handleAxiosError } from "./base"

export interface Subscription {
  role: SubscriptionRole
  tier: Tier
  since: string
  till: string
  members?: SubscriptionMember[]
}

export type SubscriptionRole = "payer" | "member"
export type Tier = "student" | "team" | "community"

export interface SubscriptionResponse {
  isActive: boolean
  subscription?: Subscription
}

export interface InvitationCodeResponse {
  invitationCode?: string
}

export interface SubscriptionMember {
  userId: string
  name: string
  email?: string
  role?: SubscriptionRole
  isAccountActive?: boolean
}

interface SubscriptionService {
  getSubscription: () => Promise<SubscriptionResponse>
  cancelSubscription: () => Promise<void>
  getSubscriptionInvitationCode: () => Promise<InvitationCodeResponse>
  updateSubscriptionInvitationCode: () => Promise<InvitationCodeResponse>
  removeSubscriptionInvitationCode: () => Promise<InvitationCodeResponse>
  joinSubscription: (invitationCode: string) => Promise<SubscriptionResponse>
  removeSubscriptionMember: (memberId: string) => Promise<void>
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
  getSubscriptionInvitationCode: async () => {
    try {
      const response = await apiClient.get("/subscription/invitationCode")
      return response.data
    } catch (error) {
      return handleAxiosError(error)
    }
  },
  updateSubscriptionInvitationCode: async () => {
    try {
      const response = await apiClient.put("/subscription/invitationCode", {})
      return response.data
    } catch (error) {
      return handleAxiosError(error)
    }
  },
  removeSubscriptionInvitationCode: async () => {
    try {
      const response = await apiClient.delete("/subscription/invitationCode")
      return response.data
    } catch (error) {
      return handleAxiosError(error)
    }
  },
  joinSubscription: async (invitationCode: string) => {
    try {
      const response = await apiClient.post("/subscription/members/me", null, {
        params: { invitationCode },
      })
      return response.data
    } catch (error) {
      return handleAxiosError(error)
    }
  },
  removeSubscriptionMember: async (memberId: string) => {
    try {
      await apiClient.delete(`/subscription/members/${memberId}`)
    } catch (error) {
      return handleAxiosError(error)
    }
  },
}
