import { cache } from "@solidjs/router";
import { useClient } from "~/providers/client";

export const getDevicesPage = cache(() => useClient().user.getDevicesPage({}).then((req) => req.response), "getDevicesPage")

export default function({ params }: any) {
  void getDevicesPage()
}
