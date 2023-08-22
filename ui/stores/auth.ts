import { defineStore, acceptHMRUpdate } from 'pinia'

export const useAuthStore = defineStore({
  id: 'auth',
  state: () => ({
    user: {
      id: 0,
      email: "",
      username: "",
      created_at: ""
    },
    tokenValid: false,
    token: "",
  }),
  getters: {
    loggedIn: (state) => state.tokenValid

  },
  actions: {
    async load(token: string | null | undefined): Promise<null> {
      if (!token) {
        return Promise.resolve(null)
      }

      this.$patch({
        token,
      })

      const { $userService } = useNuxtApp()

      const res = await $userService.me()

      this.$patch({
        user: res.user,
        tokenValid: true,
      })

      return null
    },

    async logout(): Promise<null> {
      this.$reset()

      await $fetch('/api/token', { method: 'POST', body: '' })

      navigateTo('/login')

      return null
    },

    async login(usernameOrEmail: string, password: string): Promise<null> {
      const { $authService } = useNuxtApp()

      const res = await $authService.login({ usernameOrEmail, password })

      await $fetch('/api/token', { method: 'POST', body: res.token })

      this.$patch({
        token: res.token,
        user: res.user,
        tokenValid: true
      })

      navigateTo('/')

      return null
    },
  },
})

// make sure to pass the right store definition
if (import.meta.hot) {
  import.meta.hot.accept(acceptHMRUpdate(useAuthStore, import.meta.hot))
}
