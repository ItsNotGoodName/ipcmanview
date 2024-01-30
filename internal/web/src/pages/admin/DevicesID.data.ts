import { cache } from "@solidjs/router"
import { useClient } from "~/providers/client"

export const getAdminDevicesIDPage = cache((id: bigint) => useClient().admin.getAdminDevicesIDPage({ id }).then((req) => req.response), "getAdminDevicesIDPage")

export default function({ params }: any) {
  void getAdminDevicesIDPage(BigInt(params.id || 0))
}
