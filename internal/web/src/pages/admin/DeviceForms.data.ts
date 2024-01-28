import { cache } from "@solidjs/router";
import { getListLocations } from "./data";
import { useClient } from "~/providers/client"

export function loadAdminDevicesCreate() {
  void getListLocations()
}

export function loadAdminDevicesIDUpdate({ params }: any) {
  void getDevice(params.id)
}

export const getDevice = cache((id: bigint) => useClient().admin.getDevice({ id: id }).then((req) => req.response), "getDevice")
