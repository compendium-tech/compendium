import { createRouter, createWebHistory } from 'vue-router'
import LandingPage from "../components/landing/LandingPage.vue"
import AuthFormPage from "../components/auth/AuthFormPage.vue"
import { useAuthStore } from '../stores/auth'

const routes = [
  {
    path: '/',
    name: 'Home',
    component: LandingPage,
    meta: { requiresAuth: false, redirectIfAuthenticated: false },
  },
  {
    path: '/auth/signup',
    name: 'signup',
    component: AuthFormPage,
    props: { mode: 'signup' },
    meta: { requiresAuth: false, redirectIfAuthenticated: true }
  },
  {
    path: '/auth/signin',
    name: 'signin',
    component: AuthFormPage,
    props: { mode: 'signin' },
    meta: { requiresAuth: false, redirectIfAuthenticated: true }
  },
]

const router = createRouter({
  history: createWebHistory('/'),
  routes
})

router.beforeEach((to, _from, next) => {
  const authStore = useAuthStore();
  const isAuthenticated = authStore.isAuthenticated;

  if (to.meta.requiresAuth && !isAuthenticated) {
    next('/signin');
  } else if (to.meta.redirectIfAuthenticated && isAuthenticated) {
    next('/dashboard');
  } else {
    next();
  }
});

export default router
