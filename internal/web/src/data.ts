import { cache } from "@solidjs/router";
import { useClient } from "~/providers/client";

export const getMe = cache(() => useClient().user.getMe({}).then((req) => req.response), "getMe")
