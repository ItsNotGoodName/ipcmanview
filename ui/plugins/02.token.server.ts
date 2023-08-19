import { useAuthStore } from "~/stores/auth";

export default defineNuxtPlugin(async () => {
  const authStore = useAuthStore()

  const cookieToken = useCookie("token", {
    path: "/",
    secure: true,
    sameSite: "strict",
    httpOnly: true, // change to "true" if you want only server-side access
    maxAge: 604800,
  })

  // Load token from cookie into auth store
  await authStore.load(cookieToken.value)
})