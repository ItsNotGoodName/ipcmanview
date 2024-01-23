import { TwirpFetchTransport } from "@protobuf-ts/twirp-transport";
import { AdminClient, AuthClient, HelloWorldClient, PageClient, UserClient } from "~/twirp/rpc.client";
import {
  createContext,
  ParentComponent,
  useContext,
} from "solid-js";
import { MethodInfo, NextUnaryFn, RpcError, RpcOptions, UnaryCall } from "@protobuf-ts/runtime-rpc";
import { revalidate } from "@solidjs/router";
import { getSession } from "./session";

function createStore(): ClientContextType {
  let transport = new TwirpFetchTransport({
    baseUrl: "/twirp",
    interceptors: [{
      interceptUnary(next: NextUnaryFn, method: MethodInfo, input: object, options: RpcOptions): UnaryCall {
        const call = next(method, input, options)
        call.status.catch((e: RpcError) => {
          if (e.code == "unauthenticated") {
            return revalidate(getSession.key)
          }
        })
        return call
      }
    }]
  });

  return {
    helloWorld: new HelloWorldClient(transport),
    auth: new AuthClient(transport),
    page: new PageClient(transport),
    user: new UserClient(transport),
    admin: new AdminClient(transport)
  };
}

type ClientContextType = {
  helloWorld: HelloWorldClient
  auth: AuthClient
  page: PageClient
  user: UserClient
  admin: AdminClient
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
