import { As } from "@kobalte/core";
import { Button } from "./Button";
import { PopoverArrow, PopoverCloseButton, PopoverContent, PopoverPortal, PopoverRoot, PopoverTrigger } from "./Popover";
import { ComponentProps, createSignal, splitProps } from "solid-js";

export type ConfirmButtonProps = {
  message?: string,
  pending?: boolean,
  onYes?: () => Promise<unknown>
} & Omit<ComponentProps<typeof Button>, "disabled">

export function ConfirmButton(props: ConfirmButtonProps) {
  const [_, rest] = splitProps(props, ["message", "pending", "onYes"])
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
      <PopoverTrigger asChild>
        <As component={Button} disabled={props.pending} {...rest}>
          {props.children}
        </As>
      </PopoverTrigger>
      <PopoverPortal>
        <PopoverContent class="flex flex-col gap-2">
          <PopoverArrow />
          <div>{props.message}</div>
          <div class="flex gap-4">
            <PopoverCloseButton asChild>
              <As component={Button} size="sm" disabled={props.pending}>No</As>
            </PopoverCloseButton>
            <Button
              onClick={onYes}
              disabled={props.pending}
              size="sm"
              variant={props.variant}
            >
              Yes
            </Button>
          </div>
        </PopoverContent>
      </PopoverPortal>
    </PopoverRoot>
  )
}
