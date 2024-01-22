import { makePersisted } from "@solid-primitives/storage"
import { cache } from "@solidjs/router"
import { createStore } from "solid-js/store"

export type Session = {
  valid: boolean
  username: string
  admin: boolean
}

export const [session, setSession] = makePersisted(createStore<Session>({ valid: false, username: "", admin: false }), { name: "session" })

export const getSession = cache(() => fetch("/v1/session", {
  credentials: "include",
  headers: [['Content-Type', 'application/json'], ['Accept', 'application/json']],
}).then((resp) => {
  if (resp.ok || resp.status == 401) {
    return resp.json()
  }

  throw new Error(`Invalid status code ${resp.status}`)
}).then((data: Session) => {
  setSession(data)
  return data
}), "session")
