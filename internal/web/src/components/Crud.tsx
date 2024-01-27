import { RiArrowsArrowDownSLine } from "solid-icons/ri"
import { ParentProps } from "solid-js"

import { Order, PagePaginationResult, Sort } from "~/twirp/rpc"
import { cn } from "~/lib/utils"
import { SelectContent, SelectItem, SelectListbox, SelectRoot, SelectTrigger, SelectValue } from "~/ui/Select"
import { TableCell } from "~/ui/Table"

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

function LastTableCell(props: ParentProps) {
  return (
    <TableCell class="py-0">
      <div class="flex justify-end gap-2">
        {props.children}
      </div>
    </TableCell>
  )
}

export const Crud = {
  SortButton,
  Metadata,
  PerPageSelect,
  LastTableCell,
}
