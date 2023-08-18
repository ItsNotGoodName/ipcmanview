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
    load(token: string | null | undefined): Promise<null> {
      if (!token) {
        return Promise.resolve(null)
      }
      this.$patch({
        token,
      })
      const { $userService } = useNuxtApp()
      return $userService.me().then((res) => {
        this.$patch({
          user: res.user,
          tokenValid: true,
        })
        return null
      })
    },

    logout() {
      this.$reset()
      navigateTo('/login')
    },

    login(usernameOrEmail: string, password: string): Promise<void> {
      const { $authService } = useNuxtApp()
      return $authService.login({ usernameOrEmail, password }).then((res) => {
        this.$patch({
          token: res.token,
          user: res.user,
          tokenValid: true
        })
        navigateTo('/')
      })
    },
  },
})

// make sure to pass the right store definition, `useAuth` in this case.
if (import.meta.hot) {
  import.meta.hot.accept(acceptHMRUpdate(useAuthStore, import.meta.hot))
}
