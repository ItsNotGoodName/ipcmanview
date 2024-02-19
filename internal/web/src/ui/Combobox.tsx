import { Combobox } from "@kobalte/core";
import { RiSystemAddLine, RiSystemCheckLine, RiSystemSearchLine } from "solid-icons/ri";
import { For, Show, createSignal } from "solid-js";
import { buttonVariants } from "./Button";
import { cn } from "~/lib/utils";
import { Seperator } from "./Seperator";
import { cva } from "class-variance-authority";

const ALL_OPTIONS = ["Apple", "Banana", "Blueberry", "Grapes", "Pineapple"];

const tagVariants = cva("focus:ring-ring bg-secondary text-secondary-foreground hover:bg-secondary/80 inline-flex items-center rounded-sm border border-transparent px-1 py-0.5 text-xs font-normal transition-colors focus:outline-none focus:ring-2 focus:ring-offset-2")

const itemVariants = cva("ui-highlighted:bg-accent ui-highlighted::text-accent-foreground hover:bg-accent hover:text-accent-foreground relative flex w-full cursor-default select-none items-center justify-center gap-2 rounded-sm px-2 py-1.5 text-sm outline-none transition-colors disabled:pointer-events-none disabled:opacity-50", {
  variants: {
    variant: {
      button: "w-full",
      item: "group",
    },
  },
})

export function MultipleSelectionExample() {
  const [values, setValues] = createSignal(["Blueberry", "Grapes"]);
  return (
    <Combobox.Root<string>
      multiple
      options={ALL_OPTIONS}
      value={values()}
      onChange={setValues}
      placeholder="Search"
      itemComponent={props => (
        <Combobox.Item item={props.item} class={itemVariants({ variant: "item" })}>
          <div class="border-primary group-data-[selected]:bg-primary group-data-[selected]:text-primary-foreground size-4 shrink-0 rounded-sm border shadow">
            <Combobox.ItemIndicator class="flex items-center justify-center text-current">
              <RiSystemCheckLine class="size-4" />
            </Combobox.ItemIndicator>
          </div>
          <Combobox.ItemLabel class="flex-1">{props.item.rawValue}</Combobox.ItemLabel>
        </Combobox.Item>
      )}
      allowsEmptyCollection
    >
      <Combobox.Control<string> aria-label="Fruits">
        {state => (
          <Combobox.Trigger class={cn(buttonVariants({ variant: "outline", size: "sm" }), "flex items-center gap-2")}>
            <Combobox.Icon>
              <RiSystemAddLine class="size-4" />
            </Combobox.Icon>
            Fruits
            <Show when={state.selectedOptions().length > 0}>
              <Seperator orientation="vertical" class="h-4" />
              <div class={cn(tagVariants(), "lg:hidden")}>
                {state.selectedOptions().length}
              </div>
              <div class="hidden space-x-1 lg:flex">
                <Show when={state.selectedOptions().length < 3}
                  fallback={<span class={tagVariants()}>{state.selectedOptions().length} selected</span>}
                >
                  <For each={state.selectedOptions()}>
                    {option => <span class={tagVariants()}>{option}</span>}
                  </For>
                </Show>
              </div>
            </Show>
          </Combobox.Trigger>
        )}
      </Combobox.Control>
      <Combobox.Control<string> aria-label="Fruits">
        {state => (
          <Combobox.Portal>
            <Combobox.Content class="bg-popover text-popover-foreground ui-expanded:animate-in ui-not-expanded:animate-out ui-not-expanded:fade-out-0 ui-expanded:fade-in-0 ui-not-expanded:zoom-out-95 ui-expanded:zoom-in-95 z-50 w-[200px] max-w-[var(--kb-popper-content-available-width)] origin-[var(--kb-combobox-content-transform-origin)] rounded-md border shadow-md outline-none">
              <div class="flex items-center gap-2 border-b px-3">
                <RiSystemSearchLine class="size-4 shrink-0 opacity-50" />
                <Combobox.Input class="placeholder:text-muted-foreground flex h-10 w-full rounded-md bg-transparent py-3 text-sm outline-none disabled:cursor-not-allowed disabled:opacity-50" />
              </div>
              <Combobox.Listbox class="bg-background p-1" />
              <Show when={state.selectedOptions().length > 0}>
                <Seperator />
                <div class="p-1">
                  <button onPointerDown={e => e.stopPropagation()} onClick={state.clear} class={itemVariants({ variant: "button" })} >
                    Clear filters
                  </button>
                </div>
              </Show>
            </Combobox.Content>
          </Combobox.Portal>
        )}
      </Combobox.Control>
    </Combobox.Root >
  );
}

