import { cache } from "@solidjs/router"
import { useClient } from "~/providers/client"

export const getAdminGroupsIDPage = cache((id: bigint) => useClient().admin.getAdminDevicesIDPage({ id }).then((req) => req.response), "getAdminGroupIDPage")

export default function({ params }: any) {
  void getAdminGroupsIDPage(BigInt(params.id || 0))
}
