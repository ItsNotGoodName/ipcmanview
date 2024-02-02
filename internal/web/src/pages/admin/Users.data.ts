import { cache } from "@solidjs/router"
import { parseOrder } from "~/lib/utils"
import { PageProps } from "~/lib/utils"
import { useClient } from "~/providers/client"
import { GetAdminUsersPageReq } from "~/twirp/rpc"

export const getAdminUsersPage = cache((input: GetAdminUsersPageReq) => useClient().admin.getAdminUsersPage(input).then((req) => req.response), "getAdminUsersPage")

export type AdminUsersPageSearchParams = {
  page: string
  perPage: string
  sort: string
  order: string
}

export default function({ params }: PageProps<AdminUsersPageSearchParams>) {
  void getAdminUsersPage({
    page: {
      page: Number(params.page) || 0,
      perPage: Number(params.perPage) || 0,
    },
    sort: {
      field: params.sort || "",
      order: parseOrder(params.order),
    },
  })
}
