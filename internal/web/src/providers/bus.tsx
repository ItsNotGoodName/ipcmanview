import { EventBus, EventHub, createEventBus, createEventHub } from '@solid-primitives/event-bus';
import {
  createContext,
  ParentComponent,
  useContext
} from "solid-js";

export type BusForEvent = {
  action: string,
  data: unknown
}

type BusContextType = EventHub<{
  event: EventBus<BusForEvent>;
}>

const BusContext = createContext<BusContextType>();

type BusContextProps = {};

export const BusProvider: ParentComponent<BusContextProps> = (props) => {
  const store = createEventHub({
    event: createEventBus<BusForEvent>()
  })

  return (
    <BusContext.Provider value={store}>
      {props.children}
    </BusContext.Provider>)
};

export function useBus(): BusContextType {
  const result = useContext(BusContext);
  if (!result) throw new Error("useBus must be used within a BusProvider");
  return result;
}
