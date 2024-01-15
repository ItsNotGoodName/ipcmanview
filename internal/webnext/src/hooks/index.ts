import { Accessor, createSignal } from "solid-js";

export function useLoading(fn: () => Promise<void>): [Accessor<boolean>, () => Promise<void>] {
  const [loading, setLoading] = createSignal(false)
  return [loading, () => {
    if (loading()) {
      return Promise.resolve()
    }
    setLoading(true)
    return fn().finally(() => setLoading(false))
  }]
}
