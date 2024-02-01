import { cache } from "@solidjs/router"
import { useClient } from "~/providers/client"


export const getDeviceDetail = cache((id: bigint) => useClient().user.getDeviceDetail({ id }).then(res => res.response), "getDeviceDetail")
