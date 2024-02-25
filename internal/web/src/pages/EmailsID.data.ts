import { cache } from "@solidjs/router";
import { decodeBigInts } from "~/lib/utils";
import { useClient } from "~/providers/client";
import { GetEmailsIDPageReq } from "~/twirp/rpc";

export const getEmailsIDPage = cache((input: GetEmailsIDPageReq) => useClient().user.getEmailsIDPage(input).then((req) => req.response), "getEmailsIDPage")

export default function({ params }: any) {
  void getEmailsIDPage({
    id: BigInt(params.id ?? 0),
    filterAlarmEvents: params.alarmEvent ? JSON.parse(params.alarmEvent) : [],
    filterDeviceIDs: decodeBigInts(params.device),
  })
}
