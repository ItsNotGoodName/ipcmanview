import { cache } from "@solidjs/router";
import { dotDecode, parseOrder } from "~/lib/utils";
import { useClient } from "~/providers/client";
import { GetFilesPageReq } from "~/twirp/rpc";
import { getFileMonthCount } from "./data";

export const getFilesPage = cache((input: GetFilesPageReq) => useClient().user.getFilesPage(input).then((req) => req.response), "getFilesPage")

export default function({ params }: any) {
  const filterDeviceIDs = dotDecode(params.device)
  void getFilesPage({
    page: {
      page: Number(params.page) || 0,
      perPage: Number(params.perPage) || 0
    },
    filterDeviceIDs: filterDeviceIDs,
    filterMonthID: params.month ?? "",
    order: parseOrder(params.order)
  })
  void getFileMonthCount(filterDeviceIDs)
}

