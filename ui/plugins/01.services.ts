import { useAuthStore } from "~/stores/auth";
import { AuthService, DahuaService, Fetch, UserService } from "~/core/client.gen";

export default defineNuxtPlugin(async () => {
  const authStore = useAuthStore()

  const authFetch: Fetch = async (input, init) => {
    const res = await $fetch.raw(input, {
      ...init,
      headers: {
        ...init?.headers,
        "Authorization": `Bearer ${authStore.token}`
      },
    })

    // Refresh token and try request again
    if (res.status == 401 && authStore.valid) {
      await authStore.refresh()

      return await $fetch.raw(input, {
        ...init,
        headers: {
          ...init?.headers,
          "Authorization": `Bearer ${authStore.token}`
        },
      })
    }

    return res
  }

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
