import { action, createAsync, revalidate, useAction, useNavigate, useSearchParams, useSubmission } from "@solidjs/router";
import { AlertDialogAction, AlertDialogCancel, AlertDialogModal, AlertDialogDescription, AlertDialogFooter, AlertDialogHeader, AlertDialogRoot, AlertDialogTitle, } from "~/ui/AlertDialog";
import { DropdownMenuArrow, DropdownMenuContent, DropdownMenuItem, DropdownMenuPortal, DropdownMenuRoot } from "~/ui/DropdownMenu";
import { AdminDevicesPageSearchParams, getAdminDevicesPage } from "./Devices.data";
import { ErrorBoundary, For, Show, Suspense, createSignal } from "solid-js";
import { catchAsToast, createPagePagination, createRowSelection, createToggleSortField, createModal, formatDate, isTableRowClick, parseDate, throwAsFormError, } from "~/lib/utils";
import { parseOrder } from "~/lib/utils";
import { TableBody, TableCaption, TableCell, TableHead, TableHeader, TableRoot, TableRow, } from "~/ui/Table";
import { useClient } from "~/providers/client";
import { CheckboxControl, CheckboxLabel, CheckboxRoot } from "~/ui/Checkbox";
import { Skeleton } from "~/ui/Skeleton";
import { PageError } from "~/ui/Page";
import { TooltipArrow, TooltipContent, TooltipRoot, TooltipTrigger } from "~/ui/Tooltip";
import { LayoutNormal } from "~/ui/Layout";
import { GetDeviceResp, ListDeviceFeaturesResp_Item, SetDeviceDisableReq } from "~/twirp/rpc";
import { Crud } from "~/components/Crud";
import { RiSystemLockLine } from "solid-icons/ri";
import { DialogOverflow, DialogHeader, DialogContent, DialogOverlay, DialogPortal, DialogRoot, DialogTitle } from "~/ui/Dialog";
import { Button } from "~/ui/Button";
import { FieldLabel, FieldMessage, FieldRoot, FormMessage, SelectFieldRoot, fieldControlProps } from "~/ui/Form";
import { FieldElementProps, FieldStore, FormStore, createForm, required, reset, setValue } from "@modular-forms/solid";
import { Input } from "~/ui/Input";
import { SelectContent, SelectHTML, SelectItem, SelectLabel, SelectListbox, SelectTrigger, SelectValue } from "~/ui/Select";
import { getDevice, getListDeviceFeatures, getListLocations } from "./data";
import { Shared } from "~/components/Shared";
import { Badge } from "~/ui/Badge";
import { BreadcrumbsItem, BreadcrumbsRoot } from "~/ui/Breadcrumbs";

const actionDeleteDevice = action((ids: bigint[]) => useClient()
  .admin.deleteDevice({ ids })
  .then(() => revalidate(getAdminDevicesPage.key))
  .catch(catchAsToast))

const actionSetDisable = action((input: SetDeviceDisableReq) => useClient()
  .admin.setDeviceDisable(input)
  .then(() => revalidate(getAdminDevicesPage.key))
  .catch(catchAsToast))

