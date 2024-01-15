import { Timestamp } from "~/twirp/google/protobuf/timestamp";
import { cache } from "@solidjs/router";
import { useClient } from "~/providers/client";

export const getHello = cache(() => useClient().helloWorld.hello({ subject: "World", currentTime: Timestamp.now() }).then((req) => req.response), "hello")

export function loadHello() {
  void getHello()
}


