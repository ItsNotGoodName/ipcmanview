import { ParentProps } from "solid-js";
import { cn } from "~/lib/utils";

export function LayoutNormal(props: ParentProps<{ class?: string }>) {
  return (
    <div class="flex justify-center p-4">
      <div class={cn("flex w-full flex-col gap-2", props.class)}>
        {props.children}
      </div>
    </div>
  )
}
