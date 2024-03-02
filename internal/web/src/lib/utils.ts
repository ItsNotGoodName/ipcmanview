import { PartialMessage } from "@protobuf-ts/runtime";
import { Accessor, Resource, Setter, batch, createEffect, createMemo, createSignal, onCleanup } from "solid-js";
import { Timestamp } from "~/twirp/google/protobuf/timestamp";
import { type ClassValue, clsx } from "clsx"
import { toast } from "~/ui/Toast";
import { RpcError } from "@protobuf-ts/runtime-rpc";
import { FieldStore, FieldValues, FormError, FormStore, PartialValues, reset, setValue } from "@modular-forms/solid";
import { Order, PagePaginationResult, Sort } from "~/twirp/rpc";
import { createStore } from "solid-js/store";
import { useSearchParams } from "@solidjs/router";
import { createDateNow, createTimeDifference } from "@solid-primitives/date";

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

type RowSelectionItem<T> = {
  id: T,
  checked: boolean,
  disabled: boolean
}

export type CreateRowSelectionReturn<T> = {
  items: Array<RowSelectionItem<T> | undefined>
  multiple: Accessor<boolean>
  all: Accessor<boolean>
  selections: Accessor<Array<T>>
  set: (id: T, value: boolean) => void
  setAll: (value: boolean) => void
}

export function createRowSelection<T>(ids: Accessor<Array<{ id: T, disabled?: boolean }>>): CreateRowSelectionReturn<T> {
  const [items, setItems] = createStore<Array<RowSelectionItem<T>>>(
    ids().map(v => ({ id: v.id, checked: false, disabled: v.disabled || false }))
  )
  createEffect(() =>
    setItems((prev) => ids().map(v => ({ id: v.id, disabled: v.disabled || false, checked: prev.find(p => p.id == v.id)?.checked || false })))
  )

  return {
    items,
    multiple: () => {
      for (let index = 0; index < items.length; index++) {
        if (items[index].checked) return true
      }
      return false
    },
    all: () => {
      let disabled = 0
      for (let index = 0; index < items.length; index++) {
        if (items[index].disabled) disabled++
        else if (!items[index].checked) return false
      }
      if (items.length - disabled == 0) return false
      return true
    },
    selections: () => items.filter(v => v.checked == true).map(v => v.id),
    set: (id, value) => {
      setItems(
        (v) => v.id === id && !v.disabled,
        "checked",
        value,
      );
    },
    setAll: (value) => {
      setItems(
        (v) => !v.disabled,
        "checked",
        value,
      );
    }
  }
}

export function syncForm<TFieldValues extends FieldValues>(form: FormStore<TFieldValues, any>, data: Resource<PartialValues<TFieldValues> | undefined>): Accessor<boolean> {
  createEffect(() => {
    if (!data.loading && !data.error) {
      reset(form, { initialValues: data() })
    }
  })
  return () => data.loading || !!data.error
}

export function isTableDataClick(t: MouseEvent) {
  return t.target && (t.target as any).tagName == "TD"
}

type CreateValueModalReturn<T> = {
  open: Accessor<boolean>
  value: Accessor<T>
  setClose: () => void
  setValue: Setter<T>
}

export function createModal<T>(value: T): CreateValueModalReturn<T> {
  const [getOpen, setOpen] = createSignal(false)
  const [getValue, setValue] = createSignal(value)
  return {
    open: getOpen,
    value: getValue,
    setClose: () => setOpen(false),
    setValue: (...args) => batch(() => {
      setOpen(true)
      // @ts-ignore
      return setValue(...args)
    })
  }
}

export type CreatePagePaginationReturn = {
  previousPageDisabled: Accessor<boolean>
  previousPage: () => void
  nextPageDisabled: Accessor<boolean>
  nextPage: () => void
  setPerPage: (value: number) => void
}

export function createPagePagination(pageResult: () => PagePaginationResult | undefined): CreatePagePaginationReturn {
  const [_, setSearchParams] = useSearchParams()
  return {
    previousPageDisabled: () => pageResult()?.previousPage == pageResult()?.page,
    previousPage: () => setSearchParams({ page: pageResult()?.previousPage.toString() }),
    nextPageDisabled: () => pageResult()?.nextPage == pageResult()?.page,
    nextPage: () => setSearchParams({ page: pageResult()?.nextPage.toString() }),
    setPerPage: (value: number) => value && setSearchParams({ page: 1, perPage: value })
  }
}

function toggleSortField(sort?: Sort, field?: string): { field?: string, order: Order } {
  if (field == sort?.field) {
    const order = ((sort?.order ?? Order.ORDER_UNSPECIFIED) + 1) % 3

    if (order == Order.ORDER_UNSPECIFIED) {
      return { field: undefined, order: Order.ORDER_UNSPECIFIED }
    }

    return { field: field, order: order }
  }

  return { field: field, order: Order.DESC }
}

export function createToggleSortField(sort: () => Sort | undefined) {
  const [_, setSearchParams] = useSearchParams()
  return (field: string) => {
    const s = toggleSortField(sort(), field)
    return setSearchParams({ sort: s.field, order: encodeOrder(s.order) })
  }
}

export function relativeWSURL(uri: string): string {
  return `${window.location.protocol === "https:" ? "wss:" : "ws:"}//${window.location.host}${uri}`
}

export function encodeQuery(q: URLSearchParams): string {
  if (q.size == 0)
    return ""
  return "?" + q.toString()
}

export function dotEncode(q: string[]): string {
  return q.join('.')
}

export function dotDecode(q?: string): string[] {
  return q ? q.split('.').map((v: any) => v) : []
}

export function useHiddenScrollbar(): void {
  const html = document.getElementsByTagName("html")[0]
  if (html.style.getPropertyValue("scrollbar-width") == "none") return
  html.style.setProperty("scrollbar-width", "none")
  onCleanup(() => html.style.removeProperty("scrollbar-width"))
}

export function validationState(error?: string | boolean): "invalid" | "valid" {
  return error ? "invalid" : "valid"
}

export function setFormValue(form: FormStore<any, any>, field: FieldStore<any, any>) {
  return (value: any) => setValue(form, field.name, value)
}

export function createUptime(date: Accessor<Date>) {
  const [now, update] = createDateNow(() => false);
  const [difference] = createTimeDifference(date, now)
  const timer = setInterval(update, 1000)
  onCleanup(() => clearInterval(timer))

  return createMemo(() => {
    const total = difference() / 1000
    const days = Math.floor(total / 86400)
    const hours = Math.floor((total % 86400) / 3600)
    const minutes = Math.floor((total % 3600) / 60)
    const seconds = Math.floor(total % 60)
    return {
      days,
      hasDays: days > 0,
      hours,
      hasHours: hours > 0 || days > 0,
      minutes,
      hasMinutes: minutes > 0 || hours > 0 || days > 0,
      seconds,
    }
  })
}
