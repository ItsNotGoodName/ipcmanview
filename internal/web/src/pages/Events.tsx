import { writeClipboard } from "@solid-primitives/clipboard";
import hljs from "~/lib/hljs"
import { A, createAsync, useSearchParams } from "@solidjs/router";
import { ErrorBoundary, For, Suspense, createEffect, createSignal, } from "solid-js";
import { Crud } from "~/components/Crud";
import { Shared } from "~/components/Shared";
import { createPagePagination, createToggleSortField, formatDate, parseDate, parseOrder } from "~/lib/utils";
import { LayoutNormal } from "~/ui/Layout";
import { PaginationEllipsis, PaginationEnd, PaginationItem, PaginationItems, PaginationLink, PaginationNext, PaginationPrevious, PaginationRoot } from "~/ui/Pagination";
import { TableBody, TableCaption, TableCell, TableHead, TableHeader, TableRoot, TableRow } from "~/ui/Table";
import { linkVariants } from "~/ui/Link";
import { PageError } from "~/ui/Page";
import { Skeleton } from "~/ui/Skeleton";
import { getEventsPage } from "./Events.data";
import { RiArrowsArrowDownSLine, RiDocumentClipboardLine } from "solid-icons/ri";
import { Button } from "~/ui/Button";

export function Events() {
  const [searchParams, setSearchParams] = useSearchParams()
  const data = createAsync(() => getEventsPage({
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

  const [allDataOpen, setAllDataOpen] = createSignal(false)

  return (
    <LayoutNormal class="max-w-4xl">
      <Shared.Title>Events</Shared.Title>
      <ErrorBoundary fallback={(e) => <PageError error={e} />}>
        <Suspense fallback={<Skeleton class="h-32" />}>
          <div>
            <div class="flex justify-end gap-2">
              <Crud.PageButtons
                class="sm:hidden"
                previousPageDisabled={pagination.previousPageDisabled()}
                previousPage={pagination.previousPage}
                nextPageDisabled={pagination.nextPageDisabled()}
                nextPage={pagination.nextPage}
              />
            </div>
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
              <PaginationItems />
              <PaginationEnd>
                <PaginationPrevious />
                <PaginationNext />
              </PaginationEnd>
            </PaginationRoot>
          </div>
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
                <TableHead>Code</TableHead>
                <TableHead>Action</TableHead>
                <TableHead>Index</TableHead>
                <Crud.LastTableHead>
                  <Button data-expanded={allDataOpen()} size="icon" variant="ghost" onClick={() => setAllDataOpen(!allDataOpen())} class="[&[data-expanded=true]>svg]:rotate-180" title="Data">
                    <RiArrowsArrowDownSLine class="h-5 w-5 shrink-0 transition-transform duration-200" />
                  </Button>
                </Crud.LastTableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              <For each={data()?.events}>
                {v => {
                  const [dataOpen, setDataOpen] = createSignal(allDataOpen())
                  createEffect(() => setDataOpen(allDataOpen()))

                  return (
                    <>
                      <TableRow class="border-b-0">
                        <TableCell class="truncate">
                          {formatDate(parseDate(v.createdAtTime))}
                        </TableCell>
                        <TableCell class="truncate">
                          <A href={`/devices/${v.deviceId}`} class={linkVariants()}>
                            {v.deviceName}
                          </A>
                        </TableCell>
                        <TableCell class="truncate">
                          {v.code}
                        </TableCell>
                        <TableCell class="truncate">
                          {v.action}
                        </TableCell>
                        <TableCell class="truncate">
                          {v.index.toString()}
                        </TableCell>
                        <Crud.LastTableCell>
                          <Button data-expanded={dataOpen()} size="icon" variant="ghost" onClick={() => setDataOpen(!dataOpen())} class="[&[data-expanded=true]>svg]:rotate-180" title="Data">
                            <RiArrowsArrowDownSLine class="h-5 w-5 shrink-0 transition-transform duration-200" />
                          </Button>
                        </Crud.LastTableCell>
                      </TableRow>
                      <tr class="border-b">
                        <td colspan={6} class="p-0">
                          <div data-expanded={dataOpen()} class="relative h-0 overflow-y-hidden data-[expanded=true]:h-full">
                            <Button size="icon" variant="ghost" onClick={() => writeClipboard(v.data)} class="absolute right-4 top-2" title="Copy">
                              <RiDocumentClipboardLine class="h-5 w-5" />
                            </Button>
                            <pre innerHTML={hljs.highlight("json", v.data).value} class="hljs p-4" />
                          </div>
                        </td>
                      </tr >
                    </>
                  )
                }}
              </For >
            </TableBody>
            <TableCaption>
              <Crud.PageMetadata pageResult={data()?.pageResult} />
            </TableCaption>
          </TableRoot>
          <div class="flex justify-between gap-2">
            <Crud.PerPageSelect
              class="w-20"
              perPage={data()?.pageResult?.perPage}
              onChange={(perPage) => setSearchParams({ perPage })}
            />
            <Crud.PageButtons
              previousPageDisabled={pagination.previousPageDisabled()}
              previousPage={pagination.previousPage}
              nextPageDisabled={pagination.nextPageDisabled()}
              nextPage={pagination.nextPage}
            />
          </div>
        </Suspense>
      </ErrorBoundary>
    </LayoutNormal >
  )
}

