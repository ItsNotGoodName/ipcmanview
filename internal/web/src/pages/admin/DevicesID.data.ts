import { cache } from "@solidjs/router"
import { PageProps } from "~/lib/utils"
import { useClient } from "~/providers/client"
import { GetAdminGroupIDPageReq, } from "~/twirp/rpc"

export const getAdminGroupsIDPage = cache((input: GetAdminGroupIDPageReq) => useClient().admin.getAdmind(input).then((req) => req.response), "getAdminGroupIDPage")

export type AdminGroupsIDPageSearchParams = {
  id: string
}

export default function({ params }: PageProps<AdminGroupsIDPageSearchParams>) {
  void getAdminGroupsIDPage({ id: BigInt(params.id || 0) })
}
