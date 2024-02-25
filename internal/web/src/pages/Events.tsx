import { writeClipboard } from "@solid-primitives/clipboard";
import hljs from "~/lib/hljs"
import { A, createAsync, useSearchParams } from "@solidjs/router";
import { Accessor, ErrorBoundary, For, Suspense, createEffect, createMemo, createSignal, } from "solid-js";
import { Crud } from "~/components/Crud";
import { Shared } from "~/components/Shared";
import { createPagePagination, createToggleSortField, decodeBigInts, encodeBigInts, formatDate, parseDate, parseOrder } from "~/lib/utils";
import { LayoutNormal } from "~/ui/Layout";
import { PaginationEllipsis, PaginationEnd, PaginationItem, PaginationItems, PaginationLink, PaginationNext, PaginationPrevious, PaginationRoot } from "~/ui/Pagination";
import { TableBody, TableCaption, TableCell, TableHead, TableHeader, TableRoot, TableRow } from "~/ui/Table";
import { linkVariants } from "~/ui/Link";
import { PageError } from "~/ui/Page";
import { Skeleton } from "~/ui/Skeleton";
import { getEventsPage } from "./Events.data";
import { RiArrowsArrowDownSLine, RiDocumentClipboardLine, RiSystemFilterLine } from "solid-icons/ri";
import { Button, buttonVariants } from "~/ui/Button";
import { ComboboxContent, ComboboxControl, ComboboxIcon, ComboboxInput, ComboboxItem, ComboboxItemLabel, ComboboxListbox, ComboboxReset, ComboboxRoot, ComboboxState, ComboboxTrigger } from "~/ui/Combobox";
import { getListDevices, getListEventFilters } from "./data";
import { ListDevicesResp_Device } from "~/twirp/rpc";
import { BreadcrumbsItem, BreadcrumbsRoot, } from "~/ui/Breadcrumbs";

