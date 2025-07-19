import axios from 'axios';
import Cookies from 'js-cookie';
import { useAuthStore } from "./stores/auth.js"

const API_BASE_URL = 'http://localhost:8080/api/v1';

const apiClient = axios.create({
  baseURL: API_BASE_URL,
  withCredentials: true,
});

apiClient.interceptors.request.use(config => {
  const csrfToken = Cookies.get('csrfToken');
  if (csrfToken) {
    config.headers['X-CSRF-Token'] = csrfToken;
  }
  return config;
}, error => {
  return Promise.reject(error);
});

let isRefreshing = false;

apiClient.interceptors.response.use(
  response => response,
  async error => {
    const originalRequest = error.config;

    if (error.response?.status === 401 && !originalRequest._retry) {
      originalRequest._retry = true;

      if (isRefreshing) {
        return new Promise(function(resolve, reject) {
          failedQueue.push({ resolve, reject });
        })
        .then(() => apiClient(originalRequest))
        .catch(err => Promise.reject(err));
      }

      isRefreshing = true;

      return new Promise(async (resolve, reject) => {
        try {
          const _refreshResponse = await authService.refreshToken();
          isRefreshing = false;
          processQueue(null);
          resolve(apiClient(originalRequest));
        } catch (refreshError) {
          isRefreshing = false;
          processQueue(refreshError);
          console.error('Refresh token failed, logging out:', refreshError);
          alert('Session expired. Please sign in again.');
          window.location.href = '/signin';
          reject(refreshError);
        }
      });
    }

    return Promise.reject(error);
  }
);

export const authService = {
  signUp: (name, email, password) => {
    return apiClient.post('/users', { name, email, password });
  },
  verifyMfaSignUp: (email, otp) => {
    return apiClient.post('/sessions?flow=mfa', { email, otp });
  },

  signInPassword: (email, password) => {
    return apiClient.post('/sessions?flow=password', { email, password });
  },
  verifyMfaSignIn: (email, otp) => {
    return apiClient.post('/sessions?flow=mfa', { email, otp });
  },

  initiatePasswordReset: (email) => {
    return apiClient.put('/password?flow=init', { email });
  },
  confirmPasswordReset: (email, otp, newPassword) => {
    return apiClient.put('/password?flow=finish', { email, otp, password: newPassword });
  },

  refreshToken: () => {
    return apiClient.post('/sessions?flow=refresh');
  }
};


let failedQueue = [];

const processQueue = (error, tokenRefreshed = false) => {
  failedQueue.forEach(prom => {
    if (error) {
      prom.reject(error);
    } else {
      prom.resolve(tokenRefreshed);
    }
  });
  failedQueue = [];
};

apiClient.interceptors.response.use(
  response => response,
  async error => {
    const originalRequest = error.config;
    const authStore = useAuthStore(); // Get the store instance

    // Check if it's a 401 error and not a retry attempt for a token refresh that already failed
    if (error.response?.status === 401 && !originalRequest._retry) {
      originalRequest._retry = true; // Mark this request as retried once

      // If a refresh is already in progress, queue this request
      if (authStore.isRefreshingToken) {
        return new Promise(function(resolve, reject) {
          failedQueue.push({ resolve: () => resolve(apiClient(originalRequest)), reject });
        });
      }

      // If no refresh is in progress, attempt to refresh
      const refreshed = await authStore.refreshAccessToken();

      if (refreshed) {
        // If refresh was successful, retry all queued requests and the original request
        processQueue(null, true);
        return apiClient(originalRequest); // Retry the original failed request
      } else {
        // If refresh failed (e.g., refresh token expired), clear queue and reject all
        processQueue(error);
        authStore.clearSession(); // Ensure session is cleared if refresh failed
        // The refreshAccessToken action already handles the alert and redirect
        return Promise.reject(error); // Propagate the error for the original request
      }
    }

    return Promise.reject(error);
  }
);
