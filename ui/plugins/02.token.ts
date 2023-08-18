import { useAuthStore } from "~/stores/auth";

export default defineNuxtPlugin(async () => {
  const authStore = useAuthStore()

  const cookieToken = useCookie("token", {
    path: "/",
    secure: true,
    sameSite: "strict",
    httpOnly: false, // change to "true" if you want only server-side access
    maxAge: 604800,
  })

  // Load token from cookie into auth store
  await useAsyncData(() => authStore.load(cookieToken.value))

  // Sync auth store's token with cookie
  authStore.$subscribe((_, state) => {
    cookieToken.value = state.token
  })
})
