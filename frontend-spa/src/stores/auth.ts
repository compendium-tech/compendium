import { defineStore } from "pinia"
import { authService } from "../api"

export const useAuthStore = defineStore("auth", {
  state: () => ({
    isAuthenticated: false,
    accessTokenExpiresAt: null,
    isRefreshingToken: false,
  }),
  actions: {
    setSession(accessTokenExpiresAt: Date) {
      this.isAuthenticated = true
      this.accessTokenExpiresAt = accessTokenExpiresAt
    },

    setIsRefreshingToken(isRefreshingToken: boolean) {
      this.isRefreshingToken = isRefreshingToken
    },

    clearSession() {
      this.isAuthenticated = false
      this.accessTokenExpiresAt = null
    },

    async refresh() {
      if (this.isRefreshingToken) {
        return false
      }

      this.isRefreshingToken = true
      try {
        const response = await authService.refresh()
        this.setSession(response.accessTokenExpiresAt)
        this.isRefreshingToken = false
        return true
      } catch (error) {
        if (error.response?.data?.errorKind === 8) this.clearSession()

        this.isRefreshingToken = false
        return false
      }
    },
  },
  persist: {
    key: "authState",
    storage: localStorage,
    pick: ["isAuthenticated", "accessTokenExpiresAt"],
  },
})
