import { cache } from "@solidjs/router"
import { useClient } from "~/providers/client"

export const getAdminDevicesIDPage = cache((id: string) => useClient().admin.getAdminDevicesIDPage({ id }).then((req) => req.response), "getAdminDevicesIDPage")

export default function({ params }: any) {
  void getAdminDevicesIDPage(params.id || 0)
}
