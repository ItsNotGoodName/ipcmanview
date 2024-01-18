import { makePersisted } from "@solid-primitives/storage";
import {
  Accessor,
  createContext,
  createEffect,
  ParentProps,
  useContext
} from "solid-js";
import { createStore } from "solid-js/store";

type AuthContextType = {
  session: Accessor<string>
  setSession: (token: string) => void
  valid: Accessor<boolean>
  clear: () => void
};

const AuthContext = createContext<AuthContextType>();

export function AuthProvider(props: ParentProps) {
  const [storage, setStorage] = makePersisted(createStore<{ session: string, valid: boolean }>({ session: "", valid: true }), { name: "auth" })

  createEffect(() => {
    document.cookie =
      "session=" +
      storage.session +
      ";SameSite=Strict"
  })

  const store: AuthContextType = {
    valid: () => storage.valid,
    session: () => storage.session,
    setSession: (session) => setStorage({ session, valid: true }),
    clear: () => setStorage({ session: "", valid: false })
  }

  return (
    <AuthContext.Provider value={store}>
      {props.children}
    </AuthContext.Provider>)
};

export function useAuth(): AuthContextType {
  const result = useContext(AuthContext);
  if (!result) throw new Error("useAuth must be used within a AuthProvider");
  return result;
}
