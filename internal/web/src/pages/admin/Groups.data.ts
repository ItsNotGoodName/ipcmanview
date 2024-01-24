import { cache } from "@solidjs/router"
import { parseOrder } from "~/lib/order"
import { useClient } from "~/providers/client"
import { ListGroupsReq } from "~/twirp/rpc"

export const getListGroups = cache((input: ListGroupsReq) => useClient().admin.listGroups(input).then((req) => req.response), "listGroups")

export default function({ params }: any) {
  void getListGroups({
    page: {
      page: Number(params.page) || 1,
      perPage: Number(params.perPage) || 10
    },
    sort: {
      field: params.sort || "",
      order: parseOrder(params.order),
    }
  })
}
