import { makePersisted } from "@solid-primitives/storage";
import {
  batch,
  Component,
  createContext,
  createEffect,
  createSignal,
  JSX,
  Show,
  useContext
} from "solid-js";
import { createStore } from "solid-js/store";
import { AuthService, LoginArgs, User } from "~/core/client.gen";
import { UserProvider } from "./user";
import { BACKEND_URL } from "~/env";

type AuthContextType = {
  fetch: typeof fetch
  register: AuthService["register"]
  login: (args: LoginArgs) => Promise<void>
  logout: () => void
};

const AuthContext = createContext<AuthContextType>();

type AuthContextProps = {
  authenticated: JSX.Element;
  anonymous: JSX.Element;
};

export const AuthProvider: Component<AuthContextProps> = (props) => {
  // Persist JWT token to storage
  const [storage, setStorage] = makePersisted(createStore<{ token: string }>({ token: "" }), { name: "auth" })

  // Update cookie when JWT token changes, cookie is used to fetch protected HTTP resources (e.g. images)
  createEffect(() => {
    document.cookie =
      "auth_token=" +
      storage.token +
      `;Path=/file/` + // TODO: set correct path for images
      import.meta.env.VITE_COOKIE_ATTRIBUTES;
  })

  // Custom fetcher with JWT token
  const authFetch = (input: RequestInfo | URL, init?: RequestInit) => fetch(input, {
    ...init, headers: {
      ...init?.headers,
      "Authorization": `Bearer ${storage.token}`
    }
  }).then((res) => {
    if (res.status == 401 && storage.token != "") {
      console.log("No longer authenticated.");
      setStorage({ token: "" });
    }

    return res;
  })

  const authService = new AuthService(BACKEND_URL, authFetch);

  const [initialUser, setInitialUser] = createSignal<User>({
    id: 0,
    email: "",
    username: "",
    created_at: ""
  })

  const store: AuthContextType = {
    fetch: authFetch,
    register: authService.register,
    login: async (args) => {
      const res = await authService.login(args);
      batch(() => {
        setStorage({ token: res.token });
        setInitialUser(res.user);
      })
    },
    logout: () => {
      setStorage({ token: "" });
    }
  };

  return (
    <AuthContext.Provider value={store}>
      <Show when={storage.token != ""} fallback={props.anonymous}>
        <UserProvider initialUser={initialUser()}>
          {props.authenticated}
        </UserProvider>
      </Show>
    </AuthContext.Provider>)
};


export function useAuthStore(): AuthContextType {
  const result = useContext(AuthContext);
  if (!result) throw new Error("useAuthStore must be used within a AuthProvider");
  return result;
}

