import { cva } from "class-variance-authority";
import { ParentProps } from "solid-js";
import { Seperator } from "~/ui/Seperator";

function Title(props: ParentProps) {
  return (
    <div class="flex flex-col gap-2">
      <div class="text-lg">{props.children}</div>
      <Seperator />
    </div>
  )
}

const connectionIndicatorVariants = cva("size-4 rounded-full shrink-0", {
  variants: {
    state: {
      connected: "bg-lime-500",
      connecting: "bg-orange-500",
      disconnected: "bg-red-500"
    },
  },
})

export const Shared = {
  Title,
  connectionIndicatorVariants
}
