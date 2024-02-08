import { A, createAsync, useSearchParams } from "@solidjs/router";
import Humanize from "humanize-plus"
import { RiArrowsArrowRightLine, RiEditorAttachment2 } from "solid-icons/ri";
import { ErrorBoundary, For, Show, Suspense } from "solid-js";
import { Crud } from "~/components/Crud";
import { Shared } from "~/components/Shared";
import { createPagePagination, createToggleSortField, formatDate, parseDate, parseOrder } from "~/lib/utils";
import { LayoutNormal } from "~/ui/Layout";
import { PaginationEllipsis, PaginationEnd, PaginationItem, PaginationItems, PaginationLink, PaginationNext, PaginationPrevious, PaginationRoot } from "~/ui/Pagination";
import { TableBody, TableCaption, TableCell, TableHead, TableHeader, TableRoot, TableRow } from "~/ui/Table";
import { TooltipContent, TooltipRoot, TooltipTrigger } from "~/ui/Tooltip";
import { getEmailsPage } from "./Emails.data";
import { linkVariants } from "~/ui/Link";
import { PageError } from "~/ui/Page";
import { Skeleton } from "~/ui/Skeleton";

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
  const pagination = createPagePagination(() => data()?.pageResult)

  return (
    <LayoutNormal class="max-w-4xl">
      <Shared.Title>Emails</Shared.Title>
      <ErrorBoundary fallback={(e) => <PageError error={e} />}>
        <Suspense fallback={<Skeleton class="h-32" />}>
          <div class="flex justify-between gap-2">
            <Crud.PerPageSelect
              class="w-20"
              perPage={data()?.pageResult?.perPage}
              onChange={(perPage) => setSearchParams({ perPage })}
            />
            <Crud.PageButtons
              class="sm:hidden"
              previousPageDisabled={pagination.previousPageDisabled()}
              previousPage={pagination.previousPage}
              nextPageDisabled={pagination.nextPageDisabled()}
              nextPage={pagination.nextPage}
            />
          </div>
          <PaginationRoot
            page={Math.min(data()?.pageResult?.page || 0, data()?.pageResult?.totalPages || 0)}
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
            <PaginationItems />
            <PaginationEnd>
              <PaginationPrevious />
              <PaginationNext />
            </PaginationEnd>
          </PaginationRoot>
          <TableRoot>
            <TableHeader>
              <TableRow>
                <TableHead>
                  <Crud.SortButton
                    onClick={toggleSort}
                    sort={data()?.sort}
                  >
                    Created At
                  </Crud.SortButton>
                </TableHead>
                <TableHead>Device</TableHead>
                <TableHead>From</TableHead>
                <TableHead>Subject</TableHead>
                <TableHead></TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              <For each={data()?.emails}>
                {v =>
                  <TableRow>
                    <TableCell class="w-0 truncate">
                      {formatDate(parseDate(v.createdAtTime))}
                    </TableCell>
                    <TableCell class="w-0 truncate">
                      <A href={`/devices/${v.deviceId}`} class={linkVariants()}>
                        {v.deviceName}
                      </A>
                    </TableCell>
                    <TableCell class="w-0 truncate">
                      {v.from}
                    </TableCell>
                    <TableCell class="truncate">
                      {v.subject}
                    </TableCell>
                    <TableCell class="bg-background sticky right-0 w-0">
                      <div class="flex justify-end gap-4">
                        <Show when={v.attachmentCount > 0}>
                          <A href={`/emails/${v.id}?tab=attachments`}>
                            <TooltipRoot>
                              <TooltipTrigger class="flex h-full items-center">
                                <RiEditorAttachment2 class="h-5 w-5" />
                              </TooltipTrigger>
                              <TooltipContent>
                                {v.attachmentCount} {Humanize.pluralize(v.attachmentCount, "attachment")}
                              </TooltipContent>
                            </TooltipRoot>
                          </A>
                        </Show>
                        <A href={`/emails/${v.id}`}>
                          <RiArrowsArrowRightLine class="h-5 w-5" />
                        </A>
                      </div>
                    </TableCell>
                  </TableRow>
                }
              </For>
            </TableBody>
            <TableCaption>
              <Crud.PageMetadata pageResult={data()?.pageResult} />
            </TableCaption>
          </TableRoot>
        </Suspense>
      </ErrorBoundary>
    </LayoutNormal>
  )
}
