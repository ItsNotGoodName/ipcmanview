import { Separator } from "@kobalte/core";
import { JSX, splitProps } from "solid-js";

export function Seperator(props: JSX.HTMLAttributes<HTMLDivElement> & { orientation?: "horizontal" | "vertical" }) {
  const [_, rest] = splitProps(props, ["orientation"])
  return <div {...rest}>
    <Separator.Root
      class={"bg-border ui-horizontal:h-[1px] ui-horizontal:w-full ui-vertical:w-[1px] ui-vertical:h-full shrink-0"}
      orientation={props.orientation}
    />
  </div>
}
