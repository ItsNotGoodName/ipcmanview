import { LayoutNormal } from "~/ui/Layout";
import { Shared } from "~/components/Shared";
import { BreadcrumbsItem, BreadcrumbsRoot } from "~/ui/Breadcrumbs";
import { ComboboxContent, ComboboxControl, ComboboxIcon, ComboboxInput, ComboboxItem, ComboboxItemLabel, ComboboxListbox, ComboboxReset, ComboboxRoot, ComboboxState, ComboboxTrigger } from "~/ui/Combobox";
import { A, createAsync, useSearchParams } from "@solidjs/router";
import { getFilesPage } from "./Files.data";
import { RiMediaImageLine, RiMediaVideoLine, RiSystemDownloadLine, RiSystemFilterLine } from "solid-icons/ri";
import { createPagePagination, dotDecode, dotEncode, encodeOrder, formatDate, parseDate, parseOrder, } from "~/lib/utils";
import { DeviceFilterCombobox } from "~/components/DeviceFilterCombobox";
import { For, Show, createMemo, createSignal } from "solid-js";
import { Crud } from "~/components/Crud";
import { GetFileMonthCountResp_Month, Order } from "~/twirp/rpc";
import { PaginationEllipsis, PaginationEnd, PaginationItem, PaginationItems, PaginationLink, PaginationNext, PaginationPrevious, PaginationRoot } from "~/ui/Pagination";
import { Seperator } from "~/ui/Seperator";
import { Image } from "@kobalte/core";
import { getFileMonthCount } from "./data";
import { SwitchControl, SwitchLabel, SwitchRoot } from "~/ui/Switch";

export function Files() {
  const [searchParams, setSearchParams] = useSearchParams()

  const filterDeviceIDs = createMemo(() => dotDecode(searchParams.device))
  const setFilterDeviceIDs = (value: string[]) => setSearchParams({ page: 1, device: dotEncode(value) })
  const filterMonthID = () => searchParams.month ?? ""
  const setFilterMonths = (value: GetFileMonthCountResp_Month) => setSearchParams({ page: 1, month: value?.monthId })
  const order = () => parseOrder(searchParams.order)
  const orderAscending = () => order() == Order.ASC ? "checked" : ""
  const setOrderAscending = (value: boolean) => setSearchParams({ order: value ? encodeOrder(Order.ASC) : undefined })

  const data = createAsync(() => getFilesPage({
    page: {
      page: Number(searchParams.page) || 0,
      perPage: Number(searchParams.perPage) || 0
    },
    filterDeviceIDs: filterDeviceIDs(),
    filterMonthID: filterMonthID(),
    order: order()
  }))
  const fileMonthCount = createAsync(() => getFileMonthCount(filterDeviceIDs()))

  const pagination = createPagePagination(() => data()?.pageResult)

  return (
    <LayoutNormal>
      <Shared.Title>
        <BreadcrumbsRoot>
          <BreadcrumbsItem>
            Files
          </BreadcrumbsItem>
        </BreadcrumbsRoot>
      </Shared.Title>
      <div class="flex flex-col gap-2">
        <div class="flex flex-wrap gap-2">
          <Crud.PerPageSelect
            perPage={data()?.pageResult?.perPage}
            onChange={(perPage) => setSearchParams({ perPage })}
            class="hidden w-20 sm:block"
          />
          <DeviceFilterCombobox deviceIDs={filterDeviceIDs()} setDeviceIDs={setFilterDeviceIDs} />
          <ComboboxRoot<GetFileMonthCountResp_Month>
            options={fileMonthCount()?.months || []}
            optionTextValue="monthId"
            optionValue="monthId"
            optionLabel="monthId"
            placeholder="Months"
            value={{ monthId: filterMonthID(), count: 0 }}
            onChange={setFilterMonths}
            itemComponent={props => (
              <ComboboxItem item={props.item}>
                <ComboboxItemLabel>
                  {props.item.rawValue.monthId}
                </ComboboxItemLabel>
                <div class="text-muted-foreground ml-auto">
                  {props.item.rawValue.count}
                </div>
              </ComboboxItem>
            )}
          >
            <ComboboxControl<GetFileMonthCountResp_Month> aria-label="Month">
              {state => (
                <ComboboxTrigger>
                  <ComboboxIcon as={RiSystemFilterLine} class="size-4" />
                  Month
                  <ComboboxState state={state} getOptionString={v => v.monthId} />
                  <ComboboxReset state={state} class="size-4" />
                </ComboboxTrigger>
              )}
            </ComboboxControl>
            <ComboboxContent>
              <ComboboxInput />
              <ComboboxListbox />
            </ComboboxContent>
          </ComboboxRoot>

          <SwitchRoot
            value={orderAscending()}
            onChange={setOrderAscending}
            class="flex items-center gap-2">
            <SwitchLabel>Ascending</SwitchLabel>
            <SwitchControl />
          </SwitchRoot>
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

        <Crud.PageMetadata pageResult={data()?.pageResult} />

        <div class="grid grid-cols-2 gap-4 sm:grid-cols-4 xl:grid-cols-6 2xl:grid-cols-8">
          <For each={data()?.files}>
            {v => {
              const [url, setUrl] = createSignal<string>()
              const isImage = () => v.type == "jpg"
              const loadUrl = () => isImage() && setUrl(v.url)
              const srcUrl = () => url() ?? v.thumbnailUrl

              return (
                <div>
                  <div class="flex flex-col rounded-b border transition-all">
                    <Image.Root onClick={loadUrl} class="mx-auto w-full cursor-pointer">
                      <Image.Img src={srcUrl()} class="h-full w-full object-contain" />
                      <Image.Fallback>
                        <Show when={isImage()} fallback={
                          <RiMediaVideoLine class="h-full w-full object-contain" />
                        }>
                          <RiMediaImageLine class="h-full w-full object-contain" />
                        </Show>
                      </Image.Fallback>
                    </Image.Root>
                    <Seperator />
                    <div class="flex items-center justify-between gap-2 p-2">
                      <div class="flex flex-col text-sm">
                        <div>
                          <A href={`/files/${v.id}`}>{formatDate(parseDate(v.startTime))}</A>
                        </div>
                        <div>
                          <A href={`/devices/${v.deviceId}`}>{v.deviceName}</A>
                        </div>
                      </div>
                      <a href={v.url} target="_blank" title="Download">
                        <RiSystemDownloadLine class="h-5 w-5" />
                      </a>
                    </div>
                  </div>
                </div>
              )
            }
            }
          </For>
        </div>
      </div>
    </LayoutNormal>
  )
}

export default Files
