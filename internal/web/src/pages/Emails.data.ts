import { Params, cache } from "@solidjs/router";
import { parseOrder } from "~/lib/utils";
import { useClient } from "~/providers/client";
import { GetEmailsPageReq } from "~/twirp/rpc";
import { getlistDevices, getListEmailAlarmEvents } from "./data";

export function withEmailPageQuery(q: URLSearchParams, searchParams: Partial<Params>): URLSearchParams {
  if (searchParams.alarmEvent)
    q.set("alarmEvent", searchParams.alarmEvent)
  if (searchParams.device)
    q.set("device", searchParams.device)
  return q
}

export const getEmailsPage = cache((input: GetEmailsPageReq) => useClient().user.getEmailsPage(input).then((req) => req.response), "getEmailsPage")

export default function({ params }: any) {
  void getEmailsPage({
    page: {
      page: Number(params.page) || 0,
      perPage: Number(params.perPage) || 0
    },
    sort: {
      field: params.sort || "",
      order: parseOrder(params.order)
    },
    filterDeviceIDs: params.device ? params.device.split('.').map((v: any) => BigInt(v)) : [],
    filterAlarmEvents: params.alarmEvent ? JSON.parse(params.alarmEvent) : [],
  })
  void getlistDevices()
  void getListEmailAlarmEvents()
}
