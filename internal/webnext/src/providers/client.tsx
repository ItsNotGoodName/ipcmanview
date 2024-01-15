import { TwirpFetchTransport } from "@protobuf-ts/twirp-transport";
import { HelloWorldClient } from "~/twirp/rpc.client";
import {
  createContext,
  ParentComponent,
  useContext,
} from "solid-js";

function createStore(): ClientContextType {
  let transport = new TwirpFetchTransport({ baseUrl: "/twirp" });

  return {
    helloWorld: new HelloWorldClient(transport),
  };
}

type ClientContextType = {
  helloWorld: HelloWorldClient
};

const ClientContext = createContext<ClientContextType>();

type ClientContextProps = {};

export const ClientProvider: ParentComponent<ClientContextProps> = (props) => {
  const store = createStore()

  return (
    <ClientContext.Provider value={store}>
      {props.children}
    </ClientContext.Provider>)
};

export function _useClient(): ClientContextType {
  const result = useContext(ClientContext);
  if (!result) throw new Error("useClient must be used within a ClientProvider");
  return result;
}

// HACK: ignore provider because solid-router's preloadRoute is not within the ClientProvider for some reason
const store = createStore()

export function useClient(): ClientContextType {
  return store
}
