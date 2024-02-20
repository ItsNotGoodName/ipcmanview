import { Accessor, ParentProps, onCleanup } from 'solid-js'
import { createReconnectingWS, createWSState } from '@solid-primitives/websocket';
import {
  createContext,
  useContext
} from "solid-js";
import { useBus } from './bus';
import { relativeWSURL } from '~/lib/utils';

export enum WSState {
  Connecting,
  Connected,
  Disconnecting,
  Disconnected,
}

type WSContextType = {
  state: Accessor<WSState>
};

type BaseEvent<T extends string, D> = {
  type: T
  data: D
}

type WSEvent = BaseEvent<"event", { action: string, data: any }>

const WSContext = createContext<WSContextType>();

export function WSProvider(props: ParentProps) {
  const ws = createReconnectingWS(relativeWSURL("/v1/ws"));

  const webSocketState = createWSState(ws);

  const bus = useBus()

  const onMessage = (msg: MessageEvent<string>) => {
    const event = JSON.parse(msg.data) as WSEvent

    switch (event.type) {
      case "event":
        bus.event.emit(event.data)
        break
    }
  }

  const onError = (event: Event) => {
    console.log(event)
  }

  ws.addEventListener("message", onMessage)
  ws.addEventListener("error", onError)
  onCleanup(() => {
    ws.removeEventListener("message", onMessage)
    ws.removeEventListener("error", onError)
  })

  const store: WSContextType = {
    state: webSocketState,
  };

  return (
    <WSContext.Provider value={store}>
      {props.children}
    </WSContext.Provider>
  )
};

export function useWS(): WSContextType {
  const result = useContext(WSContext);
  if (!result) throw new Error("useWS must be used within a WSProvider");
  return result;
}
