import { defineStore } from "pinia"
import { authService } from "../api"

export const useAuthStore = defineStore("auth", {
  state: () => ({
    isAuthenticated: false,
    accessTokenExpiry: null,
    email: null,
    isRefreshingToken: false,
  }),
  actions: {
    setSession(accessTokenExpiry: string, email: string) {
      this.isAuthenticated = true
      this.accessTokenExpiry = accessTokenExpiry
      this.email = email
    },

    clearSession() {
      this.isAuthenticated = false
      this.accessTokenExpiry = null
      this.email = null
    },

    async refresh() {
      if (this.isRefreshingToken) {
        return false
      }

      this.isRefreshingToken = true
      try {
        const response = await authService.refresh()
        this.setSession(
          response.data.accessTokenExpiry,
          this.email || response.data.email
        )
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
    paths: ["isAuthenticated", "accessTokenExpiry", "email"],
  },
})
