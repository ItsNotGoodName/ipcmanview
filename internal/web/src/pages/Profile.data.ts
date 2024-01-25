import { cache } from "@solidjs/router";
import { useClient } from "~/providers/client";

export const getProfilePage = cache(() => useClient().user.getProfilePage({}).then((req) => req.response), "getProfilePage")

export default function() {
  void getProfilePage()
}
