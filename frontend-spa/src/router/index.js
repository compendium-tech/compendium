import { createRouter, createWebHistory } from 'vue-router'
import LandingPage from "../components/landing/LandingPage.vue"
import SignUpPage from "../components/auth/SignUpPage.vue"

const routes = [
  {
    path: '/',
    name: 'Home',
    component: LandingPage
  },
  {
    path: '/auth/signup',
    name: 'Sign up',
    component: SignUpPage
  },
]

const router = createRouter({
  history: createWebHistory('/'),
  routes
})

export default router
