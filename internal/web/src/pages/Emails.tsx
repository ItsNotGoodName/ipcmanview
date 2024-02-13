import { A, createAsync, useNavigate, useSearchParams } from "@solidjs/router";
import Humanize from "humanize-plus"
import { RiEditorAttachment2 } from "solid-icons/ri";
import { ErrorBoundary, For, Show, Suspense } from "solid-js";
import { Crud } from "~/components/Crud";
import { Shared } from "~/components/Shared";
import { createPagePagination, createToggleSortField, formatDate, parseDate, parseOrder } from "~/lib/utils";
import { LayoutNormal } from "~/ui/Layout";
import { PaginationEllipsis, PaginationEnd, PaginationItem, PaginationItems, PaginationLink, PaginationNext, PaginationPrevious, PaginationRoot } from "~/ui/Pagination";
import { TableBody, TableCaption, TableCell, TableHead, TableHeader, TableRoot, TableRow } from "~/ui/Table";
import { TooltipArrow, TooltipContent, TooltipRoot, TooltipTrigger } from "~/ui/Tooltip";
import { getEmailsPage } from "./Emails.data";
import { linkVariants } from "~/ui/Link";
import { PageError } from "~/ui/Page";
import { Skeleton } from "~/ui/Skeleton";
import { BreadcrumbsItem, BreadcrumbsRoot } from "~/ui/Breadcrumbs";

export function Emails() {
  const navigate = useNavigate()
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
      <Shared.Title>
        <BreadcrumbsRoot>
          <BreadcrumbsItem>
            Emails
          </BreadcrumbsItem>
        </BreadcrumbsRoot>
      </Shared.Title>
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
                <TableHead>Alarm Event</TableHead>
                <TableHead>From</TableHead>
                <TableHead>Subject</TableHead>
                <TableHead></TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              <For each={data()?.emails}>
                {v =>
                  <TableRow onClick={(t) => (t as any).srcElement.tagName == "TD" && navigate(`/emails/${v.id}`)} class="cursor-pointer">
                    <TableCell class="truncate">
                      {formatDate(parseDate(v.createdAtTime))}
                    </TableCell>
                    <TableCell class="truncate">
                      <A href={`/devices/${v.deviceId}`} class={linkVariants()}>
                        {v.deviceName}
                      </A>
                    </TableCell>
                    <TableCell class="truncate">
                      {v.alarmEvent}
                    </TableCell>
                    <TableCell class="truncate">
                      {v.from}
                    </TableCell>
                    <TableCell class="truncate">
                      {v.subject}
                    </TableCell>
                    <Crud.LastTableCell>
                      <Show when={v.attachmentCount > 0}>
                        <A href={`/emails/${v.id}?tab=attachments`}>
                          <TooltipRoot>
                            <TooltipTrigger class="flex h-full items-center">
                              <RiEditorAttachment2 class="h-4 w-4" />
                            </TooltipTrigger>
                            <TooltipContent>
                              <TooltipArrow />
                              {v.attachmentCount} {Humanize.pluralize(v.attachmentCount, "attachment")}
                            </TooltipContent>
                          </TooltipRoot>
                        </A>
                      </Show>
                    </Crud.LastTableCell>
                  </TableRow>
                }
              </For>
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
    </LayoutNormal>
  )
}
