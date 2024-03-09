import { cache } from "@solidjs/router"
import { parseOrder } from "~/lib/utils"
import { PageProps } from "~/lib/utils"
import { useClient } from "~/providers/client"
import { GetAdminDevicesPageReq } from "~/twirp/rpc"

export const getAdminDevicesPage = cache((input: GetAdminDevicesPageReq) => useClient().admin.getAdminDevicesPage(input).then((req) => req.response), "getAdminDevicesPage")

export type AdminDevicesPageSearchParams = {
  page: string
  perPage: string
  sort: string
  order: string
  panel: string
  slug: string
}

export default function({ params }: PageProps<AdminDevicesPageSearchParams>) {
  void getAdminDevicesPage({
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

