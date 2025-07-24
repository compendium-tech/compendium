import { defineStore } from "pinia"
import "pinia-plugin-persistedstate"

export const useAuthStore = defineStore("auth", {
  state: () => ({
    isAuthenticated: false,
    accessTokenExpiresAt: "",
    isRefreshingToken: false,
  }),
  actions: {
    setSession(accessTokenExpiresAt: string) {
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
  },
  persist: {
    key: "authState",
    storage: localStorage,
    pick: ["isAuthenticated", "accessTokenExpiresAt"],
  },
})
