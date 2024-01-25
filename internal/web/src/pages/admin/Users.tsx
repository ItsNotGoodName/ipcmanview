import { createAsync, useNavigate, useSearchParams, } from "@solidjs/router";
import { ErrorBoundary, For, Show, Suspense, createEffect, } from "solid-js";
import { RiArrowsArrowLeftSLine, RiArrowsArrowRightSLine, RiSystemLockLine, RiUserFacesAdminLine, } from "solid-icons/ri";
import { Button } from "~/ui/Button";
import { SelectContent, SelectItem, SelectListbox, SelectRoot, SelectTrigger, SelectValue } from "~/ui/Select";
import { formatDate, parseDate, } from "~/lib/utils";
import { Order, } from "~/twirp/rpc";
import { encodeOrder, nextOrder, parseOrder } from "~/lib/order";
import { TableBody, TableCaption, TableCell, TableHead, TableHeader, TableMetadata, TableRoot, TableRow, TableSortButton } from "~/ui/Table";
import { Seperator } from "~/ui/Seperator";
import { Skeleton } from "~/ui/Skeleton";
import { PageError } from "~/ui/Page";
import { TooltipContent, TooltipRoot, TooltipTrigger } from "~/ui/Tooltip";
import { AdminUsersPageSearchParams, getAdminUsersPage } from "./Users.data";
import { unwrap } from "solid-js/store";

export function AdminUsers() {
  const navigate = useNavigate()
  const [searchParams, setSearchParams] = useSearchParams<AdminUsersPageSearchParams>()
  const data = createAsync(() => getAdminUsersPage({
    page: {
      page: Number(searchParams.page) || 1,
      perPage: Number(searchParams.perPage) || 10
    },
    sort: {
      field: searchParams.sort || "",
      order: parseOrder(searchParams.order)
    },
  }))

  const previousDisabled = () => data()?.pageResult?.previousPage == data()?.pageResult?.page
  const previous = () => !previousDisabled() && setSearchParams({ page: data()?.pageResult?.previousPage.toString() } as AdminUsersPageSearchParams)
  const nextDisabled = () => data()?.pageResult?.nextPage == data()?.pageResult?.page
  const next = () => !nextDisabled() && setSearchParams({ page: data()?.pageResult?.nextPage.toString() } as AdminUsersPageSearchParams)
  const toggleSort = (value: string) => {
    if (value == data()?.sort?.field) {
      const order = nextOrder(data()?.sort?.order ?? Order.ORDER_UNSPECIFIED)

      if (order == Order.ORDER_UNSPECIFIED) {
        return setSearchParams({ sort: undefined, order: undefined })
      }

      return setSearchParams({ sort: value, order: encodeOrder(order) } as AdminUsersPageSearchParams)
    }

    return setSearchParams({ sort: value, order: encodeOrder(Order.DESC) } as AdminUsersPageSearchParams)
  }


  return (
    <div class="flex justify-center p-4">
      <div class="flex w-full max-w-4xl flex-col gap-2">
        <div class="text-xl">Users</div>
        <Seperator />
        <ErrorBoundary fallback={(e: Error) => <PageError error={e} />}>
          <Suspense fallback={<Skeleton class="h-32" />}>
            <div class="flex justify-between gap-2">
              <SelectRoot
                class="w-20"
                value={data()?.pageResult?.perPage}
                onChange={(value) => value && setSearchParams({ page: 1, perPage: value })}
                options={[10, 25, 50, 100]}
                itemComponent={props => (
                  <SelectItem item={props.item}>
                    {props.item.rawValue}
                  </SelectItem>
                )}
              >
                <SelectTrigger aria-label="Per page">
                  <SelectValue<number>>
                    {state => state.selectedOption()}
                  </SelectValue>
                </SelectTrigger>
                <SelectContent>
                  <SelectListbox />
                </SelectContent>
              </SelectRoot>
              <div class="flex gap-2">
                <Button
                  title="Previous"
                  size="icon"
                  disabled={previousDisabled()}
                  onClick={previous}
                >
                  <RiArrowsArrowLeftSLine class="h-6 w-6" />
                </Button>
                <Button
                  title="Next"
                  size="icon"
                  disabled={nextDisabled()}
                  onClick={next}
                >
                  <RiArrowsArrowRightSLine class="h-6 w-6" />
                </Button>
              </div>
            </div>
            <TableRoot>
              <TableHeader>
                <tr class="border-b">
                  <TableHead>
                    <TableSortButton
                      name="username"
                      onClick={toggleSort}
                      sort={data()?.sort}
                    >
                      Username
                    </TableSortButton>
                  </TableHead>
                  <TableHead class="w-full">
                    <TableSortButton
                      name="email"
                      onClick={toggleSort}
                      sort={data()?.sort}
                    >
                      Email
                    </TableSortButton>
                  </TableHead>
                  <TableHead>
                    <TableSortButton
                      name="createdAt"
                      onClick={toggleSort}
                      sort={data()?.sort}
                    >
                      Created At
                    </TableSortButton>
                  </TableHead>
                  <TableHead></TableHead>
                </tr>
              </TableHeader>
              <TableBody>
                <For each={data()?.items}>
                  {(item) => {
                    const onClick = () => navigate(`./${item.id}`)

                    return (
                      <TableRow class="">
                        <TableCell onClick={onClick} class="cursor-pointer select-none">{item.username}</TableCell>
                        <TableCell onClick={onClick} class="cursor-pointer select-none">{item.email}</TableCell>
                        <TableCell onClick={onClick} class="text-nowrap cursor-pointer select-none whitespace-nowrap">{formatDate(parseDate(item.createdAtTime))}</TableCell>
                        <TableCell class="py-0">
                          <div class="flex gap-2">
                            <Show when={item.admin}>
                              <TooltipRoot>
                                <TooltipTrigger class="p-1">
                                  <RiUserFacesAdminLine class="h-5 w-5" />
                                </TooltipTrigger>
                                <TooltipContent>
                                  Admin
                                </TooltipContent>
                              </TooltipRoot>
                            </Show>
                            <Show when={item.disabled}>
                              <TooltipRoot>
                                <TooltipTrigger class="p-1">
                                  <RiSystemLockLine class="h-5 w-5" />
                                </TooltipTrigger>
                                <TooltipContent>
                                  Disabled since {formatDate(parseDate(item.disabledAtTime))}
                                </TooltipContent>
                              </TooltipRoot>
                            </Show>
                          </div>
                        </TableCell>
                      </TableRow>
                    )
                  }}
                </For>
              </TableBody>
              <TableCaption>
                <TableMetadata pageResult={data()?.pageResult} />
              </TableCaption>
            </TableRoot>
          </Suspense>
        </ErrorBoundary>
      </div>
    </div>
  )
}

