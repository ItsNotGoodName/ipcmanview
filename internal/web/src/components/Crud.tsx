import { RiArrowsArrowDownSLine, RiArrowsArrowLeftSLine, RiArrowsArrowRightSLine, RiSystemMore2Line } from "solid-icons/ri"
import { ParentProps } from "solid-js"

import { Order, PagePaginationResult, Sort } from "~/twirp/rpc"
import { cn } from "~/lib/utils"
import { SelectContent, SelectItem, SelectListbox, SelectRoot, SelectTrigger, SelectValue } from "~/ui/Select"
import { TableCell, TableHead } from "~/ui/Table"
import { Button } from "~/ui/Button"
import { DropdownMenuTrigger } from "~/ui/DropdownMenu"

function SortButton(props: ParentProps<{ onClick: (name: string) => void, name: string, sort?: Sort }>) {
  return (
    <button
      onClick={[props.onClick, props.name]}
      class={cn("text-nowrap flex items-center whitespace-nowrap text-lg", props.name == props.sort?.field && 'text-blue-500')}
    >
      {props.children}
      <RiArrowsArrowDownSLine data-selected={props.sort?.field == props.name && props.sort.order == Order.ASC} class="h-5 w-5 transition-all data-[selected=true]:rotate-180" />
    </button>
  )
}

function Metadata(props: { pageResult?: PagePaginationResult }) {
  return (
    <div class="flex justify-between">
      <div>
        {props.pageResult?.seenItems.toString() || 0} / {props.pageResult?.totalItems.toString() || 0}
      </div>
      <div>
        Page {props.pageResult?.page || 0} / {props.pageResult?.totalPages || 0}
      </div>
    </div>
  )
}

function PerPageSelect(props: { class?: string, perPage?: number, onChange: (value: number) => void }) {
  return (
    <SelectRoot
      class={props.class}
      value={props.perPage}
      onChange={props.onChange}
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

function PageButtons(props: { previousPageDisabled: boolean, previousPage: () => void, nextPageDisabled: boolean, nextPage: () => void }) {
  return (
    <div class="flex gap-2">
      <Button
        title="Previous"
        size="icon"
        disabled={props.previousPageDisabled}
        onClick={props.previousPage}
      >
        <RiArrowsArrowLeftSLine class="h-6 w-6" />
      </Button>
      <Button
        title="Next"
        size="icon"
        disabled={props.nextPageDisabled}
        onClick={props.nextPage}
      >
        <RiArrowsArrowRightSLine class="h-6 w-6" />
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
  Metadata,
  PerPageSelect,
  PageButtons,
  LastTableCell,
  LastTableHead,
  MoreDropdownMenuTrigger
}
