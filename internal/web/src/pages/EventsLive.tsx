import { A, createAsync, useSearchParams } from "@solidjs/router";
import { ErrorBoundary, For, Suspense, createEffect, createSignal, } from "solid-js";
import { Crud } from "~/components/Crud";
import { Shared } from "~/components/Shared";
import { formatDate, } from "~/lib/utils";
import { LayoutNormal } from "~/ui/Layout";
import { TableBody, TableCell, TableHead, TableHeader, TableRoot, TableRow } from "~/ui/Table";
import { linkVariants } from "~/ui/Link";
import { PageError } from "~/ui/Page";
import { Skeleton } from "~/ui/Skeleton";
import { RiArrowsArrowDownSLine } from "solid-icons/ri";
import { Button } from "~/ui/Button";
import { getListDevices } from "./data";
import { BreadcrumbsItem, BreadcrumbsLink, BreadcrumbsRoot, BreadcrumbsSeparator } from "~/ui/Breadcrumbs";
import { useBus } from "~/providers/bus";
import { createDate, createTimeAgo } from "@solid-primitives/date";
import { TooltipArrow, TooltipContent, TooltipRoot, TooltipTrigger } from "~/ui/Tooltip";
import { DahuaEvent } from "~/lib/models.gen";
import { JSONTableRow } from "./Events";

export function EventsLive() {
  const bus = useBus()
  const [searchParams, setSearchParams] = useSearchParams()

  const listDevices = createAsync(() => getListDevices())

  const dataOpen = () => Boolean(searchParams.data)
  const setDataOpen = (value: boolean) => setSearchParams({ data: value ? String(value) : "" })

  const [events, setEvents] = createSignal<DahuaEvent[]>([])
  bus.dahuaEvent.listen((e) => setEvents((prev) => [e, ...prev]))

  return (
    <LayoutNormal class="max-w-4xl">
      <Shared.Title>
        <BreadcrumbsRoot>
          <BreadcrumbsItem>
            <BreadcrumbsLink as={A} href="/events">
              Events
            </BreadcrumbsLink>
            <BreadcrumbsSeparator />
          </BreadcrumbsItem>
          <BreadcrumbsItem>
            Live
          </BreadcrumbsItem>
        </BreadcrumbsRoot>
      </Shared.Title>
      <ErrorBoundary fallback={(e) => <PageError error={e} />}>
        <Suspense fallback={<Skeleton class="h-32" />}>
          <TableRoot>
            <TableHeader>
              <TableRow>
                <TableHead>Created At</TableHead>
                <TableHead>Device</TableHead>
                <TableHead>Code</TableHead>
                <TableHead>Action</TableHead>
                <TableHead>Index</TableHead>
                <Crud.LastTableHead>
                  <Button data-expanded={dataOpen()} onClick={() => setDataOpen(!dataOpen())} title="Data" size="icon" variant="ghost" class="[&[data-expanded=true]>svg]:rotate-180">
                    <RiArrowsArrowDownSLine class="h-5 w-5 shrink-0 transition-transform duration-200" />
                  </Button>
                </Crud.LastTableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              <For each={events()}>
                {v => {
                  const [rowDataOpen, setRowDataOpen] = createSignal(dataOpen())
                  createEffect(() => setRowDataOpen(dataOpen()))

                  const [createdAt] = createDate(() => v.created_at);
                  const [createdAtAgo] = createTimeAgo(createdAt);

                  return (
                    <>
                      <TableRow class="border-b-0">
                        <TableCell>
                          <TooltipRoot>
                            <TooltipTrigger>{createdAtAgo()}</TooltipTrigger>
                            <TooltipContent>
                              <TooltipArrow />
                              {formatDate(createdAt())}
                            </TooltipContent>
                          </TooltipRoot>
                        </TableCell>
                        <TableCell>
                          <A href={`/devices/${v.device_id}`} class={linkVariants()}>
                            {listDevices()?.find((d) => d.id == String(v.device_id))?.name}
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
                          <Button data-expanded={rowDataOpen()} onClick={() => setRowDataOpen(!rowDataOpen())} title="Data" size="icon" variant="ghost" class="[&[data-expanded=true]>svg]:rotate-180">
                            <RiArrowsArrowDownSLine class="h-5 w-5 shrink-0 transition-transform duration-200" />
                          </Button>
                        </Crud.LastTableCell>
                      </TableRow>
                      <JSONTableRow colspan={6} expanded={rowDataOpen()} data={JSON.stringify(v.data, null, 2)} />
                    </>
                  )
                }}
              </For>
            </TableBody>
          </TableRoot>
        </Suspense>
      </ErrorBoundary>
    </LayoutNormal>
  )
}

export default EventsLive
