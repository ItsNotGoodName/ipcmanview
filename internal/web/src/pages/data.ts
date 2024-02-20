import { cache } from "@solidjs/router"
import { useClient } from "~/providers/client"

export const getlistDevices = cache(() => useClient().user.listDevices({}).then(res => res.response.devices), "listDevices")
export const getDeviceRPCStatus = cache((id: bigint) => useClient().user.getDeviceRPCStatus({ id }).then(res => res.response), "getDeviceRPCStatus")
export const getDeviceDetail = cache((id: bigint) => useClient().user.getDeviceDetail({ id }).then(res => res.response), "getDeviceDetail")
export const getListDeviceStorage = cache((id: bigint) => useClient().user.listDeviceStorage({ id }).then(res => res.response.items), "listDeviceStorage")
export const getDeviceSoftwareVersion = cache((id: bigint) => useClient().user.getDeviceSoftwareVersion({ id }).then(res => res.response), "getDeviceSoftwareVersion")
export const getListDeviceLicenses = cache((id: bigint) => useClient().user.listDeviceLicenses({ id }).then(res => res.response.items), "listDeviceLicenses")
export const getlistEmailAlarmEvents = cache(() => useClient().user.listEmailAlarmEvents({}).then(res => res.response.alarmEvents), "listEmailAlarmEvents")
