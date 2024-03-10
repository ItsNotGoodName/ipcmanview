import { useClient } from "~/providers/client";
import { getListEventRules } from "./data";
import { cache } from "@solidjs/router";

export const getAdminEventsPage = cache(() => useClient().admin.getAdminEventsPage({}).then((req) => req.response), "getAdminEventsPage")

export default function() {
  void getAdminEventsPage()
  void getListEventRules()
}
