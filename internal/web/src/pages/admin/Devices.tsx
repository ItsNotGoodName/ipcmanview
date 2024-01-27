import { action, createAsync, revalidate, useAction, useNavigate, useSearchParams, useSubmission } from "@solidjs/router";
import { AlertDialogAction, AlertDialogCancel, AlertDialogContent, AlertDialogDescription, AlertDialogFooter, AlertDialogHeader, AlertDialogRoot, AlertDialogTitle, } from "~/ui/AlertDialog";
import { DropdownMenuArrow, DropdownMenuContent, DropdownMenuItem, DropdownMenuPortal, DropdownMenuRoot, DropdownMenuTrigger } from "~/ui/DropdownMenu";
import { AdminDevicesPageSearchParams, getAdminDevicesPage } from "./Devices.data";
import { ErrorBoundary, For, Show, Suspense, batch, createEffect, createSignal } from "solid-js";
import { RiArrowsArrowLeftSLine, RiArrowsArrowRightSLine, RiSystemLockLine, RiSystemMore2Line, } from "solid-icons/ri";
import { Button } from "~/ui/Button";
import { catchAsToast, createPagePagination, createRowSelection, createToggleSortField, formatDate, parseDate, throwAsFormError } from "~/lib/utils";
import { parseOrder } from "~/lib/utils";
import { TableBody, TableCaption, TableCell, TableHead, TableHeader, TableRoot, TableRow, } from "~/ui/Table";
import { Seperator } from "~/ui/Seperator";
import { useClient } from "~/providers/client";
import { DialogCloseButton, DialogContent, DialogHeader, DialogOverlay, DialogPortal, DialogRoot, DialogTitle, } from "~/ui/Dialog";
import { CheckboxControl, CheckboxInput, CheckboxLabel, CheckboxRoot } from "~/ui/Checkbox";
import { Skeleton } from "~/ui/Skeleton";
import { PageError } from "~/ui/Page";
import { TooltipContent, TooltipRoot, TooltipTrigger } from "~/ui/Tooltip";
import { LayoutNormal } from "~/ui/Layout";
import { SetDeviceDisableReq } from "~/twirp/rpc";
import { Crud } from "~/components/Crud";
import { FieldControl, FieldLabel, FieldMessage, FieldRoot, FormMessage } from "~/ui/Form";
import { Textarea } from "~/ui/Textarea";
import { Input } from "~/ui/Input";
import { createForm, required, reset } from "@modular-forms/solid";
import { SheetCloseButton, SheetContent, SheetHeader, SheetOverlay, SheetPortal, SheetRoot, SheetTitle } from "~/ui/Sheet";
import { getListDeviceFeatures, getListLocations } from "./data";
import { SelectContent, SelectDescription, SelectHTML, SelectHiddenSelect, SelectItem, SelectLabel, SelectListbox, SelectRoot, SelectTrigger, SelectValue } from "~/ui/Select";
import { Select } from "@kobalte/core";
import { labelVariants } from "~/ui/Label";

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

  // Create
  const [createFormSheet, setCreateFormSheet] = createSignal(true);

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
  const [deleteDeviceRowSelection, setDeleteDeviceRowSelection] = createSignal(false)
  const deleteDeviceByRowSelection = () => deleteDeviceAction(rowSelection.selections())
    .then(() => setDeleteDeviceRowSelection(false))

  // Disable/Enable
  const setDeviceDisableSubmission = useSubmission(actionSetDeviceDisable)
  const setDeviceDisable = useAction(actionSetDeviceDisable)
  const setDeviceDisableByRowSelection = (disable: boolean) => setDeviceDisable({ items: rowSelection.selections().map(v => ({ id: v, disable })) })
    .then(() => rowSelection.setAll(false))

  return (
    <LayoutNormal>
      <SheetRoot open={createFormSheet()} onOpenChange={setCreateFormSheet}>
        <SheetPortal>
          <SheetOverlay />
          <SheetContent>
            <SheetHeader>
              <SheetCloseButton />
              <SheetTitle>Create group</SheetTitle>
            </SheetHeader>
            <CreateDeviceForm setOpen={setCreateFormSheet} />
          </SheetContent>
        </SheetPortal>
      </SheetRoot>

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

      <AlertDialogRoot open={deleteDeviceRowSelection()} onOpenChange={setDeleteDeviceRowSelection}>
        <AlertDialogContent>
          <AlertDialogHeader>
            <AlertDialogTitle>Are you sure you wish to delete {rowSelection.selections().length} groups?</AlertDialogTitle>
            <AlertDialogDescription class="max-h-32 overflow-y-auto">
              <For each={data()?.items}>
                {(e, index) =>
                  <Show when={rowSelection.rows[index()].checked}>
                    <div>
                      {e.name}
                    </div>
                  </Show>
                }
              </For>
            </AlertDialogDescription>
          </AlertDialogHeader>
          <AlertDialogFooter>
            <AlertDialogCancel>Cancel</AlertDialogCancel>
            <AlertDialogAction variant="destructive" disabled={deleteDeviceSubmission.pending} onClick={deleteDeviceByRowSelection}>
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
            <Crud.PerPageSelect
              class="w-20"
              perPage={data()?.pageResult?.perPage}
              onChange={pagination.setPerPage}
            />
            <div class="flex gap-2">
              <Button
                title="Previous"
                size="icon"
                disabled={pagination.previousPageDisabled()}
                onClick={pagination.previousPage}
              >
                <RiArrowsArrowLeftSLine class="h-6 w-6" />
              </Button>
              <Button
                title="Next"
                size="icon"
                disabled={pagination.nextPageDisabled()}
                onClick={pagination.nextPage}
              >
                <RiArrowsArrowRightSLine class="h-6 w-6" />
              </Button>
            </div>
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
                <TableHead>
                  <div class="flex items-center justify-end">
                    <DropdownMenuRoot placement="bottom-end">
                      <DropdownMenuTrigger class="hover:bg-accent hover:text-accent-foreground rounded p-1" title="Actions">
                        <RiSystemMore2Line class="h-5 w-5" />
                      </DropdownMenuTrigger>
                      <DropdownMenuPortal>
                        <DropdownMenuContent>
                          <DropdownMenuItem onSelect={() => setCreateFormSheet(true)}>
                            Create
                          </DropdownMenuItem>
                          <DropdownMenuItem
                            disabled={rowSelection.selections().length == 0 || setDeviceDisableSubmission.pending}
                            onSelect={() => setDeviceDisableByRowSelection(true)}
                          >
                            Disable
                          </DropdownMenuItem>
                          <DropdownMenuItem
                            disabled={rowSelection.selections().length == 0 || setDeviceDisableSubmission.pending}
                            onSelect={() => setDeviceDisableByRowSelection(false)}
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
                  </div>
                </TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              <For each={data()?.items}>
                {(item, index) => {
                  const onClick = () => navigate(`./${item.id}`)
                  const toggleDeviceDisable = () => setDeviceDisable({ items: [{ id: item.id, disable: !item.disabled }] })

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
                          <DropdownMenuTrigger class="hover:bg-accent hover:text-accent-foreground rounded p-1" title="Actions">
                            <RiSystemMore2Line class="h-5 w-5" />
                          </DropdownMenuTrigger>
                          <DropdownMenuPortal>
                            <DropdownMenuContent>
                              <DropdownMenuItem onSelect={() => setUpdateDeviceFormDialog(item.id)}>
                                Edit
                              </DropdownMenuItem>
                              <DropdownMenuItem
                                disabled={setDeviceDisableSubmission.pending}
                                onSelect={toggleDeviceDisable}
                              >
                                <Show when={item.disabled} fallback={<>Disable</>}>
                                  Enable
                                </Show>
                              </DropdownMenuItem>
                              <DropdownMenuItem
                                disabled={deleteDeviceSubmission.pending}
                                onSelect={() => setDeleteDeviceSelection(item)}
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

type CreateDeviceForm = {
  name: string
  url: string
  username: string
  password: string
  location: string
  features: string[]
}

const actionCreateDeviceForm = action((form: CreateDeviceForm) => useClient()
  .admin.createDevice({ model: form })
  .then(() => revalidate(getAdminDevicesPage.key))
  .catch(throwAsFormError)
)

function CreateDeviceForm(props: { setOpen: (value: boolean) => void }) {
  const [addMore, setAddMore] = createSignal(false)

  const locations = createAsync(getListLocations)
  const deviceFeatures = createAsync(getListDeviceFeatures)
  const [createDeviceForm, { Field, Form }] = createForm<CreateDeviceForm>({ initialValues: { name: "", description: "" } });
  const createDeviceFormAction = useAction(actionCreateDeviceForm)
  const submit = (form: CreateDeviceForm) => createDeviceFormAction(form)
    .then(() => batch(() => {
      props.setOpen(addMore())
      reset(createDeviceForm)
    }))

  createEffect(() => {
    console.log(locations()?.locations)
  })

  return (
    <Suspense fallback={<Skeleton class="h-32" />}>
      <Form class="flex flex-col gap-4" onSubmit={submit}>
        <input class="hidden" type="text" name="username" autocomplete="username" />
        <Field name="name">
          {(field, props) => (
            <FieldRoot class="gap-1.5">
              <FieldLabel field={field}>Name</FieldLabel>
              <FieldControl field={field}>
                <Input
                  {...props}
                  placeholder="Name"
                  value={field.value}
                />
              </FieldControl>
              <FieldMessage field={field} />
            </FieldRoot>
          )}
        </Field>
        <Field name="url">
          {(field, props) => (
            <FieldRoot class="gap-1.5">
              <FieldLabel field={field}>URL</FieldLabel>
              <FieldControl field={field}>
                <Input
                  {...props}
                  placeholder="URL"
                  value={field.value}
                />
              </FieldControl>
              <FieldMessage field={field} />
            </FieldRoot>
          )}
        </Field>
        <Field name="username">
          {(field, props) => (
            <FieldRoot class="gap-1.5">
              <FieldLabel field={field}>Username</FieldLabel>
              <FieldControl field={field}>
                <Input
                  {...props}
                  placeholder="Username"
                  value={field.value}
                />
              </FieldControl>
              <FieldMessage field={field} />
            </FieldRoot>
          )}
        </Field>
        <Field name="password">
          {(field, props) => (
            <FieldRoot class="gap-1.5">
              <FieldLabel field={field}>Password</FieldLabel>
              <FieldControl field={field}>
                <Input
                  {...props}
                  placeholder="Password"
                  type="password"
                  value={field.value}
                />
              </FieldControl>
              <FieldMessage field={field} />
            </FieldRoot>
          )}
        </Field>
        <Field name="location">
          {(field, props) => (
            <FieldRoot class="gap-1.5">
              <FieldLabel field={field}>Location</FieldLabel>
              <FieldControl field={field}>
                <SelectHTML
                  {...props}
                  value={field.value}
                >
                  <For each={locations()}>
                    {v =>
                      <option value={v}>{v}</option>
                    }
                  </For>
                </SelectHTML>
              </FieldControl>
              <FieldMessage field={field} />
            </FieldRoot>
          )}
        </Field>
        <Field name="features" type="string[]">
          {(field, props) => (
            <FieldRoot class="gap-1.5">
              <FieldLabel field={field}>Features</FieldLabel>
              <FieldControl field={field}>
                <SelectHTML
                  {...props}
                  class="h-32"
                  value={field.value}
                  multiple
                >
                  <For each={deviceFeatures()}>
                    {v =>
                      <option value={v.value}>{v.name}</option>
                    }
                  </For>
                </SelectHTML>
              </FieldControl>
              <FieldMessage field={field} />
            </FieldRoot>
          )}
        </Field>
        <Button type="submit" disabled={createDeviceForm.submitting}>
          <Show when={!createDeviceForm.submitting} fallback={<>Creating device</>}>
            Create device
          </Show>
        </Button>
        <FormMessage form={createDeviceForm} />
        <CheckboxRoot checked={addMore()} onChange={setAddMore}>
          <CheckboxInput />
          <CheckboxControl />
          <CheckboxLabel>Add more</CheckboxLabel>
        </CheckboxRoot>
      </Form>
    </Suspense>
  )
}

type UpdateDeviceForm = {
  id: any
  name: string
  description: string
}

const actionUpdateDeviceForm = action((form: UpdateDeviceForm) => useClient()
  .admin.updateDevice({ id: form.id, model: form })
  .then(() => revalidate(getAdminDevicesPage.key))
  .catch(throwAsFormError)
)

function UpdateDeviceForm(props: { setOpen: (value: boolean) => void, id: bigint }) {
  const [updateDeviceForm, { Field, Form }] = createForm<UpdateDeviceForm>();
  const updateDeviceFormAction = useAction(actionUpdateDeviceForm)
  const submit = (form: UpdateDeviceForm) => updateDeviceFormAction(form)
    .then(() => props.setOpen(false))
  const [form] = createResource(() => getDevice(props.id)
    .then((data) => setupForm(updateDeviceForm, { ...data, ...data.model })))

  return (
    <Show when={!form.error} fallback={<PageError error={form.error} />}>
      <Form class="flex flex-col gap-4" onSubmit={(form) => submit(form)}>
        <Field name="id" type="number">
          {(field, props) => <input {...props} type="hidden" value={field.value} />}
        </Field>
        <Field name="name" validate={required("Please enter a name.")}>
          {(field, props) => (
            <FieldRoot class="gap-1.5">
              <FieldLabel field={field}>Name</FieldLabel>
              <FieldControl field={field}>
                <Input
                  {...props}
                  placeholder="Name"
                  value={field.value}
                  disabled={form.loading}
                />
              </FieldControl>
              <FieldMessage field={field} />
            </FieldRoot>
          )}
        </Field>
        <Field name="description">
          {(field, props) => (
            <FieldRoot class="gap-1.5">
              <FieldLabel field={field}>Description</FieldLabel>
              <FieldControl field={field}>
                <Textarea
                  {...props}
                  placeholder="Description"
                  disabled={form.loading}
                >
                  {field.value}
                </Textarea>
              </FieldControl>
              <FieldMessage field={field} />
            </FieldRoot>
          )}
        </Field>
        <Button type="submit" disabled={form.loading || updateDeviceForm.submitting}>
          <Show when={!updateDeviceForm.submitting} fallback={<>Updating group</>}>
            Update group
          </Show>
        </Button>
        <FormMessage form={updateDeviceForm} />
      </Form>
    </Show>
  )
}
