import {
  createContext,
  ParentComponent,
  useContext
} from "solid-js";
import { useAuthStore } from "./auth";
import { BACKEND_URL } from "~/env";
import { DahuaService } from "~/core/client.gen";

type ServiceContextType = {
  dahuaService: DahuaService
};

const ServiceContext = createContext<ServiceContextType>();

type ServiceContextProps = {};

export const ServiceProvider: ParentComponent<ServiceContextProps> = (props) => {
  const authStore = useAuthStore()

  const store: ServiceContextType = {
    dahuaService: new DahuaService(BACKEND_URL, authStore.fetch),
  };

  return (
    <ServiceContext.Provider value={store}>
      {props.children}
    </ServiceContext.Provider>)
};

export function useService(): ServiceContextType {
  const result = useContext(ServiceContext);
  if (!result) throw new Error("useService must be used within a ServiceProvider");
  return result;
}
