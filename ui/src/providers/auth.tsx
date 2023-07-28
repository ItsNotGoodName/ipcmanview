import { makePersisted } from "@solid-primitives/storage";
import {
  batch,
  createContext,
  createEffect,
  createResource,
  JSX,
  Match,
  ParentComponent,
  Resource,
  Switch,
  useContext
} from "solid-js";
import { createStore } from "solid-js/store";
import { AuthService, LoginArgs, User, UserService } from "~/core/client.gen";

type AuthContextType = {
  user: Resource<User>
  fetch: typeof fetch
  login: (args: LoginArgs) => Promise<void>
  logout: () => void
};

const AuthContext = createContext<AuthContextType>();

type AuthContextProps = {
  loading: JSX.Element;
  login: JSX.Element;
};

export const AuthProvider: ParentComponent<AuthContextProps> = (props) => {
  // Persist JWT token to storage
  const [storage, setStorage] = makePersisted(createStore<{ token: string }>({ token: "" }), { name: "auth" })

  // Update cookie when JWT token changes
  createEffect(() => {
    document.cookie =
      "auth_token=" +
      storage.token +
      `;Path=/file/` +
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
      batch(() => {
        setUser();
        setStorage({ token: "" });
      });
    }

    return res;
  })

  const authService = new AuthService(import.meta.env.VITE_BACKEND_URL, authFetch);
  const userService = new UserService(import.meta.env.VITE_BACKEND_URL, authFetch);

  // Resources
  const [user, { mutate: setUser }] = createResource(() => storage.token != "", () => userService.me().then((res) => res.user));
  createEffect(() => {
    console.log(user.state)
  })

  const store: AuthContextType = {
    user: user,
    fetch: authFetch,
    login: async (args) => {
      const res = await authService.login(args);
      batch(() => {
        setUser(res.user);
        setStorage({ token: res.token });
      });
    },
    logout: () => {
      batch(() => {
        setUser();
        setStorage({ token: "" });
      });
    }
  };

  return (
    <AuthContext.Provider value={store}>
      <Switch fallback={props.login}>
        <Match when={user.loading}>{props.loading}</Match>
        <Match when={storage.token != "" && user.state == "ready"}>{props.children}</Match>
      </Switch>
    </AuthContext.Provider>)
};

export function useAuth(): AuthContextType {
  const result = useContext(AuthContext);
  if (!result) throw new Error("useAuth must be used within a AuthProvider");
  return result;
}

