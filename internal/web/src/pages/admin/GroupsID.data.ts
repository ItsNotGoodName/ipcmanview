import { cache } from "@solidjs/router"
import { useClient } from "~/providers/client"

export const getAdminGroupsIDPage = cache((id: string) => useClient().admin.getAdminGroupsIDPage({ id }).then((req) => req.response), "getAdminGroupIDPage")

export default function({ params }: any) {
  void getAdminGroupsIDPage(params.id || 0)
}
