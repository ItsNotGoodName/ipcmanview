import { cache } from "@solidjs/router";
import { useClient } from "~/providers/client";

export const getProfile = cache(() => useClient().page.profile({}).then((req) => req.response), "getProfile")

export const getListGroup = cache(() => useClient().user.listGroup({}).then((req) => req.response), "getListGroup")

export function loadProfile() {
  void getProfile()
  void getListGroup()
}
