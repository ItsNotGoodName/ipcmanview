import { A, createAsync, useSearchParams } from "@solidjs/router";
import Humanize from "humanize-plus"
import { RiEditorAttachment2 } from "solid-icons/ri";
import { For, Show } from "solid-js";
import { Crud } from "~/components/Crud";
import { Shared } from "~/components/Shared";
import { createToggleSortField, formatDate, parseDate, parseOrder } from "~/lib/utils";
import { LayoutNormal } from "~/ui/Layout";
import { PaginationContent, PaginationEllipsis, PaginationItem, PaginationItems, PaginationLink, PaginationNext, PaginationPrevious, PaginationRoot } from "~/ui/Pagination";
import { TableBody, TableCaption, TableCell, TableHead, TableHeader, TableRoot, TableRow } from "~/ui/Table";
import { TooltipContent, TooltipRoot, TooltipTrigger } from "~/ui/Tooltip";
import { getEmailsPage } from "./Emails.data";
import { linkVariants } from "~/ui/Link";

export function Emails() {
  const [searchParams, setSearchParams] = useSearchParams()
  const data = createAsync(() => getEmailsPage({
    page: {
      page: Number(searchParams.page) || 0,
      perPage: Number(searchParams.perPage) || 0
    },
    sort: {
      field: searchParams.sort || "",
      order: parseOrder(searchParams.order)
    },
  }))
  const toggleSort = createToggleSortField(() => data()?.sort)

  return (
    <LayoutNormal>
      <Shared.Title>Emails</Shared.Title>
      <Crud.PerPageSelect
        class="w-20"
        perPage={data()?.pageResult?.perPage}
        onChange={(perPage) => setSearchParams({ perPage })}
      />
      <PaginationRoot
        page={data()?.pageResult?.page}
        count={data()?.pageResult?.totalPages || 0}
        onPageChange={(page) => setSearchParams({ page })}
        itemComponent={props => (
          <PaginationItem page={props.page}>
            <PaginationLink isActive={props.page == data()?.pageResult?.page}>
              {props.page}
            </PaginationLink>
          </PaginationItem>
        )}
        ellipsisComponent={() => (
          <PaginationEllipsis />
        )}
      >
        <PaginationContent>
          <PaginationPrevious />
          <PaginationItems />
          <PaginationNext />
        </PaginationContent>
      </PaginationRoot>
      <TableRoot>
        <TableHeader>
          <TableRow>
            <TableHead>Device</TableHead>
            <TableHead>
              <Crud.SortButton
                name=""
                onClick={toggleSort}
                sort={data()?.sort}
              >
                Created At
              </Crud.SortButton>
            </TableHead>
            <TableHead>From</TableHead>
            <TableHead>Subject</TableHead>
            <TableHead></TableHead>
          </TableRow>
        </TableHeader>
        <TableBody>
          <For each={data()?.emails}>
            {v =>
              <TableRow>
                <TableCell class="w-0 text-nowrap">
                  <A href={`/devices/${v.deviceId}`} class={linkVariants()}>
                    {v.deviceName}
                  </A>
                </TableCell>
                <TableCell class="w-0 text-nowrap">
                  {formatDate(parseDate(v.createdAtTime))}
                </TableCell>
                <TableCell class="w-0 text-nowrap">
                  {v.from}
                </TableCell>
                <TableCell>
                  {v.subject}
                </TableCell>
                <TableCell class="w-0 text-nowrap">
                  <Show when={v.attachmentCount > 0}>
                    <A href={`/emails/${v.id}?tab=attachments`}>
                      <TooltipRoot>
                        <TooltipTrigger class="flex h-full items-center">
                          <RiEditorAttachment2 />
                        </TooltipTrigger>
                        <TooltipContent>
                          {v.attachmentCount} {Humanize.pluralize(v.attachmentCount, "attachment")}
                        </TooltipContent>
                      </TooltipRoot>
                    </A>
                  </Show>
                </TableCell>
              </TableRow>
            }
          </For>
        </TableBody>
        <TableCaption>
          <Crud.Metadata pageResult={data()?.pageResult} />
        </TableCaption>
      </TableRoot>
    </LayoutNormal>
  )
}
