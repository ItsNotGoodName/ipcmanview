import { PartialMessage } from "@protobuf-ts/runtime";
import { Accessor, createEffect, createSignal } from "solid-js";
import { Timestamp } from "~/twirp/google/protobuf/timestamp";
import { type ClassValue, clsx } from "clsx"
import { toast } from "~/ui/Toast";
import { RpcError } from "@protobuf-ts/runtime-rpc";
import { FieldValues, FormError, FormStore, PartialValues, reset } from "@modular-forms/solid";
import { Order, Sort } from "~/twirp/rpc";
import { createStore } from "solid-js/store";

export function cn(...inputs: ClassValue[]) {
  return clsx(inputs)
}

export function createLoading(fn: () => Promise<void>): [Accessor<boolean>, () => Promise<void>] {
  const [loading, setLoading] = createSignal(false)
  return [loading, async () => {
    if (loading()) {
      return
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

export function catchAsToast(e: Error) {
  toast.error("Error", e.message)
}

export function throwAsFormError(e: unknown) {
  if (e instanceof RpcError)
    // @ts-ignore
    throw new FormError(e.message, e.meta ?? {})
  if (e instanceof Error)
    throw new FormError(e.message)
  throw new FormError("Unknown error has occured.")
}

export type PageProps<T> = {
  params: Partial<T>
}

export const defaultPerPageOptions = [10, 25, 50, 100]

export function parseOrder(s?: string): Order {
  if (s == "desc")
    return Order.DESC
  if (s == "asc")
    return Order.ASC
  return Order.ORDER_UNSPECIFIED
}

export function encodeOrder(o: Order): string {
  if (o == Order.DESC)
    return "desc"
  if (o == Order.ASC)
    return "asc"
  return ""
}

export function toggleSortField(sort?: Sort, field?: string): { field?: string, order: Order } {
  if (field == sort?.field) {
    const order = ((sort?.order ?? Order.ORDER_UNSPECIFIED) + 1) % 3

    if (order == Order.ORDER_UNSPECIFIED) {
      return { field: undefined, order: Order.ORDER_UNSPECIFIED }
    }

    return { field: field, order: order }
  }

  return { field: field, order: Order.DESC }
}

export type CreateRowSelectionReturn<T> = {
  rows: Array<{ id: T, checked: boolean }>
  selections: Accessor<Array<T>>
  multiple: Accessor<boolean>
  indeterminate: Accessor<boolean>
  check: (id: T, value: boolean) => void
  checkAll: (value: boolean) => void
}

export function createRowSelection<T>(ids: Accessor<Array<T>>): CreateRowSelectionReturn<T> {
  const [rows, setRows] = createStore<Array<{ id: T, checked: boolean }>>(ids().map(v => ({ id: v, checked: false })))
  createEffect(() =>
    setRows((prev) => ids().map(v => ({ id: v, checked: prev.find(p => p.id == v)?.checked || false }))))

  const selections = () => rows.filter(v => v.checked == true).map(v => v.id)

  return {
    rows,
    selections,
    multiple: () => {
      const length = selections().length
      return length != 0 && length == rows.length
    },
    indeterminate: () => {
      const length = selections().length
      return length != 0 && length != rows.length
    },
    check: (id, value) => {
      setRows(
        (todo) => todo.id === id,
        "checked",
        value,
      );
    },
    checkAll: (value) => {
      setRows(
        () => true,
        "checked",
        value,
      );
    }
  }
}

export function syncForm<TFieldValues extends FieldValues>(form: FormStore<TFieldValues, any>, data: PartialValues<TFieldValues> | undefined): boolean {
  if (form.submitted || form.dirty) {
    return false
  }

  reset(form, { initialValues: data })

  return false
}
