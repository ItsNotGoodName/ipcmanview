import { createResizeObserver } from "@solid-primitives/resize-observer"
import { Index, JSX, Show, batch, createMemo, createSignal, onMount, splitProps } from "solid-js"
import { createStore } from "solid-js/store"
import { Component } from "solid-js/types/server/rendering.js"
import { Dynamic } from "solid-js/web"
import { cn } from "~/lib/utils"
import { Button } from "~/ui/Button"

type Rectangles = {
  x: number
  y: number
  w: number
  h: number
}

export function View() {
  const [rectangles, setRectangles] = createStore<Rectangles[]>([
    { x: 0, y: 0, w: 4096, h: 4096 },
    { x: 4096, y: 4096, w: 4096, h: 4096 },
    { x: 0, y: 4096, w: 4096, h: 4096 },
    { x: 4096, y: 0, w: 4096, h: 4096 },
    { x: 4096, y: 0, w: 4096, h: 4096 },
  ])

  const panels = [
    () => {
      console.log("render", 1)
      return <div class="h-full w-full bg-red-500">1</div>
    },
    () => {
      console.log("render", 2)
      return <div class="h-full w-full bg-green-500">2</div>
    },
    () => {
      console.log("render", 3)
      return <div class="h-full w-full bg-blue-500">3</div>
    },
    () => {
      console.log("render", 4)
      return <div class="h-full w-full bg-yellow-500">4</div>
    },
  ];

  return (
    <div class="flex h-screen flex-col">
      <Positioner class="flex-1 overflow-hidden" rectangles={rectangles} panels={panels} />
      <div class="flex justify-between gap-2 overflow-x-auto p-1">
        <div class="flex gap-1">
          <Button onClick={() => setRectangles((_, i) => i == 1, "y", y => y - 10)}>Up</Button>
          <Button onClick={() => setRectangles((_, i) => i == 1, "x", x => x - 10)}>Left</Button>
          <Button onClick={() => setRectangles((_, i) => i == 1, "x", x => x + 10)}>Right</Button>
          <Button onClick={() => setRectangles((_, i) => i == 1, "y", y => y + 10)}>Down</Button>
        </div>
        <div class="flex gap-1">
          <Button>Edit</Button>
        </div>
      </div>
    </div>
  )
}

type PositionProps = JSX.HTMLAttributes<HTMLDivElement> & { rectangles: Rectangles[], panels: Component[] }

function Positioner(props: PositionProps) {
  const [_, rest] = splitProps(props, ["class", "rectangles", "panels"])
  const { rectangles, resize } = useScaleRectangles(props.rectangles)

  let ref: HTMLDivElement;
  onMount(() => {
    createResizeObserver(ref, ({ width, height }, el) => {
      if (el === ref) resize(width, height);
    });
  });

  return (
    <div ref={ref!} class={cn("relative", props.class)} {...rest}>
      <Index each={rectangles()}>
        {(rectangle, i) =>
          <div class="absolute" style={{ "top": `${rectangle().y}px`, "left": `${rectangle().x}px`, "width": `${rectangle().w}px`, "height": `${rectangle().h}px` }}>
            <Show when={props.panels[i]}>
              <Dynamic component={props.panels[i]} />
            </Show>
          </div>
        }
      </Index>
    </div>
  )
}

const useScaleRectangles = (original: Rectangles[]) => {
  const scale = 8192
  const [width, setWidth] = createSignal(scale)
  const [height, setHeight] = createSignal(scale)

  const resize = (width: number, height: number) => {
    batch(() => {
      setWidth(width)
      setHeight(height)
    })
  }

  const rectangles = createMemo(() => {
    const next: Rectangles[] = []
    for (let i = 0; i < original.length; i++) {
      next.push({
        x: (original[i].x * width()) / scale,
        y: (original[i].y * height()) / scale,
        w: (original[i].w * width()) / scale,
        h: (original[i].h * height()) / scale,
      })
    }
    return next
  })

  return {
    rectangles,
    resize,
  }
}
