import { As, Popover } from "@kobalte/core";
import { createSignal, splitProps } from "solid-js";

import { PopoverArrow, PopoverCloseButton, PopoverContent, PopoverPortal, PopoverRoot, PopoverTrigger } from "./Popover";
import { Button } from "./Button";

export type ConfirmButtonProps = {
  message?: string,
  disabled?: boolean,
  onYes?: () => Promise<unknown>
} & Popover.PopoverTriggerProps

export function ConfirmButton(props: ConfirmButtonProps) {
  const [_, rest] = splitProps(props, ["message", "onYes"])
  const [open, setOpen] = createSignal(false);
  const onYes = () => {
    if (props.onYes) {
      props.onYes().then(() => setOpen(false))
    } else {
      setOpen(false)
    }
  }

  return (
    <PopoverRoot open={open()} onOpenChange={setOpen}>
      <PopoverTrigger {...rest} />
      <PopoverPortal>
        <PopoverContent class="flex flex-col gap-2">
          <PopoverArrow />
          <div>{props.message}</div>
          <div class="flex justify-end gap-2">
            <PopoverCloseButton asChild>
              <As component={Button} size="sm" disabled={props.disabled}>No</As>
            </PopoverCloseButton>
            <Button
              onClick={onYes}
              disabled={props.disabled}
              size="sm"
              variant="destructive"
            >
              Yes
            </Button>
          </div>
        </PopoverContent>
      </PopoverPortal>
    </PopoverRoot>
  )
}
