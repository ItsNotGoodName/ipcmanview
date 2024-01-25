import { ParentProps } from "solid-js";

export function LayoutNormal(props: ParentProps) {
  return (
    <div class="flex justify-center p-4">
      <div class="flex w-full max-w-4xl flex-col gap-2">
        {props.children}
      </div>
    </div>
  )
}
