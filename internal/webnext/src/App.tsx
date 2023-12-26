import { TwirpFetchTransport } from "@protobuf-ts/twirp-transport";
import { createSignal } from "solid-js";
import { HelloWorldClient } from "./twirp/rpc.client";
import { Timestamp } from "./twirp/google/protobuf/timestamp";

function App() {
  let transport = new TwirpFetchTransport({ baseUrl: "/twirp" });
  let client = new HelloWorldClient(transport);

  const [text, setText] = createSignal("")

  client.hello({ subject: "World", currentTime: Timestamp.now() }).then((req) => {
    setText(req.response.text + "! " + Timestamp.toDate(Timestamp.create(req.response.currentTime)))
  })

  return (
    <>
      {text()}
    </>
  )
}

export default App
