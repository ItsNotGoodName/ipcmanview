import { InvalidSessionError } from "~/core/client.gen";
import { useAuthStore } from "~/stores/auth";

export default defineNuxtPlugin(async () => {
  const cookieSession = useCookie("session", {
    path: "/",
    secure: true,
    sameSite: "strict",
    httpOnly: true,
  })
  if (!cookieSession.value) {
    return
  }

  const { $authService } = useNuxtApp()
  const authStore = useAuthStore()

  // Refresh token
  try {
    authStore.token = await $authService.refresh({ session: cookieSession.value }).then((res) => res.token)
  } catch (e) {
    if (e instanceof InvalidSessionError) {
      cookieSession.value = ""
    } else {
      throw e
    }
  }
})
