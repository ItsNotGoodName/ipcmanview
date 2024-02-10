import { RiArrowsArrowDownSLine, RiArrowsArrowLeftSLine, RiArrowsArrowRightSLine, RiSystemMore2Line } from "solid-icons/ri"
import { ParentProps } from "solid-js"

import { Order, PagePaginationResult, Sort } from "~/twirp/rpc"
import { cn } from "~/lib/utils"
import { SelectContent, SelectItem, SelectListbox, SelectRoot, SelectTrigger, SelectValue } from "~/ui/Select"
import { TableCell, TableHead } from "~/ui/Table"
import { Button } from "~/ui/Button"
import { DropdownMenuTrigger } from "~/ui/DropdownMenu"

function SortButton(props: ParentProps<{ onClick: (name: string) => void, name?: string, sort?: Sort }>) {
  const name = () => props.name ?? ""
  return (
    <button
      onClick={[props.onClick, name()]}
      class={cn("text-nowrap flex items-center whitespace-nowrap", name() == props.sort?.field && 'text-blue-500')}
    >
      {props.children}
      <RiArrowsArrowDownSLine data-selected={props.sort?.field == name() && props.sort.order == Order.ASC} class="h-5 w-5 transition-all data-[selected=true]:rotate-180" />
    </button>
  )
}

function PageMetadata(props: { pageResult?: PagePaginationResult }) {
  return (
    <div class="flex justify-between">
      <div>
        Seen {props.pageResult?.seenItems.toString() || 0} of {props.pageResult?.totalItems.toString() || 0}
      </div>
      <div>
        Page {props.pageResult?.page || 0} of {props.pageResult?.totalPages || 0}
      </div>
    </div>
  )
}

function PerPageSelect(props: { class?: string, perPage?: number, onChange: (value: number) => void }) {
  return (
    <SelectRoot
      class={props.class}
      value={props.perPage}
      onChange={(value) => {
        if (value == null)
          return
        props.onChange(value)
      }}
      options={[10, 25, 50, 100]}
      itemComponent={props => (
        <SelectItem item={props.item}>
          {props.item.rawValue}
        </SelectItem>
      )}
    >
      <SelectTrigger aria-label="Per page">
        <SelectValue<number>>
          {state => state.selectedOption()}
        </SelectValue>
      </SelectTrigger>
      <SelectContent>
        <SelectListbox />
      </SelectContent>
    </SelectRoot>
  )
}

function PageButtons(props: { class?: string, previousPageDisabled?: boolean, previousPage?: () => void, nextPageDisabled?: boolean, nextPage?: () => void }) {
  return (
    <div class={cn("flex gap-1", props.class)}>
      <Button
        aria-label="Go to previous page"
        title="Previous"
        size="icon"
        variant="outline"
        disabled={props.previousPageDisabled}
        onClick={props.previousPage}
      >
        <RiArrowsArrowLeftSLine class="h-5 w-5" />
      </Button>
      <Button
        aria-label="Go to next page"
        title="Next"
        size="icon"
        variant="outline"
        disabled={props.nextPageDisabled}
        onClick={props.nextPage}
      >
        <RiArrowsArrowRightSLine class="h-5 w-5" />
      </Button>
    </div>
  )
}

function LastTableCell(props: ParentProps) {
  return (
    <TableCell class="py-0">
      <div class="flex justify-end gap-2">
        {props.children}
      </div>
    </TableCell>
  )
}

function LastTableHead(props: ParentProps) {
  return (
    <TableHead>
      <div class="flex items-center justify-end">
        {props.children}
      </div>
    </TableHead>
  )
}

function MoreDropdownMenuTrigger() {
  return (
    <DropdownMenuTrigger class="hover:bg-accent hover:text-accent-foreground rounded p-1" title="Actions">
      <RiSystemMore2Line class="h-5 w-5" />
    </DropdownMenuTrigger>
  )
}

export const Crud = {
  SortButton,
  PageMetadata,
  PerPageSelect,
  PageButtons,
  LastTableCell,
  LastTableHead,
  MoreDropdownMenuTrigger
}
