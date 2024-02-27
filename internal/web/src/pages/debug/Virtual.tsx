import { createAsync } from "@solidjs/router"
import { Show, } from "solid-js"
import { getListLocations } from "../admin/data"
import { SelectContent, SelectListBoxVirtual, SelectPortal, SelectRoot, SelectTrigger, SelectValue } from "~/ui/Select"
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
            <SelectContent>
              <SelectListBoxVirtual options={listLocations()!} />
            </SelectContent>
          </SelectPortal>
        </SelectRoot>
      </Show>
    </LayoutNormal>
  )
}
