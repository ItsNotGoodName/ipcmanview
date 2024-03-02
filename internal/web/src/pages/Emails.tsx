import { A, createAsync, useNavigate, useSearchParams } from "@solidjs/router";
import Humanize from "humanize-plus"
import { RiEditorAttachment2, RiSystemFilterLine } from "solid-icons/ri";
import { Accessor, ErrorBoundary, For, Show, Suspense, createMemo } from "solid-js";
import { Crud } from "~/components/Crud";
import { Shared } from "~/components/Shared";
import { encodeQuery, createPagePagination, createToggleSortField, formatDate, parseDate, parseOrder, dotDecode, dotEncode, isTableDataClick } from "~/lib/utils";
import { LayoutNormal } from "~/ui/Layout";
import { PaginationEllipsis, PaginationEnd, PaginationItem, PaginationItems, PaginationLink, PaginationNext, PaginationPrevious, PaginationRoot } from "~/ui/Pagination";
import { TableBody, TableCaption, TableCell, TableHead, TableHeader, TableRoot, TableRow } from "~/ui/Table";
import { TooltipArrow, TooltipContent, TooltipRoot, TooltipTrigger } from "~/ui/Tooltip";
import { getEmailsPage, withEmailPageQuery } from "./Emails.data";
import { linkVariants } from "~/ui/Link";
import { PageError } from "~/ui/Page";
import { Skeleton } from "~/ui/Skeleton";
import { BreadcrumbsItem, BreadcrumbsRoot } from "~/ui/Breadcrumbs";
import { getListEmailAlarmEvents } from "./data";
import { ComboboxContent, ComboboxControl, ComboboxIcon, ComboboxInput, ComboboxItem, ComboboxItemLabel, ComboboxListbox, ComboboxReset, ComboboxRoot, ComboboxState, ComboboxTrigger } from "~/ui/Combobox";
import { DeviceFilterCombobox } from "~/components/DeviceFilterCombobox"

export function Emails() {
  const navigate = useNavigate()
  const [searchParams, setSearchParams] = useSearchParams()

  const filterDeviceIDs: Accessor<string[]> = createMemo(() => dotDecode(searchParams.device))
  const setFilterDeviceIDs = (value: string[]) => setSearchParams({ device: dotEncode(value) })
  const filterAlarmEvents: Accessor<string[]> = createMemo(() => searchParams.alarmEvent ? JSON.parse(searchParams.alarmEvent) : [])

  const data = createAsync(() => getEmailsPage({
    page: {
      page: Number(searchParams.page) || 0,
      perPage: Number(searchParams.perPage) || 0
    },
    sort: {
      field: searchParams.sort || "",
      order: parseOrder(searchParams.order)
    },
    filterDeviceIDs: filterDeviceIDs(),
    filterAlarmEvents: filterAlarmEvents()
  }))
  const listEmailAlarmEvents = createAsync(() => getListEmailAlarmEvents())

  const toggleSort = createToggleSortField(() => data()?.sort)
  const pagination = createPagePagination(() => data()?.pageResult)
  const query = (init?: string) => encodeQuery(withEmailPageQuery(new URLSearchParams(init), searchParams))

  return (
    <LayoutNormal class="max-w-4xl">
      <Shared.Title>
        <BreadcrumbsRoot>
          <BreadcrumbsItem>
            Emails
          </BreadcrumbsItem>
        </BreadcrumbsRoot>
      </Shared.Title>
      <div class="flex flex-col gap-2">
        <ErrorBoundary fallback={(e) => <PageError error={e} />}>
          <Suspense fallback={<Skeleton class="h-32" />}>
            <div class="flex flex-wrap gap-2">
              <Crud.PerPageSelect
                perPage={data()?.pageResult?.perPage}
                onChange={(perPage) => setSearchParams({ perPage })}
                class="hidden w-20 sm:block"
              />
              <DeviceFilterCombobox deviceIDs={filterDeviceIDs()} setDeviceIDs={setFilterDeviceIDs} />
              <ComboboxRoot<string>
                multiple
                options={listEmailAlarmEvents() || []}
                placeholder="Alarm Event"
                value={listEmailAlarmEvents()?.filter(v => filterAlarmEvents().includes(v))}
                onChange={(value) => setSearchParams({ alarmEvent: value.length != 0 ? JSON.stringify(value) : "" })}
                itemComponent={props => (
                  <ComboboxItem item={props.item}>
                    <ComboboxItemLabel>{props.item.rawValue}</ComboboxItemLabel>
                  </ComboboxItem>
                )}
              >
                <ComboboxControl<string> aria-label="Alarm Event">
                  {state => (
                    <ComboboxTrigger>
                      <ComboboxIcon as={RiSystemFilterLine} class="size-4" />
                      Alarm Event
                      <ComboboxState state={state} />
                      <ComboboxReset state={state} class="size-4" />
                    </ComboboxTrigger>
                  )}
                </ComboboxControl>
                <ComboboxContent>
                  <ComboboxInput />
                  <ComboboxListbox />
                </ComboboxContent>
              </ComboboxRoot>
            </div>

            <div class="flex sm:hidden">
              <Crud.PerPageSelect
                perPage={data()?.pageResult?.perPage}
                onChange={(perPage) => setSearchParams({ perPage })}
                class="w-20"
              />

              <Crud.PageButtons
                previousPageDisabled={pagination.previousPageDisabled()}
                previousPage={pagination.previousPage}
                nextPageDisabled={pagination.nextPageDisabled()}
                nextPage={pagination.nextPage}
                class="flex-1 justify-end"
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
                    <TableRow
                      onClick={(t) => isTableDataClick(t) && navigate(`/emails/${v.id}${query()}`)}
                      class="[&>td]:cursor-pointer"
                    >
                      <TableCell>
                        {formatDate(parseDate(v.createdAtTime))}
                      </TableCell>
                      <TableCell>
                        <A href={`/devices/${v.deviceId}`} class={linkVariants()}>
                          {v.deviceName}
                        </A>
                      </TableCell>
                      <TableCell>
                        {v.alarmEvent}
                      </TableCell>
                      <TableCell>
                        {v.from}
                      </TableCell>
                      <TableCell>
                        {v.subject}
                      </TableCell>
                      <Crud.LastTableCell>
                        <Show when={v.attachmentCount > 0}>
                          <A href={`/emails/${v.id}${query("?tab=attachments")}`}>
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
          </Suspense>
        </ErrorBoundary>
      </div>
    </LayoutNormal>
  )
}

export default Emails