export function Events() {
  const [searchParams, setSearchParams] = useSearchParams()

  const filterDeviceIDs: Accessor<bigint[]> = createMemo(() => decodeBigInts(searchParams.device))
  const filterCodes: Accessor<string[]> = createMemo(() => searchParams.code ? JSON.parse(searchParams.code) : [])
  const filterActions: Accessor<string[]> = createMemo(() => searchParams.action ? JSON.parse(searchParams.action) : [])

  const data = createAsync(() => getEventsPage({
    page: {
      page: Number(searchParams.page) || 0,
      perPage: Number(searchParams.perPage) || 0
    },
    sort: {
      field: searchParams.sort || "",
      order: parseOrder(searchParams.order)
    },
    filterDeviceIDs: filterDeviceIDs(),
    filterCodes: filterCodes(),
    filterActions: filterActions(),
  }))
  const listDevices = createAsync(() => getListDevices())
  const listEventFilters = createAsync(() => getListEventFilters())

  const toggleSort = createToggleSortField(() => data()?.sort)
  const pagination = createPagePagination(() => data()?.pageResult)

  const dataOpen = () => Boolean(searchParams.data)
  const setDataOpen = (value: boolean) => setSearchParams({ data: value ? String(value) : "" })

  return (
    <LayoutNormal class="max-w-4xl">
      <Shared.Title>
        <BreadcrumbsRoot>
          <BreadcrumbsItem>
            Events
          </BreadcrumbsItem>
        </BreadcrumbsRoot>
      </Shared.Title>
      <ErrorBoundary fallback={(e) => <PageError error={e} />}>
        <Suspense fallback={<Skeleton class="h-32" />}>
          <div class="flex flex-col gap-2">
            <div class="flex flex-wrap gap-2">
              <Crud.PerPageSelect
                class="hidden w-20 sm:block"
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
                onChange={(value) => setSearchParams({ device: encodeBigInts(value.map(v => v.id)) })}
                itemComponent={props => (
                  <ComboboxItem item={props.item}>
                    <ComboboxItemLabel>{props.item.rawValue.name}</ComboboxItemLabel>
                  </ComboboxItem>
                )}
              >
                <ComboboxControl<ListDevicesResp_Device> aria-label="Device">
                  {state => (
                    <ComboboxTrigger>
                      <ComboboxIcon as={RiSystemFilterLine} class="size-4" />
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
                options={listEventFilters()?.codes || []}
                placeholder="Code"
                value={listEventFilters()?.codes.filter(v => filterCodes().includes(v))}
                onChange={(value) => setSearchParams({ code: value.length != 0 ? JSON.stringify(value) : "" })}
                itemComponent={props => (
                  <ComboboxItem item={props.item}>
                    <ComboboxItemLabel>{props.item.rawValue}</ComboboxItemLabel>
                  </ComboboxItem>
                )}
              >
                <ComboboxControl<string> aria-label="Code">
                  {state => (
                    <ComboboxTrigger>
                      <ComboboxIcon as={RiSystemFilterLine} class="size-4" />
                      Code
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
              <ComboboxRoot<string>
                multiple
                options={listEventFilters()?.actions || []}
                placeholder="Action"
                value={listEventFilters()?.actions.filter(v => filterActions().includes(v))}
                onChange={(value) => setSearchParams({ action: value.length != 0 ? JSON.stringify(value) : "" })}
                itemComponent={props => (
                  <ComboboxItem item={props.item}>
                    <ComboboxItemLabel>{props.item.rawValue}</ComboboxItemLabel>
                  </ComboboxItem>
                )}
              >
                <ComboboxControl<string> aria-label="Action">
                  {state => (
                    <ComboboxTrigger>
                      <ComboboxIcon as={RiSystemFilterLine} class="size-4" />
                      Action
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
              <A class={buttonVariants({ variant: "link" })} href="/events/live">Live</A>
            </div>

            <div class="flex sm:hidden">
              <Crud.PerPageSelect
                class="w-20"
                perPage={data()?.pageResult?.perPage}
                onChange={(perPage) => setSearchParams({ perPage })}
              />

              <Crud.PageButtons
                class="flex-1 justify-end"
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
                  <Button data-expanded={dataOpen()} size="icon" variant="ghost" onClick={() => setDataOpen(!dataOpen())} class="[&[data-expanded=true]>svg]:rotate-180" title="Data">
                    <RiArrowsArrowDownSLine class="h-5 w-5 shrink-0 transition-transform duration-200" />
                  </Button>
                </Crud.LastTableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              <For each={data()?.events}>
                {v => {
                  const [rowDataOpen, setRowDataOpen] = createSignal(dataOpen())
                  createEffect(() => setRowDataOpen(dataOpen()))

                  return (
                    <>
                      <TableRow class="border-b-0">
                        <TableCell>
                          {formatDate(parseDate(v.createdAtTime))}
                        </TableCell>
                        <TableCell>
                          <A href={`/devices/${v.deviceId}`} class={linkVariants()}>
                            {v.deviceName}
                          </A>
                        </TableCell>
                        <TableCell>
                          {v.code}
                        </TableCell>
                        <TableCell>
                          {v.action}
                        </TableCell>
                        <TableCell>
                          {v.index.toString()}
                        </TableCell>
                        <Crud.LastTableCell>
                          <Button data-expanded={rowDataOpen()} size="icon" variant="ghost" onClick={() => setRowDataOpen(!rowDataOpen())} class="[&[data-expanded=true]>svg]:rotate-180" title="Data">
                            <RiArrowsArrowDownSLine class="h-5 w-5 shrink-0 transition-transform duration-200" />
                          </Button>
                        </Crud.LastTableCell>
                      </TableRow>
                      <JSONTableRow colspan={6} expanded={rowDataOpen()} data={v.data} />
                    </>
                  )
                }}
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

export function JSONTableRow(props: { colspan?: number, expanded?: boolean, data: string }) {
  return (
    <tr class="border-b">
      <td colspan={props.colspan} class="p-0">
        <div data-expanded={props.expanded} class="relative h-0 overflow-y-hidden data-[expanded=true]:h-full">
          <Button size="icon" variant="ghost" onClick={() => writeClipboard(props.data)} class="absolute right-4 top-2" title="Copy">
            <RiDocumentClipboardLine class="size-5" />
          </Button>
          <pre><code innerHTML={hljs.highlight(props.data, { language: "json" }).value} class="hljs" /></pre>
        </div>
      </td>
    </tr>
  )
}

export default Events
