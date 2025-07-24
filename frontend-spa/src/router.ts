import { createRouter, createWebHistory } from "vue-router"
import LandingPage from "./components/landing/LandingPage.vue"
import SignUpPage from "./components/auth/SignUpPage.vue"
import SignInPage from "./components/auth/SignInPage.vue"
import NotFoundPage from "./components/NotFoundPage.vue"
import DashboardPage from "./components/dashboard/DashboardPage.vue"
import { useAuthStore } from "./stores/auth.ts"

const routes = [
  {
    path: "/",
    name: "home",
    component: LandingPage,
    meta: {
      title: "Welcome",
      requiresAuth: false,
      redirectIfAuthenticated: false,
    },
  },
  {
    path: "/dashboard",
    name: "dashboard",
    component: DashboardPage,
    meta: {
      title: "Dashboard",
      requiresAuth: true,
    },
  },
  {
    path: "/auth/signup",
    name: "signUp",
    component: SignUpPage,
    meta: {
      title: "Sign up",
      requiresAuth: false,
      redirectIfAuthenticated: true,
    },
  },
  {
    path: "/auth/signin",
    name: "signIn",
    component: SignInPage,
    meta: {
      title: "Sign in",
      requiresAuth: false,
      redirectIfAuthenticated: true,
    },
  },
  {
    path: "/:pathMatch(.*)*",
    name: "notFound",
    component: NotFoundPage,
    meta: {
      title: "Not found",
      requiresAuth: false,
      redirectIfAuthenticated: false,
    },
  },
]

const router = createRouter({
  history: createWebHistory("/"),
  routes,
})

const DEFAULT_TITLE = "Compendium"

router.beforeEach((to, _from, next) => {
  const authStore = useAuthStore()
  const isAuthenticated = authStore.isAuthenticated

  document.title = to.meta.title || DEFAULT_TITLE

  if (to.meta.requiresAuth && !isAuthenticated) {
    next("/auth/signin")
  } else if (to.meta.redirectIfAuthenticated && isAuthenticated) {
    next("/dashboard")
  } else {
    next()
  }
})

export default router
