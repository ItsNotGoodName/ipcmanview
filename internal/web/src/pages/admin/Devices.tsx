import { A, action, createAsync, revalidate, useAction, useNavigate, useSearchParams, useSubmission } from "@solidjs/router";
import { AlertDialogAction, AlertDialogCancel, AlertDialogModal, AlertDialogDescription, AlertDialogFooter, AlertDialogHeader, AlertDialogRoot, AlertDialogTitle, } from "~/ui/AlertDialog";
import { DropdownMenuArrow, DropdownMenuContent, DropdownMenuItem, DropdownMenuPortal, DropdownMenuRoot } from "~/ui/DropdownMenu";
import { AdminDevicesPageSearchParams, getAdminDevicesPage } from "./Devices.data";
import { ErrorBoundary, For, Show, Suspense, createSignal } from "solid-js";
import { catchAsToast, createPagePagination, createRowSelection, createToggleSortField, formatDate, parseDate, } from "~/lib/utils";
import { parseOrder } from "~/lib/utils";
import { TableBody, TableCaption, TableCell, TableHead, TableHeader, TableRoot, TableRow, } from "~/ui/Table";
import { Seperator } from "~/ui/Seperator";
import { useClient } from "~/providers/client";
import { CheckboxControl, CheckboxRoot } from "~/ui/Checkbox";
import { Skeleton } from "~/ui/Skeleton";
import { PageError } from "~/ui/Page";
import { TooltipContent, TooltipRoot, TooltipTrigger } from "~/ui/Tooltip";
import { LayoutNormal } from "~/ui/Layout";
import { SetDeviceDisableReq } from "~/twirp/rpc";
import { Crud } from "~/components/Crud";
import { As } from "@kobalte/core";
import { RiSystemLockLine } from "solid-icons/ri";

const actionDeleteDevice = action((ids: bigint[]) => useClient()
  .admin.deleteDevice({ ids })
  .then(() => revalidate(getAdminDevicesPage.key))
  .catch(catchAsToast)
)

const actionSetDeviceDisable = action((input: SetDeviceDisableReq) => useClient()
  .admin.setDeviceDisable(input)
  .then(() => revalidate(getAdminDevicesPage.key))
  .catch(catchAsToast)
)

