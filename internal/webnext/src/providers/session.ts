import { makePersisted } from "@solid-primitives/storage"
import { cache, createAsync } from "@solidjs/router"
import { createEffect, createSignal } from "solid-js"

// HACK: allow App.tsx access to the session
export const [session, setSession] = makePersisted(createSignal(false), { name: "session" })
export function useSession() {
  const session = createAsync(getSession)

  createEffect(() => {
    const value = session()
    if (value != undefined) {
      setSession(value)
    }
  })
}

export const getSession = cache(() => fetch("/v1/session", {
  credentials: "include",
  headers: [['Content-Type', 'application/json'], ['Accept', 'application/json']],
}).then((resp) => {
  if (resp.ok) {
    return true
  }
  if (resp.status == 401) {
    return false
  }
  throw new Error(`Invalid status code ${resp.status}`)
}), "session")

export default function() {
  void getSession()
}
