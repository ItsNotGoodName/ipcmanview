import { Combobox } from "@kobalte/core";
import { RiSystemCheckLine, RiSystemCloseLine, RiSystemSearchLine } from "solid-icons/ri";
import { Accessor, For, Show, mergeProps, splitProps } from "solid-js";
import { cva } from "class-variance-authority";

import { buttonVariants } from "./Button";
import { cn } from "~/lib/utils";
import { Seperator } from "./Seperator";

export interface ComboboxControlState<T> {
  /** The selected options. */
  selectedOptions: Accessor<T[]>;
  /** A function to remove an option from the selection. */
  remove: (option: T) => void;
  /** A function to clear the selection. */
  clear: () => void;
}

const tagVariants = cva("focus:ring-ring bg-secondary text-secondary-foreground hover:bg-secondary/80 inline-flex items-center rounded-sm border border-transparent px-1 py-0.5 text-xs font-normal transition-colors focus:outline-none focus:ring-2 focus:ring-offset-2")

export function ComboboxRoot<Option, OptGroup = never>(props: Combobox.ComboboxRootProps<Option, OptGroup>) {
  return <Combobox.Root allowsEmptyCollection {...props} />
}

export function ComboboxItem(props: Combobox.ComboboxItemProps) {
  const [_, rest] = splitProps(props, ["class"])
  return <Combobox.Item class={cn("ui-highlighted:bg-accent ui-highlighted:text-accent-foreground hover:bg-accent hover:text-accent-foreground group relative flex w-full cursor-default select-none items-center justify-start gap-2 rounded-sm px-2 py-1.5 text-sm outline-none transition-colors disabled:pointer-events-none disabled:opacity-50", props.class)} {...rest}>
    <div class="size-4 border-primary group-data-[selected]:bg-primary group-data-[selected]:text-primary-foreground flex shrink-0 items-center justify-center rounded-sm border">
      <Combobox.ItemIndicator class="flex items-center justify-center text-current">
        <RiSystemCheckLine class="size-4" />
      </Combobox.ItemIndicator>
    </div>
    {props.children}
  </Combobox.Item>
}

export const ComboboxControl = Combobox.Control

export function ComboboxTrigger(props: Combobox.ComboboxTriggerProps) {
  const [_, rest] = splitProps(props, ["class", "children"])
  return <Combobox.Trigger class={cn(buttonVariants({ variant: "outline", size: "sm" }), "flex items-center gap-2", props.class)} {...rest}>
    {props.children}
  </Combobox.Trigger>
}

export const ComboboxIcon = Combobox.Icon

export function ComboboxState<Option>(props: { state: ComboboxControlState<Option>, optionToString?: (option: Option) => string }) {
  const mergedProps = mergeProps(props, { optionToString: (option: any) => (option as string) })
  return <Show when={mergedProps.state.selectedOptions().length > 0}>
    <Seperator orientation="vertical" class="h-4" />
    <div class={cn(tagVariants(), "lg:hidden")}>
      {mergedProps.state.selectedOptions().length}
    </div>
    <div class="hidden space-x-1 lg:flex">
      <Show when={mergedProps.state.selectedOptions().length < 3} fallback={
        <span class={tagVariants()}>{mergedProps.state.selectedOptions().length} selected</span>
      }>
        <For each={mergedProps.state.selectedOptions()}>
          {option => <span class={tagVariants()}>{mergedProps.optionToString(option)}</span>}
        </For>
      </Show>
    </div>
  </Show>
}

export function ComboboxReset<Option>(props: { class?: string, state: ComboboxControlState<Option> }) {
  return <Show when={props.state.selectedOptions().length > 0}>
    <button class="h-full" onPointerDown={e => e.stopPropagation()} onClick={props.state.clear}>
      <RiSystemCloseLine class={props.class} />
      <span class="sr-only">Reset</span>
    </button>
  </Show>
}

export const ComboboxItemLabel = Combobox.ItemLabel

export function ComboboxContent(props: Combobox.ComboboxContentProps) {
  const [_, rest] = splitProps(props, ["class"])
  return <Combobox.Portal>
    <Combobox.Content class={cn("bg-popover text-popover-foreground ui-expanded:animate-in ui-not-expanded:animate-out ui-not-expanded:fade-out-0 ui-expanded:fade-in-0 ui-not-expanded:zoom-out-95 ui-expanded:zoom-in-95 z-50 w-[200px] max-w-[var(--kb-popper-content-available-width)] origin-[var(--kb-combobox-content-transform-origin)] rounded-md border shadow-md outline-none", props.class)} {...rest} />
  </Combobox.Portal>
}

export function ComboboxInput(props: Omit<Combobox.ComboboxInputProps, "class">) {
  return <div class="flex items-center gap-2 border-b px-3">
    <RiSystemSearchLine class="size-4 shrink-0 opacity-50" />
    <Combobox.Input class="placeholder:text-muted-foreground flex h-10 w-full rounded-md bg-transparent py-3 text-sm outline-none disabled:cursor-not-allowed disabled:opacity-50" {...props} />
  </div>
}

export function ComboboxListbox<Option, OptGroup>(props: Omit<Combobox.ComboboxListboxProps<Option, OptGroup>, "class">) {
  return <Combobox.Listbox class="max-h-48 overflow-y-auto p-1" {...props} />
}