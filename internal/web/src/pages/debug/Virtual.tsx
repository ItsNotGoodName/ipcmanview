import { createAsync } from "@solidjs/router"
import { createVirtualizer } from "@tanstack/solid-virtual"
import { For, Show, } from "solid-js"
import { getListLocations } from "../admin/data"
import { SelectContent, SelectItem, SelectListbox, SelectPortal, SelectRoot, SelectTrigger, SelectValue } from "~/ui/Select"
import { LayoutNormal } from "~/ui/Layout"

export function Virtual() {
  const listLocations = createAsync(() => getListLocations())

  return (
    <LayoutNormal>
      <Show when={listLocations()}>
        <SelectRoot
          virtualized={true}
          options={listLocations()!}
          placeholder="Select an item…"

        >
          <SelectTrigger aria-label="Food">
            <SelectValue<string>>
              {state => state.selectedOption()}
            </SelectValue>
          </SelectTrigger>
          <SelectPortal>
            <SelectContentVirtual options={listLocations()!} />
          </SelectPortal>
        </SelectRoot>
      </Show>
    </LayoutNormal>
  )
}

function SelectContentVirtual(props: { options: string[] }) {
  let listboxRef: HTMLUListElement | null;
  const virtualizer = createVirtualizer({
    count: props.options.length,
    getScrollElement: () => listboxRef,
    getItemKey: (index) => props.options[index],
    estimateSize: () => 32,
    overscan: 5,
  });

  return (
    <SelectContent>
      <SelectListbox<number>
        ref={listboxRef!}
        scrollToItem={(item) => virtualizer.scrollToIndex(props.options.indexOf(item))}
      >
        {items => (
          <div
            style={{
              height: `${virtualizer.getTotalSize()}px`,
              width: "100%",
              position: "relative",
            }}
          >
            <For each={virtualizer.getVirtualItems()}>
              {virtualRow => {
                const item = items().getItem(virtualRow.key as string);
                if (item) {
                  return (
                    <SelectItem
                      item={item}
                      style={{
                        position: "absolute",
                        top: 0,
                        left: 0,
                        width: "100%",
                        height: `${virtualRow.size}px`,
                        transform: `translateY(${virtualRow.start}px)`,
                      }}
                    >
                      {item.rawValue}
                    </SelectItem>
                  );
                }
              }}
            </For>
          </div>
        )}
      </SelectListbox>
    </SelectContent>
  );
}
