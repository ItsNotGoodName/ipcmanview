import { useAuthStore } from "~/stores/auth";
import { AuthService, DahuaService, Fetch, UserService } from "~/core/client.gen";

export default defineNuxtPlugin(async () => {
  const authStore = useAuthStore()

  const authFetch: Fetch = (input, init) =>
    $fetch.raw(input, {
      ...init,
      headers: {
        ...init?.headers,
        "Authorization": `Bearer ${authStore.token}`
      },
    }).then(async (res) => {
      if (res.status == 401 && authStore.token != "") {
        console.log("No longer authenticated.");
        await authStore.logout()
      }

      return res
    })

  const runtimeConfig = useRuntimeConfig()

  return {
    provide: {
      // @ts-ignore -- typescript bug with wrapping $fetch -- https://github.com/unjs/nitro/issues/470
      authService: new AuthService(runtimeConfig.public.apiBase, authFetch),
      userService: new UserService(runtimeConfig.public.apiBase, authFetch),
      dahuaService: new DahuaService(runtimeConfig.public.apiBase, authFetch)
    }
  }
});
