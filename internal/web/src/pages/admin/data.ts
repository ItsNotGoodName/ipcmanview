import { cache } from "@solidjs/router"
import { useClient } from "~/providers/client"


export const getListLocations = cache(() => useClient().admin.listLocations({}).then(res => res.response.locations), "listLocations")

export const getListDeviceFeatures = cache(() => useClient().admin.listDeviceFeatures({}).then(res => res.response.features), "listDeviceFeatures")
