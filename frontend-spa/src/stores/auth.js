import { defineStore } from 'pinia';
import router from '../router';
import { authService } from '../api';

export const useAuthStore = defineStore('auth', {
  state: () => ({
    isAuthenticated: false,
    accessTokenExpiry: null,
    userEmail: null,
    isRefreshingToken: false,
  }),
  actions: {
    setSession(accessTokenExpiry, email) {
      this.isAuthenticated = true;
      this.accessTokenExpiry = accessTokenExpiry;
      this.userEmail = email;
    },
    clearSession() {
      this.isAuthenticated = false;
      this.accessTokenExpiry = null;
      this.userEmail = null;
    },
    async refreshAccessToken() {
      if (this.isRefreshingToken) {
        return false;
      }

      this.isRefreshingToken = true;
      try {
        const response = await authService.refreshToken();
        this.setSession(response.data.accessTokenExpiry, this.userEmail);
        return true;
      } catch (error) {
        const errorKind = error.response?.data?.errorKind;

        if (errorKind === 8) {
          this.clearSession();
        }
        return false;
      } finally {
        this.isRefreshingToken = false;
      }
    },
  },
  persist: {
    key: 'authState',
    storage: localStorage,
    paths: ['accessTokenExpiry', 'userEmail', 'isAuthenticated'],
  },
});
