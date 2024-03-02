import { cache } from "@solidjs/router";
import { dotDecode, parseOrder } from "~/lib/utils";
import { useClient } from "~/providers/client";
import { GetEventsPageReq } from "~/twirp/rpc";
import { getListDevices, getListEventFilters } from "./data";

export const getEventsPage = cache((input: GetEventsPageReq) => useClient().user.getEventsPage(input).then((req) => req.response), "getEventsPage")

export default function({ params }: any) {
  void getEventsPage({
    page: {
      page: Number(params.page) || 0,
      perPage: Number(params.perPage) || 0
    },
    sort: {
      field: params.sort || "",
      order: parseOrder(params.order)
    },
    filterDeviceIDs: dotDecode(params.device),
    filterCodes: params.code ? JSON.parse(params.code) : [],
    filterActions: params.action ? JSON.parse(params.action) : [],
  })
  void getListDevices()
  void getListEventFilters()
}

