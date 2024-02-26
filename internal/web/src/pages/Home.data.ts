import { cache } from "@solidjs/router";
import { useClient } from "~/providers/client";
import { getListLatestFiles } from "./data";

export const getHomePage = cache(() => useClient().user.getHomePage({}).then((req) => req.response), "getHomePage")

export default function() {
  void getHomePage()
  void getListLatestFiles()
}
