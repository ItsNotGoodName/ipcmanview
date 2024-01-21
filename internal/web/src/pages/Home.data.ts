import { cache } from "@solidjs/router";
import { useClient } from "~/providers/client";

export const getHome = cache(() => useClient().page.home({}).then((req) => req.response), "getHome")

export function loadHome() {
  void getHome()
}
