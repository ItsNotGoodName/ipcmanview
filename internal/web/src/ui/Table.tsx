import { RiArrowsArrowDownSLine } from "solid-icons/ri"
import { JSX, ParentProps, splitProps } from "solid-js"

import { cn } from "~/lib/utils"
import { Order, PagePaginationResult, Sort } from "~/twirp/rpc"

export function TableRoot(props: JSX.HTMLAttributes<HTMLTableElement>) {
  const [_, rest] = splitProps(props, ["class"])
  return <div class="relative w-full overflow-auto">
    <table
      class={cn("w-full caption-bottom text-sm", props.class)}
      {...rest}
    />
  </div>
}

export function TableHeader(props: JSX.HTMLAttributes<HTMLTableSectionElement>) {
  const [_, rest] = splitProps(props, ["class"])
  return <thead class={cn("[&_tr]:border-b", props.class)} {...rest} />
}

export function TableBody(props: JSX.HTMLAttributes<HTMLTableSectionElement>) {
  const [_, rest] = splitProps(props, ["class"])
  return <tbody
    class={cn("[&_tr:last-child]:border-0", props.class)}
    {...rest}
  />
}

export function TableFooter(props: JSX.HTMLAttributes<HTMLTableSectionElement>) {
  const [_, rest] = splitProps(props, ["class"])
  return <tfoot
    class={cn(
      "bg-muted/50 border-t font-medium [&>tr]:last:border-b-0",
      props.class
    )}
    {...rest}
  />
}

export function TableRow(props: JSX.HTMLAttributes<HTMLTableRowElement>) {
  const [_, rest] = splitProps(props, ["class"])
  return <tr
    class={cn(
      "hover:bg-muted/50 data-[state=selected]:bg-muted border-b transition-colors",
      props.class
    )}
    {...rest}
  />
}

export function TableHead(props: JSX.ThHTMLAttributes<HTMLTableCellElement>) {
  const [_, rest] = splitProps(props, ["class"])
  return <th
    class={cn(
      "text-muted-foreground h-12 px-4 text-left align-middle font-medium [&:has([role=checkbox])]:pr-0",
      props.class
    )}
    {...rest}
  />
}

export function TableCell(props: JSX.TdHTMLAttributes<HTMLTableCellElement>) {
  const [_, rest] = splitProps(props, ["class"])
  return <td
    class={cn("p-4 align-middle [&:has([role=checkbox])]:pr-0", props.class)}
    {...rest}
  />
}

export function TableCaption(props: JSX.HTMLAttributes<HTMLElement>) {
  const [_, rest] = splitProps(props, ["class"])
  return <caption
    class={cn("text-muted-foreground mt-4 text-sm", props.class)}
    {...rest}
  />
}

export function TableSortButton(props: ParentProps<{ onClick: (name: string) => void, name: string, sort?: Sort }>) {
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

export function TableMetadata(props: { pageResult?: PagePaginationResult }) {
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
