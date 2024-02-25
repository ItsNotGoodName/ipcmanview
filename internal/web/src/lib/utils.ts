import { PartialMessage } from "@protobuf-ts/runtime";
import { Accessor, Resource, Setter, batch, createEffect, createSignal, onCleanup } from "solid-js";
import { Timestamp } from "~/twirp/google/protobuf/timestamp";
import { type ClassValue, clsx } from "clsx"
import { toast } from "~/ui/Toast";
import { RpcError } from "@protobuf-ts/runtime-rpc";
import { FieldValues, FormError, FormStore, PartialValues, reset } from "@modular-forms/solid";
import { Order, PagePaginationResult, Sort } from "~/twirp/rpc";
import { createStore } from "solid-js/store";
import { useSearchParams } from "@solidjs/router";

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

export type CreateRowSelectionReturn<T> = {
  rows: Array<{ id: T, checked: boolean } | undefined>
  indeterminate: Accessor<boolean>
  multiple: Accessor<boolean>
  selections: Accessor<Array<T>>
  set: (id: T, value: boolean) => void
  setAll: (value: boolean) => void
}

export function createRowSelection<T>(ids: Accessor<Array<T>>, disabled = (_index: number) => false): CreateRowSelectionReturn<T> {
  const calculateDisabledCount = () => {
    const range = ids()
    let i = 0
    for (let index = 0; index < range.length; index++)
      if (disabled(index)) i++
    return i
  }

  const [disabledCount, setDisabledCount] = createSignal(calculateDisabledCount())
  const [rows, setRows] = createStore<Array<{ id: T, checked: boolean }>>(ids().map(v => ({ id: v, checked: false })))
  createEffect(() => batch(() => {
    setDisabledCount(calculateDisabledCount())
    setRows((prev) => ids().map(v => ({ id: v, checked: prev.find(p => p.id == v)?.checked || false })))
  }))

  const selections = () => rows.filter(v => v.checked == true).map(v => v.id)

  return {
    rows,
    indeterminate: () => {
      const length = selections().length
      return length != 0 && length != (rows.length - disabledCount())
    },
    multiple: () => {
      const length = selections().length
      return length != 0 && length == (rows.length - disabledCount())
    },
    selections,
    set: (id, value) => {
      setRows(
        (todo, index) => todo.id === id && !disabled(index),
        "checked",
        value,
      );
    },
    setAll: (value) => {
      setRows(
        (_, index) => !disabled(index),
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

export function isTableRowClick(t: MouseEvent) {
  return t.target && (t.target as any).tagName == "TD"
}

type CreateValueModalReturn<T> = {
  open: Accessor<boolean>
  value: Accessor<T>
  close: () => void
  setValue: Setter<T>
}

export function createModal<T>(value: T): CreateValueModalReturn<T> {
  const [getOpen, setOpen] = createSignal(false)
  const [getValue, setValue] = createSignal(value)
  return {
    open: getOpen,
    value: getValue,
    close: () => setOpen(false),
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

export function encodeBigInts(q: bigint[]): string {
  return q.join('.')
}

export function decodeBigInts(q?: string): bigint[] {
  return q ? q.split('.').map((v: any) => BigInt(v)) : []
}

export function useHiddenScrollbar(): void {
  const html = document.getElementsByTagName("html")[0]
  if (html.style.getPropertyValue("scrollbar-width") == "none") return
  html.style.setProperty("scrollbar-width", "none")
  onCleanup(() => html.style.removeProperty("scrollbar-width"))
}
