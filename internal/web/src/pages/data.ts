import { cache } from "@solidjs/router"
import { useClient } from "~/providers/client"

export const getDeviceStatus = cache((id: bigint) => useClient().user.getDeviceStatus({ id }).then(res => res.response), "getDeviceStatus")
export const getDeviceDetail = cache((id: bigint) => useClient().user.getDeviceDetail({ id }).then(res => res.response), "getDeviceDetail")
export const getListDeviceStorage = cache((id: bigint) => useClient().user.listDeviceStorage({ id }).then(res => res.response.items), "listDeviceStorage")
export const getDeviceSoftwareVersion = cache((id: bigint) => useClient().user.getDeviceSoftwareVersion({ id }).then(res => res.response), "getDeviceSoftwareVersion")
export const getListDeviceLicenses = cache((id: bigint) => useClient().user.listDeviceLicenses({ id }).then(res => res.response.items), "listDeviceLicenses")
