import { Accessor, ParentProps, onCleanup } from 'solid-js'
import { createReconnectingWS, createWSState } from '@solid-primitives/websocket';
import {
  createContext,
  useContext
} from "solid-js";
import { useBus } from './bus';
import { relativeWSURL } from '~/lib/utils';
import { DahuaEvent, WSData, WSEvent } from '~/lib/models.gen';
import { revalidate } from '@solidjs/router';
import { getSession } from './session';

export enum WSState {
  Connecting,
  Connected,
  Disconnecting,
  Disconnected,
}

type WSContextType = {
  state: Accessor<WSState>
};

const WSContext = createContext<WSContextType>();

export function WSProvider(props: ParentProps) {
  const bus = useBus()

  const ws = createReconnectingWS(relativeWSURL("/v1/ws"));

  const onOpen = () => revalidate(getSession.key)

  const onMessage = (msg: MessageEvent<string>) => {
    const event = new WSData(msg.data)
    switch (event.type) {
      case "event":
        bus.event.emit(new WSEvent(event.data))
        break
      case "dahua-event":
        bus.dahuaEvent.emit(new DahuaEvent(event.data))
        break
    }
  }

  const onError = (event: Event) => {
    console.log(event)
  }

  ws.addEventListener("open", onOpen)
  ws.addEventListener("message", onMessage)
  ws.addEventListener("error", onError)
  onCleanup(() => {
    ws.removeEventListener("open", onOpen)
    ws.removeEventListener("message", onMessage)
    ws.removeEventListener("error", onError)
  })

  const store: WSContextType = {
    state: createWSState(ws),
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
