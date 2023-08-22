import { defineStore, acceptHMRUpdate } from 'pinia'

export type Alert = {
  id: number
  type: 'info' | 'error' | 'success',
  title: string,
  message: string,
  timeout: number,
}

export const useAlertStore = defineStore({
  id: 'alert',
  state: () => ({
    lastId: 0,
    alerts: [] as Alert[]
  }),
  actions: {
    toast({ title = 'Alert', message = "", type = 'info', timeout = 5 }: { title?: string, message?: string, type?: 'info' | 'error' | 'success', timeout?: number }) {
      this.$patch((state) => {
        state.lastId++
        state.alerts.push({ id: state.lastId, title, message, type, timeout })
      })
    },
    dismiss(id: number) {
      this.$patch((state) => {
        state.alerts = state.alerts.filter((r) => r.id != id)
      })
    }
  }
})

// make sure to pass the right store definition
if (import.meta.hot) {
  import.meta.hot.accept(acceptHMRUpdate(useAlertStore, import.meta.hot))
}
