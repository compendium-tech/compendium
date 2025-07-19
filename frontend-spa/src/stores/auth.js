import { defineStore } from 'pinia';
import { authService } from '../api';
import router from '../router';

export const useAuthStore = defineStore('auth', {
  state: () => ({
    accessTokenExpiry: null,
    userEmail: null,
    isAuthenticated: false,
    isRefreshingToken: false,
  }),

  actions: {
    setSession(accessTokenExpiry, email) {
      this.accessTokenExpiry = accessTokenExpiry;
      this.userEmail = email;
      this.isAuthenticated = true;
      console.log('Session set.');
    },

    clearSession() {
      this.accessTokenExpiry = null;
      this.userEmail = null;
      this.isAuthenticated = false;
      this.isRefreshingToken = false;
      console.log('Session cleared.');
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
        this.clearSession();
        router.push('/signin');
        return false;
      } finally {
        this.isRefreshingToken = false;
      }
    },

    async checkAuthenticationOnLoad() {
      const now = Date.now();
      if (!this.accessTokenExpiry) {
        this.isAuthenticated = false;
        return false;
      }

      const accessTokenExpiresAt = new Date(this.accessTokenExpiry).getTime();

      if (accessTokenExpiresAt > now) {
        this.isAuthenticated = true;
        return true;
      } else {
        return await this.refreshAccessToken();
      }
    },
  },
  persist: {
    key: 'auth-state',
    storage: localStorage,
    paths: ['accessTokenExpiry', 'userEmail', 'isAuthenticated'],
  },
});
