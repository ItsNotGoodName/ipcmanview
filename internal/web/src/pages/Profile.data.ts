import { cache } from "@solidjs/router";
import { useClient } from "~/providers/client";

export const getProfilePage = cache(() => useClient().page.getProfilePage({}).then((req) => req.response), "getProfilePage")

export default function() {
  void getProfilePage()
}
