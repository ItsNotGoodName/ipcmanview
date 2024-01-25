import { A } from "@solidjs/router";
import { linkVariants } from "~/ui/Link";

export function AdminHome() {
  return (
    <div class="flex flex-col">
      <A class={linkVariants()} href="./users">Users</A>
      <A class={linkVariants()} href="./groups">Groups</A>
    </div>
  )
}
