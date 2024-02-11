import { cache } from "@solidjs/router";
import { useClient } from "~/providers/client";

export const getEmailsIDPage = cache((id: bigint) => useClient().user.getEmailsIDPage({ id }).then((req) => req.response), "getEmailsIDPage")

export default function({ params }: any) {
  void getEmailsIDPage(BigInt(params.id))
}
