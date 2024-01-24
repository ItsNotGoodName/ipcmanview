import { A } from "@solidjs/router";
import { linkVariants } from "~/ui/Link";

export function AdminHome() {
  return (
    <>
      <A class={linkVariants()} href="./groups">Groups</A>
    </>
  )
}
