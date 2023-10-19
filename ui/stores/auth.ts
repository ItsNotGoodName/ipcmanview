import { FetchError } from 'ofetch'
import { defineStore, acceptHMRUpdate } from 'pinia'
import { LoginRequest } from 'server/api/login.post'
import { useAlertStore } from './alert'

export const useAuthStore = defineStore({
  id: 'auth',
  state: (): { token: string, refreshPromise?: Promise<string> } => ({
    token: "",
  }),
  getters: {
    valid: (state) => !!state.token,
  },
  actions: {
    async refresh(): Promise<null> {
      if (this.$state.refreshPromise) {
        return this.$state.refreshPromise.then(() => null)
      }
      this.$state.refreshPromise = $fetch('/api/refresh', { method: "POST" })

      try {
        const token = await this.$state.refreshPromise
        this.$patch({ token })
      } catch (e) {
        if (e instanceof FetchError && e.statusCode == 401) {
          this.$reset()

          navigateTo('/login')
        } else if (e instanceof Error) {
          useAlertStore().toast({ title: "Auth Refresh", message: e.message, type: 'error' })

          throw e
        } else {
          throw e
        }
      } finally {
        this.$state.refreshPromise = undefined

        return null
      }
    },

    async logout(): Promise<null> {
      await $fetch('/api/logout', { method: "POST" })

      this.$reset()

      navigateTo('/login')

      return null
    },

    async login(usernameOrEmail: string, password: string): Promise<null> {
      const { token } = await $fetch('/api/login', {
        method: 'POST', body: {
          usernameOrEmail,
          password
        } satisfies LoginRequest
      })

      this.$patch({ token })

      navigateTo('/')

      return null
    },
  },
})

// make sure to pass the right store definition
if (import.meta.hot) {
  import.meta.hot.accept(acceptHMRUpdate(useAuthStore, import.meta.hot))
}
