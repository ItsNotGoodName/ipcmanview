import { cache } from "@solidjs/router";
import { useClient } from "~/providers/client";

export const getProfile = cache(() => useClient().page.profile({}).then((req) => req.response), "getProfile")

export const getListMyGroups = cache(() => useClient().user.listMyGroups({}).then((req) => req.response), "getListGroup")

export default function() {
  void getProfile()
  void getListMyGroups()
}
