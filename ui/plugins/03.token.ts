import { useAuthStore } from "~/stores/auth";

export default defineNuxtPlugin(() => {
  const authStore = useAuthStore()
  const jwtCookie = useCookie("jwt", {
    path: "/",
    secure: true,
    sameSite: "strict",
    httpOnly: false,
  })
  const { pause, resume, isActive } = useIntervalFn(() => authStore.refresh(), 5 * 60 * 1000) // 5 Minutes
  const sync = () => {
    if (authStore.valid) {
      if (!isActive.value) {
        resume()
      }
    } else {
      if (isActive.value) {
        pause()
      }
    }

    jwtCookie.value = authStore.token
  }
  sync()
  authStore.$subscribe(() => sync())
})
