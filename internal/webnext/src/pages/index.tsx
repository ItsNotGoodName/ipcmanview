import { ParentProps } from "solid-js";

export function Layout(props: ParentProps) {
  return (
    <>
      <div>IPCManView</div>
      <>{props.children}</>
    </>
  )
}