export function AdminDevices() {
  const navigate = useNavigate()
  const [searchParams] = useSearchParams<AdminDevicesPageSearchParams>()

  const data = createAsync(() => getAdminDevicesPage({
    page: {
      page: Number(searchParams.page) || 0,
      perPage: Number(searchParams.perPage) || 0
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
  const [createFormModal, setCreateFormModal] = createSignal(false);

  // Update
  const updateFormModal = createModal(BigInt(0))

  // Delete
  const deleteSubmission = useSubmission(actionDeleteDevice)
  const deleteAction = useAction(actionDeleteDevice)
  // Single
  const deleteModal = createModal({ name: "", id: BigInt(0) })
  const deleteSubmit = () =>
    deleteAction([deleteModal.value().id])
      .then(deleteModal.close)
  // Multiple
  const [deleteMultipleModal, setDeleteMultipleModal] = createSignal(false)
  const deleteMultipleSubmit = () =>
    deleteAction(rowSelection.selections())
      .then(() => setDeleteMultipleModal(false))

  // Disable/Enable
  const setDisableSubmission = useSubmission(actionSetDisable)
  const setDisableAction = useAction(actionSetDisable)
  const setDisableSubmit = (disable: boolean) =>
    setDisableAction({ items: rowSelection.selections().map(v => ({ id: v, disable })) })
      .then(() => rowSelection.setAll(false))

  return (
    <LayoutNormal class="max-w-4xl">
      <DialogRoot open={createFormModal()} onOpenChange={setCreateFormModal}>
        <DialogPortal>
          <DialogOverlay />
          <DialogContent>
            <DialogHeader>
              <DialogTitle>Create device</DialogTitle>
            </DialogHeader>
            <DialogOverflow>
              <CreateForm onSubmit={() => setCreateFormModal(false)} />
            </DialogOverflow>
          </DialogContent>
        </DialogPortal>
      </DialogRoot>

      <DialogRoot open={updateFormModal.open()} onOpenChange={updateFormModal.close}>
        <DialogPortal>
          <DialogOverlay />
          <DialogContent>
            <DialogHeader>
              <DialogTitle>Update device</DialogTitle>
            </DialogHeader>
            <DialogOverflow>
              <UpdateForm onSubmit={updateFormModal.close} id={updateFormModal.value()} />
            </DialogOverflow>
          </DialogContent>
        </DialogPortal>
      </DialogRoot>

      <AlertDialogRoot open={deleteModal.open()} onOpenChange={deleteModal.close}>
        <AlertDialogModal>
          <AlertDialogHeader>
            <AlertDialogTitle>Are you sure you wish to delete {deleteModal.value().name}?</AlertDialogTitle>
          </AlertDialogHeader>
          <AlertDialogFooter>
            <AlertDialogCancel>Cancel</AlertDialogCancel>
            <AlertDialogAction variant="destructive" disabled={deleteSubmission.pending} onClick={deleteSubmit}>
              Delete
            </AlertDialogAction>
          </AlertDialogFooter>
        </AlertDialogModal>
      </AlertDialogRoot>

      <AlertDialogRoot open={deleteMultipleModal()} onOpenChange={setDeleteMultipleModal}>
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
            <AlertDialogAction variant="destructive" disabled={deleteSubmission.pending} onClick={deleteMultipleSubmit}>
              Delete
            </AlertDialogAction>
          </AlertDialogFooter>
        </AlertDialogModal>
      </AlertDialogRoot>

      <Shared.Title>
        <BreadcrumbsRoot>
          <BreadcrumbsItem>
            Devices
          </BreadcrumbsItem>
        </BreadcrumbsRoot>
      </Shared.Title>

      <ErrorBoundary fallback={(e) => <PageError error={e} />}>
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
                        <DropdownMenuItem onSelect={() => setCreateFormModal(true)}>
                          Create
                        </DropdownMenuItem>
                        <DropdownMenuItem
                          disabled={rowSelection.selections().length == 0 || setDisableSubmission.pending}
                          onSelect={() => setDisableSubmit(true)}
                        >
                          Disable
                        </DropdownMenuItem>
                        <DropdownMenuItem
                          disabled={rowSelection.selections().length == 0 || setDisableSubmission.pending}
                          onSelect={() => setDisableSubmit(false)}
                        >
                          Enable
                        </DropdownMenuItem>
                        <DropdownMenuItem
                          disabled={rowSelection.selections().length == 0 || deleteSubmission.pending}
                          onSelect={() => setDeleteMultipleModal(true)}
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
                  const toggleDisableSubmit = () => setDisableAction({ items: [{ id: item.id, disable: !item.disabled }] })

                  return (
                    <TableRow
                      data-state={rowSelection.rows[index()]?.checked ? "selected" : ""}
                      onClick={(t) => isTableRowClick(t) && navigate(`./${item.id}`)}
                      class="cursor-pointer"
                    >
                      <TableHead>
                        <CheckboxRoot
                          checked={rowSelection.rows[index()]?.checked}
                          onChange={(v) => rowSelection.set(item.id, v)}
                        >
                          <CheckboxControl />
                        </CheckboxRoot>
                      </TableHead>
                      <TableCell>{item.name}</TableCell>
                      <TableCell>{item.url}</TableCell>
                      <TableCell>{formatDate(parseDate(item.createdAtTime))}</TableCell>
                      <Crud.LastTableCell>
                        <Show when={item.disabled}>
                          <TooltipRoot>
                            <TooltipTrigger>
                              <RiSystemLockLine class="size-5" />
                            </TooltipTrigger>
                            <TooltipContent>
                              <TooltipArrow />
                              Disabled since {formatDate(parseDate(item.disabledAtTime))}
                            </TooltipContent>
                          </TooltipRoot>
                        </Show>
                        <DropdownMenuRoot placement="bottom-end">
                          <Crud.MoreDropdownMenuTrigger />
                          <DropdownMenuPortal>
                            <DropdownMenuContent>
                              <DropdownMenuItem onSelect={() => updateFormModal.setValue(item.id)}>
                                Edit
                              </DropdownMenuItem>
                              <DropdownMenuItem
                                disabled={setDisableSubmission.pending}
                                onSelect={toggleDisableSubmit}
                              >
                                <Show when={item.disabled} fallback="Disable">Enable</Show>
                              </DropdownMenuItem>
                              <DropdownMenuItem
                                disabled={deleteSubmission.pending}
                                onSelect={() => deleteModal.setValue(item)}
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
              <Crud.PageMetadata pageResult={data()?.pageResult} />
            </TableCaption>
          </TableRoot>
        </Suspense>
      </ErrorBoundary>
    </LayoutNormal>
  )
}

type CreateForm = {
  name: string
  url: string
  username: string
  password: string
  location: string
  features: {
    array: string[]
  }
}

function CreateForm(props: { onSubmit?: () => void }) {
  const locations = createAsync(() => getListLocations())

  const [addMore, setAddMore] = createSignal(false)

  const [form, { Field, Form }] = createForm<CreateForm>({});
  const submit = async (data: CreateForm) => {
    await useClient()
      .admin.createDevice({ ...data, features: data.features.array })
      .then(() => revalidate(getAdminDevicesPage.key))
      .catch(throwAsFormError)
      .then(() => {
        if (addMore()) {
          reset(form, {
            initialValues: {
              ...data,
              name: "",
              url: ""
            },
          })
        } else {
          props.onSubmit && props.onSubmit()
        }
      })
  }

  return (
    <ErrorBoundary fallback={(e) => <PageError error={e} />}>
      <Suspense fallback={<Skeleton class="h-32" />}>
        <Form class="flex flex-col gap-4" onSubmit={submit}>
          <Field name="name">
            {(field, props) => (
              <FieldRoot>
                <FieldLabel field={field}>Name</FieldLabel>
                <Input
                  {...props}
                  {...fieldControlProps(field)}
                  placeholder="Name"
                  value={field.value}
                />
                <FieldMessage field={field} />
              </FieldRoot>
            )}
          </Field>
          <Field name="url" validate={required("Please enter a URL.")}>
            {(field, props) => (
              <FieldRoot>
                <FieldLabel field={field}>URL</FieldLabel>
                <Input
                  {...props}
                  {...fieldControlProps(field)}
                  placeholder="URL"
                  value={field.value}
                />
                <FieldMessage field={field} />
              </FieldRoot>
            )}
          </Field>
          <Field name="username">
            {(field, props) => (
              <FieldRoot>
                <FieldLabel field={field}>Username</FieldLabel>
                <Input
                  {...props}
                  {...fieldControlProps(field)}
                  placeholder="Username"
                  value={field.value}
                />
                <FieldMessage field={field} />
              </FieldRoot>
            )}
          </Field>
          <Field name="password">
            {(field, props) => (
              <FieldRoot>
                <FieldLabel field={field}>Password</FieldLabel>
                <Input
                  {...props}
                  {...fieldControlProps(field)}
                  autocomplete="off"
                  placeholder="Password"
                  type="password"
                  value={field.value}
                />
                <FieldMessage field={field} />
              </FieldRoot>
            )}
          </Field>
          <Field name="location">
            {(field, props) => (
              <FieldRoot>
                <FieldLabel field={field}>Location</FieldLabel>
                <SelectHTML
                  {...props}
                  {...fieldControlProps(field)}
                  value={field.value}
                >
                  <For each={locations()}>
                    {v => <option value={v}>{v}</option>}
                  </For>
                </SelectHTML>
                <FieldMessage field={field} />
              </FieldRoot>
            )}
          </Field>
          <Field name="features.array" type="string[]">
            {(field, props) => <DeviceFeaturesField form={form} field={field} props={props} />}
          </Field>
          <Button type="submit" disabled={form.submitting}>
            <Show when={!form.submitting} fallback="Creating device">Create device</Show>
          </Button>
          <FormMessage form={form} />
          <CheckboxRoot checked={addMore()} onChange={setAddMore} class="flex items-center gap-2">
            <CheckboxControl />
            <CheckboxLabel>Add more</CheckboxLabel>
          </CheckboxRoot>
        </Form>
      </Suspense>
    </ErrorBoundary>
  )
}

function UpdateForm(props: { onSubmit: () => void | Promise<void>, id: bigint }) {
  const device = createAsync(() => getDevice(props.id))
  const refetchDevice = () => revalidate(getDevice.keyFor(props.id))
  const onSubmit = () =>
    revalidate([getAdminDevicesPage.key, getDevice.keyFor(props.id)])
      .then(props.onSubmit)

  return (
    <ErrorBoundary fallback={(e) => <PageError error={e} />}>
      <Suspense fallback={<Skeleton class="h-32" />}>
        <Show when={device()}>
          <UpdateFormForm onSubmit={onSubmit} device={device()!} refetchDevice={refetchDevice} />
        </Show>
      </Suspense>
    </ErrorBoundary>
  )
}

type UpdateForm = {
  id: any
  name: string
  url: string
  username: string
  newPassword: string
  location: string
  features: {
    array: string[]
  }
}

function UpdateFormForm(props: { onSubmit: () => void | Promise<void>, device: GetDeviceResp, refetchDevice: () => Promise<void> }) {
  const locations = createAsync(() => getListLocations())

  const formInitialValues = (): UpdateForm => ({
    ...props.device,
    features: { array: props.device.features || [] },
    newPassword: ""
  })
  const [form, { Field, Form }] = createForm<UpdateForm>({
    initialValues: formInitialValues()
  });
  const formReset = () => props.refetchDevice().then(() => reset(form, { initialValues: formInitialValues() }))
  const submit = (data: UpdateForm) => useClient()
    .admin.updateDevice({ ...data, features: data.features.array })
    .then(props.onSubmit)
    .catch(throwAsFormError)

  return (
    <Form class="flex flex-col gap-4" onSubmit={submit}>
      <Field name="id" type="number">
        {(field, props) => <input {...props} type="hidden" value={field.value} />}
      </Field>
      <Field name="name">
        {(field, props) => (
          <FieldRoot>
            <FieldLabel field={field}>Name</FieldLabel>
            <Input
              {...props}
              {...fieldControlProps(field)}
              placeholder="Name"
              value={field.value}
            />
            <FieldMessage field={field} />
          </FieldRoot>
        )}
      </Field>
      <Field name="url" validate={required("Please enter a URL.")}>
        {(field, props) => (
          <FieldRoot>
            <FieldLabel field={field}>URL</FieldLabel>
            <Input
              {...props}
              {...fieldControlProps(field)}
              placeholder="URL"
              value={field.value}
            />
            <FieldMessage field={field} />
          </FieldRoot>
        )}
      </Field>
      <Field name="username">
        {(field, props) => (
          <FieldRoot>
            <FieldLabel field={field}>Username</FieldLabel>
            <Input
              {...props}
              {...fieldControlProps(field)}
              placeholder="Username"
              value={field.value}
            />
            <FieldMessage field={field} />
          </FieldRoot>
        )}
      </Field>
      <Field name="newPassword">
        {(field, props) => (
          <FieldRoot>
            <FieldLabel field={field}>New password</FieldLabel>
            <Input
              {...props}
              {...fieldControlProps(field)}
              autocomplete="off"
              placeholder="New password"
              type="password"
              value={field.value}
            />
            <FieldMessage field={field} />
          </FieldRoot>
        )}
      </Field>
      <Field name="location">
        {(field, props) => (
          <FieldRoot>
            <FieldLabel field={field}>Location</FieldLabel>
            <SelectHTML
              {...props}
              {...fieldControlProps(field)}
              value={field.value}
            >
              <For each={locations()}>
                {v => <option value={v}>{v}</option>}
              </For>
            </SelectHTML>
            <FieldMessage field={field} />
          </FieldRoot>
        )}
      </Field>
      <Field name="features.array" type="string[]">
        {(field, props) => <DeviceFeaturesField form={form} field={field} props={props} />}
      </Field>
      <div class="flex flex-col gap-4 sm:flex-row-reverse">
        <Button type="submit" disabled={form.submitting} class="sm:flex-1">
          <Show when={!form.submitting} fallback="Updating device">Update device</Show>
        </Button>
        <Button type="button" onClick={formReset} variant="destructive" disabled={form.submitting}>Reset</Button>
      </div>
      <FormMessage form={form} />
    </Form >
  )
}

function DeviceFeaturesField(props: { form: FormStore<any, undefined>, field: FieldStore<any, any>, props: FieldElementProps<any, any> }) {
  const deviceFeatures = createAsync(() => getListDeviceFeatures())

  return (
    <Show when={deviceFeatures()}>
      <SelectFieldRoot<ListDeviceFeaturesResp_Item>
        field={props.field}
        selectProps={props.props}
        options={deviceFeatures()!}
        optionValue="value"
        optionTextValue="name"
        placeholder="Features"
        itemComponent={props => (
          <SelectItem item={props.item} >
            <div class="flex gap-2">
              {props.item.rawValue.name} <p class="text-muted-foreground" >{props.item.rawValue.description}</p>
            </div>
          </SelectItem>
        )}
        value={deviceFeatures()?.filter(v => props.field.value?.includes(v.value))}
        onChange={(v) => setValue(props.form, props.field.name, v.map(v => v.value))}
        multiple
        class="space-y-2"
      >
        <SelectLabel>Features</SelectLabel>
        <SelectTrigger aria-label="Feature">
          <SelectValue<ListDeviceFeaturesResp_Item>>
            {state =>
              <div class="flex gap-2">
                <For each={state.selectedOptions()}>
                  {v => <Badge>{v.name}</Badge>}
                </For>
              </div>
            }
          </SelectValue>
        </SelectTrigger>
        <SelectContent>
          <SelectListbox />
        </SelectContent>
      </SelectFieldRoot>
    </Show >
  )
}
