import { cache } from "@solidjs/router";
import { useClient } from "~/providers/client";

export const getHomePage = cache(() => useClient().page.getHomePage({}).then((req) => req.response), "getHomePage")

export default function() {
  void getHomePage()
}
