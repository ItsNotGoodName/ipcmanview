import { cache } from "@solidjs/router"
import { useClient } from "~/providers/client"
import { ListDeviceFeaturesResp_Item } from "~/twirp/rpc"

let locations: string[]

export const getListLocations = cache(async () => {
  if (locations)
    return locations
  locations = await useClient().admin.listLocations({}).then(res => res.response.locations)
  return locations
}, "listLocations")


let features: ListDeviceFeaturesResp_Item[]

export const getListDeviceFeatures = cache(async () => {
  if (features)
    return features
  features = await useClient().admin.listDeviceFeatures({}).then(res => res.response.features)
  return features
}, "listDeviceFeatures")
