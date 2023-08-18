import { WebrpcError } from "core/client.gen"

export function useMutation(fn: () => Promise<unknown>) {
  const loading = ref(false)
  const error = ref<WebrpcError>()

  const mutate = () => {
    error.value = undefined
    loading.value = true
    return fn().catch((e) => {
      error.value = e
    }).finally(() => {
      loading.value = false
    })
  }

  return {
    loading, error, mutate
  }
}
