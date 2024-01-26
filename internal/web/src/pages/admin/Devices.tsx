import { action, createAsync, revalidate, useAction, useNavigate, useSearchParams, useSubmission } from "@solidjs/router";
import { AlertDialogAction, AlertDialogCancel, AlertDialogContent, AlertDialogFooter, AlertDialogHeader, AlertDialogRoot, AlertDialogTitle, } from "~/ui/AlertDialog";
import { DropdownMenuArrow, DropdownMenuContent, DropdownMenuItem, DropdownMenuPortal, DropdownMenuRoot, DropdownMenuTrigger } from "~/ui/DropdownMenu";
import { AdminDevicesPageSearchParams, getAdminDevicesPage, getDevice } from "./Devices.data";
import { ErrorBoundary, For, Show, Suspense, batch, createSignal } from "solid-js";
import { RiArrowsArrowLeftSLine, RiArrowsArrowRightSLine, RiSystemLockLine, RiSystemMore2Line, } from "solid-icons/ri";
import { Button } from "~/ui/Button";
import { SelectContent, SelectItem, SelectListbox, SelectRoot, SelectTrigger, SelectValue } from "~/ui/Select";
import { catchAsToast, createPagePagination, createRowSelection, formatDate, parseDate, syncForm, throwAsFormError } from "~/lib/utils";
import { encodeOrder, toggleSortField, parseOrder } from "~/lib/utils";
import { TableBody, TableCaption, TableCell, TableHead, TableHeader, TableMetadata, TableRoot, TableRow, TableSortButton } from "~/ui/Table";
import { Seperator } from "~/ui/Seperator";
import { useClient } from "~/providers/client";
import { createForm, required, reset } from "@modular-forms/solid";
import { FieldControl, FieldLabel, FieldMessage, FieldRoot, FormMessage } from "~/ui/Form";
import { Input } from "~/ui/Input";
import { Textarea } from "~/ui/Textarea";
import { DialogCloseButton, DialogContent, DialogHeader, DialogOverlay, DialogPortal, DialogRoot, DialogTitle, } from "~/ui/Dialog";
import { CheckboxControl, CheckboxInput, CheckboxLabel, CheckboxRoot } from "~/ui/Checkbox";
import { Skeleton } from "~/ui/Skeleton";
import { PageError } from "~/ui/Page";
import { TooltipContent, TooltipRoot, TooltipTrigger } from "~/ui/Tooltip";
import { defaultPerPageOptions } from "~/lib/utils";
import { LayoutNormal } from "~/ui/Layout";
import { SetDeviceDisableReq } from "~/twirp/rpc";

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
  const [searchParams, setSearchParams] = useSearchParams<AdminDevicesPageSearchParams>()
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
  const toggleSort = (field: string) => {
    const sort = toggleSortField(data()?.sort, field)
    return setSearchParams({ sort: sort.field, order: encodeOrder(sort.order) } as AdminDevicesPageSearchParams)
  }

  // Create
  const [createFormDialog, setCreateFormDialog] = createSignal(false);

  // Update
  const [updateDeviceFormDialog, setUpdateDeviceFormDialog] = createSignal<bigint>(BigInt(0))

  // Delete
  const deleteDeviceSubmission = useSubmission(actionDeleteDevice)
  const deleteDeviceAction = useAction(actionDeleteDevice)
  // Single
  const [deleteDeviceSelection, setDeleteDeviceSelection] = createSignal<{ name: string, id: bigint } | undefined>()
  const deleteDeviceBySelection = () => deleteDeviceAction([deleteDeviceSelection()!.id])
    .then(() => setDeleteDeviceSelection(undefined))
  // Multiple
  const [deleteDeviceRowSelector, setDeleteDeviceRowSelector] = createSignal(false)
  const deleteDeviceByRowSelector = () => deleteDeviceAction(rowSelection.selections())
    .then(() => setDeleteDeviceRowSelector(false))

  // Disable/Enable
  const setDeviceDisableSubmission = useSubmission(actionSetDeviceDisable)
  const setDeviceDisable = useAction(actionSetDeviceDisable)
  const setDeviceDisableByRowSelector = (disable: boolean) => setDeviceDisable({ items: rowSelection.selections().map(v => ({ id: v, disable })) })
    .then(() => rowSelection.setAll(false))
  const setDeviceDisableDisabled = (disable: boolean) => {
    for (let i = 0; i < rowSelection.rows.length; i++) {
      if (rowSelection.rows[i].checked && (disable != data()?.items[i].disabled))
        return false;
    }
    return true
  }

  return (
    <LayoutNormal>
      <DialogRoot open={createFormDialog()} onOpenChange={setCreateFormDialog}>
        <DialogPortal>
          <DialogOverlay />
          <DialogContent>
            <DialogHeader>
              <DialogCloseButton />
              <DialogTitle>Create group</DialogTitle>
            </DialogHeader>
            <CreateDeviceForm setOpen={setCreateFormDialog} />
          </DialogContent>
        </DialogPortal>
      </DialogRoot>

      <DialogRoot open={updateDeviceFormDialog() != BigInt(0)} onOpenChange={() => setUpdateDeviceFormDialog(BigInt(0))}>
        <DialogPortal>
          <DialogOverlay />
          <DialogContent>
            <DialogHeader>
              <DialogCloseButton />
              <DialogTitle>Update group</DialogTitle>
            </DialogHeader>
            <UpdateDeviceForm setOpen={() => setUpdateDeviceFormDialog(BigInt(0))} id={updateDeviceFormDialog()} />
          </DialogContent>
        </DialogPortal>
      </DialogRoot>

      <AlertDialogRoot open={deleteDeviceSelection() != undefined} onOpenChange={() => setDeleteDeviceSelection(undefined)}>
        <AlertDialogContent>
          <AlertDialogHeader>
            <AlertDialogTitle>Are you sure you wish to delete {deleteDeviceSelection()?.name}?</AlertDialogTitle>
          </AlertDialogHeader>
          <AlertDialogFooter>
            <AlertDialogCancel>Cancel</AlertDialogCancel>
            <AlertDialogAction variant="destructive" disabled={deleteDeviceSubmission.pending} onClick={deleteDeviceBySelection}>
              Delete
            </AlertDialogAction>
          </AlertDialogFooter>
        </AlertDialogContent>
      </AlertDialogRoot>

      <AlertDialogRoot open={deleteDeviceRowSelector()} onOpenChange={setDeleteDeviceRowSelector}>
        <AlertDialogContent>
          <AlertDialogHeader>
            <AlertDialogTitle>Are you sure you wish to delete {rowSelection.selections().length} groups?</AlertDialogTitle>
          </AlertDialogHeader>
          <AlertDialogFooter>
            <AlertDialogCancel>Cancel</AlertDialogCancel>
            <AlertDialogAction variant="destructive" disabled={deleteDeviceSubmission.pending} onClick={deleteDeviceByRowSelector}>
              Delete
            </AlertDialogAction>
          </AlertDialogFooter>
        </AlertDialogContent>
      </AlertDialogRoot>

      <div class="text-xl">Devices</div>
      <Seperator />

      <ErrorBoundary fallback={(e: Error) => <PageError error={e} />}>
        <Suspense fallback={<Skeleton class="h-32" />}>
          <div class="flex justify-between gap-2">
            <SelectRoot
              class="w-20"
              value={data()?.pageResult?.perPage}
              onChange={setPerPage}
              options={defaultPerPageOptions}
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
                disabled={previousPageDisabled()}
                onClick={previousPage}
              >
                <RiArrowsArrowLeftSLine class="h-6 w-6" />
              </Button>
              <Button
                title="Next"
                size="icon"
                disabled={nextPageDisabled()}
                onClick={nextPage}
              >
                <RiArrowsArrowRightSLine class="h-6 w-6" />
              </Button>
            </div>
          </div>
          <TableRoot>
            <TableHeader>
              <tr class="border-b">
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
                  <TableSortButton
                    name="name"
                    onClick={toggleSort}
                    sort={data()?.sort}
                  >
                    Name
                  </TableSortButton>
                </TableHead>
                <TableHead>
                  <TableSortButton
                    name="userCount"
                    onClick={toggleSort}
                    sort={data()?.sort}
                  >
                    Users
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
                <TableHead>
                  <div class="flex items-center justify-end">
                    <DropdownMenuRoot placement="bottom-end">
                      <DropdownMenuTrigger class="hover:bg-accent hover:text-accent-foreground rounded p-1" title="Actions">
                        <RiSystemMore2Line class="h-5 w-5" />
                      </DropdownMenuTrigger>
                      <DropdownMenuPortal>
                        <DropdownMenuContent>
                          <DropdownMenuItem onSelect={() => setCreateFormDialog(true)}>
                            Create
                          </DropdownMenuItem>
                          <DropdownMenuItem
                            onSelect={() => setDeviceDisableByRowSelector(true)}
                            disabled={setDeviceDisableDisabled(true) || rowSelection.selections().length == 0}
                          >
                            Disable
                          </DropdownMenuItem>
                          <DropdownMenuItem
                            onSelect={() => setDeviceDisableByRowSelector(false)}
                            disabled={setDeviceDisableDisabled(false) || rowSelection.selections().length == 0}
                          >
                            Enable
                          </DropdownMenuItem>
                          <DropdownMenuItem
                            onSelect={() => setDeleteDeviceRowSelector(true)}
                            disabled={rowSelection.selections().length == 0}
                          >
                            Delete
                          </DropdownMenuItem>
                          <DropdownMenuArrow />
                        </DropdownMenuContent>
                      </DropdownMenuPortal>
                    </DropdownMenuRoot>
                  </div>
                </TableHead>
              </tr>
            </TableHeader>
            <TableBody>
              <For each={data()?.items}>
                {(item, index) => {
                  const onClick = () => navigate(`./${item.id}`)
                  const toggleDeviceDisable = () => setDeviceDisable({ items: [{ id: item.id, disable: !item.disabled }] })

                  return (
                    <TableRow>
                      <TableHead>
                        <CheckboxRoot checked={rowSelection.rows[index()]?.checked} onChange={(v) => rowSelection.set(item.id, v)}>
                          <CheckboxControl />
                        </CheckboxRoot>
                      </TableHead>
                      <TableCell onClick={onClick} class="max-w-48 cursor-pointer select-none" title={item.name}>
                        <div class="truncate">{item.name}</div>
                      </TableCell>
                      <TableCell onClick={onClick} class="text-nowrap cursor-pointer select-none whitespace-nowrap">{item.userCount.toString()}</TableCell>
                      <TableCell onClick={onClick} class="text-nowrap cursor-pointer select-none whitespace-nowrap">{formatDate(parseDate(item.createdAtTime))}</TableCell>
                      <TableCell class="py-0">
                        <div class="flex justify-end gap-2">
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
                          <DropdownMenuRoot placement="bottom-end">
                            <DropdownMenuTrigger class="hover:bg-accent hover:text-accent-foreground rounded p-1" title="Actions">
                              <RiSystemMore2Line class="h-5 w-5" />
                            </DropdownMenuTrigger>
                            <DropdownMenuPortal>
                              <DropdownMenuContent>
                                <DropdownMenuItem disabled={setDeviceDisableSubmission.pending} onSelect={toggleDeviceDisable}>
                                  <Show when={item.disabled} fallback={<>Disable</>}>
                                    Enable
                                  </Show>
                                </DropdownMenuItem>
                                <DropdownMenuItem onSelect={() => setUpdateDeviceFormDialog(item.id)}>
                                  Update
                                </DropdownMenuItem>
                                <DropdownMenuItem onSelect={() => setDeleteDeviceSelection(item)}>
                                  Delete
                                </DropdownMenuItem>
                                <DropdownMenuArrow />
                              </DropdownMenuContent>
                            </DropdownMenuPortal>
                          </DropdownMenuRoot>
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
    </LayoutNormal>)
}

// type CreateDeviceForm = {
//   name: string
//   description: string
// }
//
// const actionCreateDeviceForm = action((form: CreateDeviceForm) => useClient()
//   .admin.createDevice({ model: form })
//   .then(() => revalidate(getAdminDevicesPage.key))
//   .catch(throwAsFormError)
// )
//
// function CreateDeviceForm(props: { setOpen: (value: boolean) => void }) {
//   const [addMore, setAddMore] = createSignal(false)
//
//   const [createDeviceForm, { Field, Form }] = createForm<CreateDeviceForm>({ initialValues: { name: "", description: "" } });
//   const createDeviceFormAction = useAction(actionCreateDeviceForm)
//   const submit = (form: CreateDeviceForm) => createDeviceFormAction(form)
//     .then(() => batch(() => {
//       props.setOpen(addMore())
//       reset(createDeviceForm)
//     }))
//
//   return (
//     <Form class="flex flex-col gap-4" onSubmit={submit}>
//       <input class="hidden" type="text" name="username" autocomplete="username" />
//       <Field name="name" validate={required("Please enter a name.")}>
//         {(field, props) => (
//           <FieldRoot class="gap-1.5">
//             <FieldLabel field={field}>Name</FieldLabel>
//             <FieldControl field={field}>
//               <Input
//                 {...props}
//                 placeholder="Name"
//                 value={field.value}
//               />
//             </FieldControl>
//             <FieldMessage field={field} />
//           </FieldRoot>
//         )}
//       </Field>
//       <Field name="description">
//         {(field, props) => (
//           <FieldRoot class="gap-1.5">
//             <FieldLabel field={field}>Description</FieldLabel>
//             <FieldControl field={field}>
//               <Textarea
//                 {...props}
//                 value={field.value}
//                 placeholder="Description"
//               />
//             </FieldControl>
//             <FieldMessage field={field} />
//           </FieldRoot>
//         )}
//       </Field>
//       <Button type="submit" disabled={createDeviceForm.submitting}>
//         <Show when={!createDeviceForm.submitting} fallback={<>Creating group</>}>
//           Create group
//         </Show>
//       </Button>
//       <FormMessage form={createDeviceForm} />
//       <CheckboxRoot checked={addMore()} onChange={setAddMore}>
//         <CheckboxInput />
//         <CheckboxControl />
//         <CheckboxLabel>Add more</CheckboxLabel>
//       </CheckboxRoot>
//     </Form>
//   )
// }
//
// type UpdateDeviceForm = {
//   id: BigInt | any
//   name: string
//   description: string
// }
//
// const actionUpdateDeviceForm = action((form: UpdateDeviceForm) => useClient()
//   .admin.updateDevice({ id: form.id, model: form })
//   .then(() => revalidate(getAdminDevicesPage.key))
//   .catch(throwAsFormError)
// )
//
// // TODO: make the initial form data for update more readable
//
// function UpdateDeviceForm(props: { setOpen: (value: boolean) => void, id: bigint }) {
//   const [updateDeviceForm, { Field, Form }] = createForm<UpdateDeviceForm>();
//   const updateDeviceFormAction = useAction(actionUpdateDeviceForm)
//   const submit = (form: UpdateDeviceForm) => updateDeviceFormAction(form).then(() => props.setOpen(false))
//   const disabled = createAsync(async () => {
//     if (props.id == BigInt(0)) return false
//     return getDevice(props.id).then(res => syncForm(updateDeviceForm, { ...res, ...res.model }))
//   })
//
//   return (
//     <ErrorBoundary fallback={(e: Error) => <PageError error={e} />}>
//       <Suspense fallback={<Skeleton class="h-32" />}>
//         <Form class="flex flex-col gap-4" onSubmit={(form) => submit(form)}>
//           <Field name="id" type="number">
//             {(field, props) => <input {...props} type="hidden" value={field.value} />}
//           </Field>
//           <Field name="name" validate={required("Please enter a name.")}>
//             {(field, props) => (
//               <FieldRoot class="gap-1.5">
//                 <FieldLabel field={field}>Name</FieldLabel>
//                 <FieldControl field={field}>
//                   <Input
//                     {...props}
//                     disabled={disabled()}
//                     placeholder="Name"
//                     value={field.value}
//                   />
//                 </FieldControl>
//                 <FieldMessage field={field} />
//               </FieldRoot>
//             )}
//           </Field>
//           <Field name="description">
//             {(field, props) => (
//               <FieldRoot class="gap-1.5">
//                 <FieldLabel field={field}>Description</FieldLabel>
//                 <FieldControl field={field}>
//                   <Textarea
//                     {...props}
//                     disabled={disabled()}
//                     placeholder="Description"
//                   >
//                     {field.value}
//                   </Textarea>
//                 </FieldControl>
//                 <FieldMessage field={field} />
//               </FieldRoot>
//             )}
//           </Field>
//           <Button type="submit" disabled={disabled() || updateDeviceForm.submitting}>
//             <Show when={!updateDeviceForm.submitting} fallback={<>Updating group</>}>
//               Update group
//             </Show>
//           </Button>
//           <FormMessage form={updateDeviceForm} />
//         </Form>
//       </Suspense>
//     </ErrorBoundary>
//   )
// }
//
