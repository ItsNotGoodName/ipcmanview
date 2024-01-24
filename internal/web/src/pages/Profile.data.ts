import { cache } from "@solidjs/router";
import { useClient } from "~/providers/client";

export const getProfilePage = cache(() => useClient().page.getProfilePage({}).then((req) => req.response), "getProfilePage")

export const getListMyGroups = cache(() => useClient().user.listMyGroups({}).then((req) => req.response), "listMyGroups")

export default function() {
  void getProfilePage()
  void getListMyGroups()
}
