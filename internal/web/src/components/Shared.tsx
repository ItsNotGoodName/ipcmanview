import { ParentProps } from "solid-js";
import { Seperator } from "~/ui/Seperator";

function Title(props: ParentProps) {
  return (
    <div class="flex flex-col gap-2">
      <div class="text-xl">{props.children}</div>
      <Seperator />
    </div>
  )
}

export const Shared = {
  Title
}
