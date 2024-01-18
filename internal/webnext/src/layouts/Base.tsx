import { RouteSectionProps } from "@solidjs/router";

export function BaseLayout(props: RouteSectionProps) {
  return (
    <>
      {props.children}
    </>
  )
}
