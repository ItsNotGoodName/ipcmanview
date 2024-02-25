import { cache } from "@solidjs/router"
import { useClient } from "~/providers/client"

let locations: string[]
export const getListLocations = cache(async () => {
  if (locations)
    return locations
  locations = await useClient().admin.listLocations({}).then(res => res.response.locations)
  return locations
}, "listLocations")

export const getListDeviceFeatures = cache(() => useClient().admin.listDeviceFeatures({}).then(res => res.response.features), "listDeviceFeatures")
export const getListEventRules = cache(() => useClient().admin.listEventRules({}).then(res => res.response.items), "listEventRules")
export const getGroup = cache((id: bigint) => useClient().admin.getGroup({ id }).then((req) => req.response), "getGroup")
export const getDevice = cache((id: bigint) => useClient().admin.getDevice({ id }).then((req) => req.response), "getDevice")
