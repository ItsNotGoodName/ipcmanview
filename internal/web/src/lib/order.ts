import { Order, Sort } from "~/twirp/rpc"

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

export function nextSort(sort?: Sort, field?: string): { field?: string, order: Order } {
  if (field == sort?.field) {
    const order = ((sort?.order ?? Order.ORDER_UNSPECIFIED) + 1) % 3

    if (order == Order.ORDER_UNSPECIFIED) {
      return { field: undefined, order: Order.ORDER_UNSPECIFIED }
    }

    return { field: field, order: order }
  }

  return { field: field, order: Order.DESC }
}
