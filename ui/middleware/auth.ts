import { useAuthStore } from "~/stores/auth"

export default defineNuxtRouteMiddleware(() => {
  // Don't run when we are hydrating from ssr
  const nuxtApp = useNuxtApp()
  if (process.client && nuxtApp.isHydrating && nuxtApp.payload.serverRendered) return

  const authStore = useAuthStore()
  if (!authStore.loggedIn) {
    return navigateTo('/login')
  }
})
