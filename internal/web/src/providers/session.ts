import { makePersisted } from "@solid-primitives/storage"
import { cache } from "@solidjs/router"
import { createStore } from "solid-js/store"

export type Session = {
  valid: boolean
  username: string
  admin: boolean
  user_id: number
  disabled: boolean
}

// HACK: this is to allow switching routes based on session
export const [sessionCache, setSessionCache] = makePersisted(createStore<Session>({ valid: false, username: "", admin: false, user_id: 0, disabled: false }), { name: "session" })

export const getSession = cache(() =>
  fetch("/v1/session", {
    credentials: "include",
    headers: [['Content-Type', 'application/json'], ['Accept', 'application/json']],
  }).then(async (resp) => {
    if (resp.ok || resp.status == 401) {
      return resp.json()
    }

    throw new Error(`Invalid status code ${resp.status}`)
  }).then((data: Session) => {
    setSessionCache(data)
    return data
  }), "session")
