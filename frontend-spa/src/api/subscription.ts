import { apiClient } from "./base"

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
    return (await apiClient.get("/subscription")).data
  },
  cancelSubscription: async () => {
    await apiClient.delete("/subscription")
  },
  getSubscriptionInvitationCode: async () => {
    return (await apiClient.get("/subscription/invitationCode")).data
  },
  updateSubscriptionInvitationCode: async () => {
    return (await apiClient.put("/subscription/invitationCode", {})).data
  },
  removeSubscriptionInvitationCode: async () => {
    return (await apiClient.delete("/subscription/invitationCode")).data
  },
  joinSubscription: async (invitationCode: string) => {
    return (
      await apiClient.post("/subscription/members/me", null, {
        params: { invitationCode },
      })
    ).data
  },
  removeSubscriptionMember: async (memberId: string) => {
    await apiClient.delete(`/subscription/members/${memberId}`)
  },
}
