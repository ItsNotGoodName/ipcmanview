import { A } from "@solidjs/router";
import { LayoutNormal } from "~/ui/Layout";
import { linkVariants } from "~/ui/Link";

export function AdminHome() {
  return (
    <LayoutNormal class="max-w-4xl">
      <A class={linkVariants()} href="./users">Users</A>
      <A class={linkVariants()} href="./groups">Groups</A>
      <A class={linkVariants()} href="./devices">Devices</A>
    </LayoutNormal>
  )
}
