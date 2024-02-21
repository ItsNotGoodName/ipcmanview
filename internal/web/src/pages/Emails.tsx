import { A, createAsync, useNavigate, useSearchParams } from "@solidjs/router";
import Humanize from "humanize-plus"
import { RiEditorAttachment2, RiSystemAddCircleLine } from "solid-icons/ri";
import { Accessor, ErrorBoundary, For, Show, Suspense, createMemo } from "solid-js";
import { Crud } from "~/components/Crud";
import { Shared } from "~/components/Shared";
import { encodeQuery, createPagePagination, createToggleSortField, formatDate, parseDate, parseOrder } from "~/lib/utils";
import { LayoutNormal } from "~/ui/Layout";
import { PaginationEllipsis, PaginationEnd, PaginationItem, PaginationItems, PaginationLink, PaginationNext, PaginationPrevious, PaginationRoot } from "~/ui/Pagination";
import { TableBody, TableCaption, TableCell, TableHead, TableHeader, TableRoot, TableRow } from "~/ui/Table";
import { TooltipArrow, TooltipContent, TooltipRoot, TooltipTrigger } from "~/ui/Tooltip";
import { getEmailsPage, withEmailPageQuery } from "./Emails.data";
import { linkVariants } from "~/ui/Link";
import { PageError } from "~/ui/Page";
import { Skeleton } from "~/ui/Skeleton";
import { BreadcrumbsItem, BreadcrumbsRoot } from "~/ui/Breadcrumbs";
import { getlistDevices, getListEmailAlarmEvents } from "./data";
import { ComboboxContent, ComboboxControl, ComboboxIcon, ComboboxInput, ComboboxItem, ComboboxItemLabel, ComboboxListbox, ComboboxReset, ComboboxRoot, ComboboxState, ComboboxTrigger } from "~/ui/Combobox";
import { ListDevicesResp_Device } from "~/twirp/rpc";

export function Emails() {
  const navigate = useNavigate()
  const [searchParams, setSearchParams] = useSearchParams()
  const filterDeviceIDs: Accessor<bigint[]> = createMemo(() => searchParams.device ? searchParams.device.split('.').map(v => BigInt(v)) : [])
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
  const listDevices = createAsync(() => getlistDevices())
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
      <ErrorBoundary fallback={(e) => <PageError error={e} />}>
        <Suspense fallback={<Skeleton class="h-32" />}>
          <div class="flex flex-col gap-2">
            <div class="flex flex-wrap gap-2">
              <Crud.PerPageSelect
                class="w-20"
                perPage={data()?.pageResult?.perPage}
                onChange={(perPage) => setSearchParams({ perPage })}
              />
              <ComboboxRoot<ListDevicesResp_Device>
                multiple
                optionValue="id"
                optionTextValue="name"
                optionLabel="name"
                options={listDevices() || []}
                placeholder="Device"
                value={listDevices()?.filter(v => filterDeviceIDs().includes(v.id))}
                onChange={(value) => setSearchParams({ device: value.map(v => v.id).join('.') })}
                itemComponent={props => (
                  <ComboboxItem item={props.item}>
                    <ComboboxItemLabel>{props.item.rawValue.name}</ComboboxItemLabel>
                  </ComboboxItem>
                )}
              >
                <ComboboxControl<ListDevicesResp_Device> aria-label="Device">
                  {state => (
                    <ComboboxTrigger>
                      <ComboboxIcon as={RiSystemAddCircleLine} class="size-4" />
                      Device
                      <ComboboxState state={state} optionToString={(option) => option.name} />
                      <ComboboxReset state={state} class="size-4" />
                    </ComboboxTrigger>
                  )}
                </ComboboxControl>
                <ComboboxContent>
                  <ComboboxInput />
                  <ComboboxListbox />
                </ComboboxContent>
              </ComboboxRoot>
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
                      <ComboboxIcon as={RiSystemAddCircleLine} class="size-4" />
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
              <Crud.PageButtons
                class="flex flex-1 items-center justify-end sm:hidden"
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
                  <TableRow onClick={(t) => (t as any).srcElement.tagName == "TD" && navigate(`/emails/${v.id}${query()}`)} class="cursor-pointer">
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
    </LayoutNormal>
  )
}
