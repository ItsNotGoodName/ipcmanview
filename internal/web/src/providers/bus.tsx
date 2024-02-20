import { EventBus, createEventBus } from '@solid-primitives/event-bus';
import {
  createContext,
  ParentComponent,
  useContext
} from "solid-js";

export type EventBusType = {
  event: { action: string, data: unknown }
}

type BusContextType = {
  bus: EventBus<EventBusType>
};

const BusContext = createContext<BusContextType>();

type BusContextProps = {};

export const BusProvider: ParentComponent<BusContextProps> = (props) => {
  const bus = createEventBus<EventBusType>();

  const store: BusContextType = {
    bus
  };

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
