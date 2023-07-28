import {
  createContext,
  createResource,
  InitializedResourceReturn,
  ParentComponent,
  Show,
  useContext
} from "solid-js";
import { User, UserService } from "~/core/client.gen";
import { useAuthStore } from "./auth";
import { BACKEND_URL } from "~/env";

type UserContextType = {
  userResource: InitializedResourceReturn<User>
};

const UserContext = createContext<UserContextType>();

type UserContextProps = {
  initialUser: User
};

export const UserProvider: ParentComponent<UserContextProps> = (props) => {
  const authStore = useAuthStore()

  const userService = new UserService(BACKEND_URL, authStore.fetch);

  const userResource = createResource(() => userService.me().then((res) => res.user), {
    initialValue: props.initialUser,
  });
  const [user] = userResource

  const store: UserContextType = {
    userResource: userResource
  };

  return (
    <UserContext.Provider value={store}>
      <Show when={user.state == "refreshing" || user.state == "ready"} fallback={<div>{user.error}</div>}>
        {props.children}
      </Show>
    </UserContext.Provider>)
};


export function useUserStore(): UserContextType {
  const result = useContext(UserContext);
  if (!result) throw new Error("useUserStore must be used within a userProvider");
  return result;
}