export function AdminDevices() {
  const navigate = useNavigate()
  const [searchParams] = useSearchParams<AdminDevicesPageSearchParams>()
  const data = createAsync(() => getAdminDevicesPage({
    page: {
      page: Number(searchParams.page) || 1,
      perPage: Number(searchParams.perPage) || 10
    },
    sort: {
      field: searchParams.sort || "",
      order: parseOrder(searchParams.order)
    },
  }))
  const rowSelection = createRowSelection(() => data()?.items.map(v => v.id) || [])

  // List
  const pagination = createPagePagination(() => data()?.pageResult)
  const toggleSort = createToggleSortField(() => data()?.sort)

  // Delete
  const deleteDeviceSubmission = useSubmission(actionDeleteDevice)
  const deleteDeviceAction = useAction(actionDeleteDevice)
  // Single
  const [openDeleteConfirm, setOpenDeleteConfirm] = createSignal<{ name: string, id: bigint } | undefined>()
  const deleteSubmit = () => deleteDeviceAction([openDeleteConfirm()!.id])
    .then(() => setOpenDeleteConfirm(undefined))
  // Multiple
  const [openDeleteMultipleConfirm, setDeleteDeviceRowSelection] = createSignal(false)
  const deleteMultipleSubmit = () => deleteDeviceAction(rowSelection.selections())
    .then(() => setDeleteDeviceRowSelection(false))

  // Disable/Enable
  const setDisableSubmission = useSubmission(actionSetDeviceDisable)
  const setDisableAction = useAction(actionSetDeviceDisable)
  const setDisableMultipleSubmit = (disable: boolean) => setDisableAction({ items: rowSelection.selections().map(v => ({ id: v, disable })) })
    .then(() => rowSelection.setAll(false))

  return (
    <LayoutNormal class="max-w-4xl">
      <AlertDialogRoot open={openDeleteConfirm() != undefined} onOpenChange={() => setOpenDeleteConfirm(undefined)}>
        <AlertDialogModal>
          <AlertDialogHeader>
            <AlertDialogTitle>Are you sure you wish to delete {openDeleteConfirm()?.name}?</AlertDialogTitle>
          </AlertDialogHeader>
          <AlertDialogFooter>
            <AlertDialogCancel>Cancel</AlertDialogCancel>
            <AlertDialogAction variant="destructive" disabled={deleteDeviceSubmission.pending} onClick={deleteSubmit}>
              Delete
            </AlertDialogAction>
          </AlertDialogFooter>
        </AlertDialogModal>
      </AlertDialogRoot>

      <AlertDialogRoot open={openDeleteMultipleConfirm()} onOpenChange={setDeleteDeviceRowSelection}>
        <AlertDialogModal>
          <AlertDialogHeader>
            <AlertDialogTitle>Are you sure you wish to delete {rowSelection.selections().length} devices?</AlertDialogTitle>
            <AlertDialogDescription>
              <ul>
                <For each={data()?.items}>
                  {(e, index) =>
                    <Show when={rowSelection.rows[index()].checked}>
                      <li>{e.name}</li>
                    </Show>
                  }
                </For>
              </ul>
            </AlertDialogDescription>
          </AlertDialogHeader>
          <AlertDialogFooter>
            <AlertDialogCancel>Cancel</AlertDialogCancel>
            <AlertDialogAction variant="destructive" disabled={deleteDeviceSubmission.pending} onClick={deleteMultipleSubmit}>
              Delete
            </AlertDialogAction>
          </AlertDialogFooter>
        </AlertDialogModal>
      </AlertDialogRoot>

      <div class="text-xl">Devices</div>
      <Seperator />

      <ErrorBoundary fallback={(e: Error) => <PageError error={e} />}>
        <Suspense fallback={<Skeleton class="h-32" />}>
          <div class="flex justify-between gap-2">
            <Crud.PerPageSelect
              class="w-20"
              perPage={data()?.pageResult?.perPage}
              onChange={pagination.setPerPage}
            />
            <Crud.PageButtons
              previousPageDisabled={pagination.previousPageDisabled()}
              previousPage={pagination.previousPage}
              nextPageDisabled={pagination.nextPageDisabled()}
              nextPage={pagination.nextPage}
            />
          </div>
          <TableRoot>
            <TableHeader>
              <TableRow>
                <TableHead>
                  <CheckboxRoot
                    checked={rowSelection.multiple()}
                    indeterminate={rowSelection.indeterminate()}
                    onChange={(v) => rowSelection.setAll(v)}
                  >
                    <CheckboxControl />
                  </CheckboxRoot>
                </TableHead>
                <TableHead>
                  <Crud.SortButton
                    name="name"
                    onClick={toggleSort}
                    sort={data()?.sort}
                  >
                    Name
                  </Crud.SortButton>
                </TableHead>
                <TableHead>
                  <Crud.SortButton
                    name="url"
                    onClick={toggleSort}
                    sort={data()?.sort}
                  >
                    URL
                  </Crud.SortButton>
                </TableHead>
                <TableHead>
                  <Crud.SortButton
                    name="createdAt"
                    onClick={toggleSort}
                    sort={data()?.sort}
                  >
                    Created At
                  </Crud.SortButton>
                </TableHead>
                <Crud.LastTableHead>
                  <DropdownMenuRoot placement="bottom-end">
                    <Crud.MoreDropdownMenuTrigger />
                    <DropdownMenuPortal>
                      <DropdownMenuContent>
                        <DropdownMenuItem asChild>
                          <As component={A} href="./create">Create</As>
                        </DropdownMenuItem>
                        <DropdownMenuItem
                          disabled={rowSelection.selections().length == 0 || setDisableSubmission.pending}
                          onSelect={() => setDisableMultipleSubmit(true)}
                        >
                          Disable
                        </DropdownMenuItem>
                        <DropdownMenuItem
                          disabled={rowSelection.selections().length == 0 || setDisableSubmission.pending}
                          onSelect={() => setDisableMultipleSubmit(false)}
                        >
                          Enable
                        </DropdownMenuItem>
                        <DropdownMenuItem
                          disabled={rowSelection.selections().length == 0 || deleteDeviceSubmission.pending}
                          onSelect={() => setDeleteDeviceRowSelection(true)}
                        >
                          Delete
                        </DropdownMenuItem>
                        <DropdownMenuArrow />
                      </DropdownMenuContent>
                    </DropdownMenuPortal>
                  </DropdownMenuRoot>
                </Crud.LastTableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              <For each={data()?.items}>
                {(item, index) => {
                  const onClick = () => navigate(`./${item.id}`)
                  const toggleDeviceDisable = () => setDisableAction({ items: [{ id: item.id, disable: !item.disabled }] })

                  return (
                    <TableRow>
                      <TableHead>
                        <CheckboxRoot
                          checked={rowSelection.rows[index()]?.checked}
                          onChange={(v) => rowSelection.set(item.id, v)}
                        >
                          <CheckboxControl />
                        </CheckboxRoot>
                      </TableHead>
                      <TableCell class="cursor-pointer select-none" onClick={onClick}>{item.name}</TableCell>
                      <TableCell class="cursor-pointer select-none" onClick={onClick}>{item.url}</TableCell>
                      <TableCell class="cursor-pointer select-none" onClick={onClick}>{formatDate(parseDate(item.createdAtTime))}</TableCell>
                      <Crud.LastTableCell>
                        <Show when={item.disabled}>
                          <TooltipRoot>
                            <TooltipTrigger>
                              <RiSystemLockLine class="h-5 w-5" />
                            </TooltipTrigger>
                            <TooltipContent>
                              Disabled since {formatDate(parseDate(item.disabledAtTime))}
                            </TooltipContent>
                          </TooltipRoot>
                        </Show>
                        <DropdownMenuRoot placement="bottom-end">
                          <Crud.MoreDropdownMenuTrigger />
                          <DropdownMenuPortal>
                            <DropdownMenuContent>
                              <DropdownMenuItem asChild>
                                <As component={A} href={`./${item.id}/update`}>Edit</As>
                              </DropdownMenuItem>
                              <DropdownMenuItem
                                disabled={setDisableSubmission.pending}
                                onSelect={toggleDeviceDisable}
                              >
                                <Show when={item.disabled} fallback={<>Disable</>}>
                                  Enable
                                </Show>
                              </DropdownMenuItem>
                              <DropdownMenuItem
                                disabled={deleteDeviceSubmission.pending}
                                onSelect={() => setOpenDeleteConfirm(item)}
                              >
                                Delete
                              </DropdownMenuItem>
                              <DropdownMenuArrow />
                            </DropdownMenuContent>
                          </DropdownMenuPortal>
                        </DropdownMenuRoot>
                      </Crud.LastTableCell>
                    </TableRow>
                  )
                }}
              </For>
            </TableBody>
            <TableCaption>
              <Crud.Metadata pageResult={data()?.pageResult} />
            </TableCaption>
          </TableRoot>
        </Suspense>
      </ErrorBoundary>
    </LayoutNormal>)
}

