import { createRouter, createWebHistory } from 'vue-router'
import LandingPage from "../components/landing/LandingPage.vue"

const routes = [
  {
    path: '/',
    name: 'Home',
    component: LandingPage
  },
]

const router = createRouter({
  history: createWebHistory('/'),
  routes
})

export default router
