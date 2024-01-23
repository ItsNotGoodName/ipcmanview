import { cache } from "@solidjs/router";
import { useClient } from "~/providers/client";
import { ListGroupsReq } from "~/twirp/rpc";

export const getListGroups = cache((input: ListGroupsReq) => useClient().admin.listGroups(input).then((req) => req.response), "listGroups")

export function loadListGroups() {
  void getListGroups({ page: 1, perPage: 100 })
}
