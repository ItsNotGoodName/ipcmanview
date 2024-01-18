import { PartialMessage } from "@protobuf-ts/runtime";
import { Accessor, createSignal } from "solid-js";
import { Timestamp } from "~/twirp/google/protobuf/timestamp";

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

export function parseDate(value: PartialMessage<Timestamp> | undefined): Date {
  return Timestamp.toDate(Timestamp.create(value))
}

export function formatDate(value: Date): string {
  return value.toLocaleString()
}
