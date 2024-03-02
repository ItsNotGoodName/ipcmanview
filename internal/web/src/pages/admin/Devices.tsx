import { action, createAsync, revalidate, useAction, useNavigate, useSearchParams, useSubmission } from "@solidjs/router";
import { AlertDialogAction, AlertDialogCancel, AlertDialogModal, AlertDialogDescription, AlertDialogFooter, AlertDialogHeader, AlertDialogRoot, AlertDialogTitle, } from "~/ui/AlertDialog";
import { DropdownMenuArrow, DropdownMenuContent, DropdownMenuItem, DropdownMenuPortal, DropdownMenuRoot } from "~/ui/DropdownMenu";
import { AdminDevicesPageSearchParams, getAdminDevicesPage } from "./Devices.data";
import { ErrorBoundary, For, Show, Suspense, createSignal } from "solid-js";
import { catchAsToast, createPagePagination, createRowSelection, createToggleSortField, createModal, formatDate, isTableDataClick, parseDate, throwAsFormError, validationState, setFormValue, createLoading, } from "~/lib/utils";
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
import { FormMessage } from "~/ui/Form";
import { FieldStore, FormStore, createForm, required, reset, setValue } from "@modular-forms/solid";
import { SelectContent, SelectErrorMessage, SelectItem, SelectLabel, SelectListBoxVirtual, SelectListbox, SelectPortal, SelectRoot, SelectTrigger, SelectValue } from "~/ui/Select";
import { getDevice, getListDeviceFeatures, getListLocations } from "./data";
import { Shared } from "~/components/Shared";
import { Badge } from "~/ui/Badge";
import { BreadcrumbsItem, BreadcrumbsRoot } from "~/ui/Breadcrumbs";
import { TextFieldErrorMessage, TextFieldInput, TextFieldLabel, TextFieldRoot } from "~/ui/TextField";

const actionDeleteDevice = action((ids: string[]) => useClient()
  .admin.deleteDevice({ ids })
  .then()
  .catch(catchAsToast))

const actionSetDisable = action((input: SetDeviceDisableReq) => useClient()
  .admin.setDeviceDisable(input)
  .then()
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
  const rowSelection = createRowSelection(() => data()?.items.map((v) => ({ id: v.id })) || [])

  // List
  const pagination = createPagePagination(() => data()?.pageResult)
  const toggleSort = createToggleSortField(() => data()?.sort)

  // Create
  const [createFormModal, setCreateFormModal] = createSignal(false);

  // Update
  const updateFormModal = createModal("")

  // Delete
  const deleteSubmission = useSubmission(actionDeleteDevice)
  const deleteAction = useAction(actionDeleteDevice)
  // Single
  const deleteModal = createModal({ name: "", id: "" })
  const submitDelete = () =>
    deleteAction([deleteModal.value().id])
      .then(deleteModal.setClose)
  // Multiple
  const [deleteMultipleModal, setDeleteMultipleModal] = createSignal(false)
  const submitDeleteMultiple = () =>
    deleteAction(rowSelection.selections())
      .then(() => setDeleteMultipleModal(false))

  // Disable/Enable
  const setDisableSubmission = useSubmission(actionSetDisable)
  const setDisableAction = useAction(actionSetDisable)
  const submitSetDisable = (disable: boolean) =>
    setDisableAction({ items: rowSelection.selections().map(id => ({ id, disable })) })
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
              <CreateForm onClose={() => setCreateFormModal(false)} />
            </DialogOverflow>
          </DialogContent>
        </DialogPortal>
      </DialogRoot>

      <DialogRoot open={updateFormModal.open()} onOpenChange={updateFormModal.setClose}>
        <DialogPortal>
          <DialogOverlay />
          <DialogContent>
            <DialogHeader>
              <DialogTitle>Update device</DialogTitle>
            </DialogHeader>
            <DialogOverflow>
              <UpdateForm onClose={updateFormModal.setClose} id={updateFormModal.value()} />
            </DialogOverflow>
          </DialogContent>
        </DialogPortal>
      </DialogRoot>

      <AlertDialogRoot open={deleteModal.open()} onOpenChange={deleteModal.setClose}>
        <AlertDialogModal>
          <AlertDialogHeader>
            <AlertDialogTitle>Are you sure you wish to delete {deleteModal.value().name}?</AlertDialogTitle>
          </AlertDialogHeader>
          <AlertDialogFooter>
            <AlertDialogCancel>Cancel</AlertDialogCancel>
            <AlertDialogAction variant="destructive" disabled={deleteSubmission.pending} onClick={submitDelete}>
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
                    <Show when={rowSelection.items[index()]?.checked}>
                      <li>{e.name}</li>
                    </Show>
                  }
                </For>
              </ul>
            </AlertDialogDescription>
          </AlertDialogHeader>
          <AlertDialogFooter>
            <AlertDialogCancel>Cancel</AlertDialogCancel>
            <AlertDialogAction variant="destructive" disabled={deleteSubmission.pending} onClick={submitDeleteMultiple}>
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

      <div class="flex flex-col gap-2">
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
                      checked={rowSelection.all()}
                      indeterminate={rowSelection.multiple()}
                      onChange={rowSelection.setAll}
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
                            disabled={!rowSelection.multiple() || setDisableSubmission.pending}
                            onSelect={() => submitSetDisable(true)}
                          >
                            Disable
                          </DropdownMenuItem>
                          <DropdownMenuItem
                            disabled={!rowSelection.multiple() || setDisableSubmission.pending}
                            onSelect={() => submitSetDisable(false)}
                          >
                            Enable
                          </DropdownMenuItem>
                          <DropdownMenuItem
                            disabled={!rowSelection.multiple() || deleteSubmission.pending}
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
                    const submitToggleDisable = () => setDisableAction({ items: [{ id: item.id, disable: !item.disabled }] })

                    return (
                      <TableRow
                        data-state={rowSelection.items[index()]?.checked ? "selected" : ""}
                        onClick={(t) => isTableDataClick(t) && navigate(`./${item.id}`)}
                        class="[&>td]:cursor-pointer"
                      >
                        <TableHead>
                          <CheckboxRoot
                            checked={rowSelection.items[index()]?.checked}
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
                                  onSelect={submitToggleDisable}
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
      </div>
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

function CreateForm(props: { onClose: () => void }) {
  const [addMore, setAddMore] = createSignal(false)

  const [form, { Field, Form }] = createForm<CreateForm>({
    initialValues: {
      name: "",
      url: "",
      username: "",
      password: "",
      location: "",
      features: { array: [] },
    }
  });
  const submitForm = async (data: CreateForm) => {
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
          props.onClose()
        }
      })
  }

  return (
    <ErrorBoundary fallback={(e) => <PageError error={e} />}>
      <Suspense fallback={<Skeleton class="h-32" />}>
        <Form class="flex flex-col gap-4" onSubmit={submitForm}>
          <Field name="name">
            {(field, props) => (
              <TextFieldRoot
                validationState={validationState(field.error)}
                value={field.value}
                class="space-y-2"
              >
                <TextFieldLabel>Name</TextFieldLabel>
                <TextFieldInput
                  {...props}
                  placeholder="Name"
                />
                <TextFieldErrorMessage>{field.error}</TextFieldErrorMessage>
              </TextFieldRoot>
            )}
          </Field>
          <Field name="url" validate={required("Please enter a URL.")}>
            {(field, props) => (
              <TextFieldRoot
                validationState={validationState(field.error)}
                value={field.value}
                class="space-y-2"
              >
                <TextFieldLabel>URL</TextFieldLabel>
                <TextFieldInput
                  {...props}
                  placeholder="URL"
                />
                <TextFieldErrorMessage>{field.error}</TextFieldErrorMessage>
              </TextFieldRoot>
            )}
          </Field>
          <Field name="username">
            {(field, props) => (
              <TextFieldRoot
                validationState={validationState(field.error)}
                value={field.value}
                class="space-y-2"
              >
                <TextFieldLabel>Username</TextFieldLabel>
                <TextFieldInput
                  {...props}
                  placeholder="Username"
                />
                <TextFieldErrorMessage>{field.error}</TextFieldErrorMessage>
              </TextFieldRoot>
            )}
          </Field>
          <Field name="password">
            {(field, props) => (
              <TextFieldRoot
                validationState={validationState(field.error)}
                value={field.value}
                class="space-y-2"
              >
                <TextFieldLabel>Password</TextFieldLabel>
                <TextFieldInput
                  {...props}
                  placeholder="Password"
                  type="password"
                />
                <TextFieldErrorMessage>{field.error}</TextFieldErrorMessage>
              </TextFieldRoot>
            )}
          </Field>
          <Field name="location">
            {(field) => <LocationsField form={form} field={field} />}
          </Field>
          <Field name="features.array" type="string[]">
            {(field) => <FeaturesField form={form} field={field} />}
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

function UpdateForm(props: { onClose: () => void, id: string }) {
  const device = createAsync(() => getDevice(props.id))
  const refetchDevice = () => revalidate(getDevice.keyFor(props.id))

  return (
    <ErrorBoundary fallback={(e) => <PageError error={e} />}>
      <Suspense fallback={<Skeleton class="h-32" />}>
        <Show when={device()}>
          <UpdateFormForm onClose={props.onClose} device={device()!} refetchDevice={refetchDevice} />
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

function UpdateFormForm(props: { onClose: () => void | Promise<void>, device: GetDeviceResp, refetchDevice: () => Promise<void> }) {
  const formInitialValues = (): UpdateForm => ({
    ...props.device,
    features: { array: props.device.features || [] },
    newPassword: ""
  })
  const [form, { Field, Form }] = createForm<UpdateForm>({
    initialValues: formInitialValues()
  });
  const [refreshFormLoading, refreshForm] = createLoading(() => props.refetchDevice()
    .then(() => reset(form, { initialValues: formInitialValues() })))
  const submitForm = (data: UpdateForm) => useClient()
    .admin.updateDevice({ ...data, features: data.features.array })
    .then(() => revalidate())
    .then(props.onClose)
    .catch(throwAsFormError)
  const formDisabled = () => refreshFormLoading() || form.submitting

  return (
    <Form class="flex flex-col gap-4" onSubmit={submitForm}>
      <Field name="id" type="number">
        {(field, props) => <input {...props} type="hidden" value={field.value} />}
      </Field>
      <Field name="name">
        {(field, props) => (
          <TextFieldRoot
            validationState={validationState(field.error)}
            value={field.value}
            class="space-y-2"
          >
            <TextFieldLabel>Name</TextFieldLabel>
            <TextFieldInput
              {...props}
              placeholder="Name"
            />
            <TextFieldErrorMessage>{field.error}</TextFieldErrorMessage>
          </TextFieldRoot>
        )}
      </Field>
      <Field name="url" validate={required("Please enter a URL.")}>
        {(field, props) => (
          <TextFieldRoot
            validationState={validationState(field.error)}
            value={field.value}
            class="space-y-2"
          >
            <TextFieldLabel>URL</TextFieldLabel>
            <TextFieldInput
              {...props}
              placeholder="URL"
            />
            <TextFieldErrorMessage>{field.error}</TextFieldErrorMessage>
          </TextFieldRoot>
        )}
      </Field>
      <Field name="username">
        {(field, props) => (
          <TextFieldRoot
            validationState={validationState(field.error)}
            value={field.value}
            class="space-y-2"
          >
            <TextFieldLabel>Username</TextFieldLabel>
            <TextFieldInput
              {...props}
              placeholder="Username"
            />
            <TextFieldErrorMessage>{field.error}</TextFieldErrorMessage>
          </TextFieldRoot>
        )}
      </Field>
      <Field name="newPassword">
        {(field, props) => (
          <TextFieldRoot
            validationState={validationState(field.error)}
            value={field.value}
            class="space-y-2"
          >
            <TextFieldLabel>New password</TextFieldLabel>
            <TextFieldInput
              {...props}
              placeholder="New password"
              type="password"
            />
            <TextFieldErrorMessage>{field.error}</TextFieldErrorMessage>
          </TextFieldRoot>
        )}
      </Field>
      <Field name="location">
        {(field) => <LocationsField form={form} field={field} />}
      </Field>
      <Field name="features.array" type="string[]">
        {(field) => <FeaturesField form={form} field={field} />}
      </Field>
      <div class="flex flex-col gap-4 sm:flex-row-reverse">
        <Button type="submit" disabled={formDisabled()} class="sm:flex-1">
          <Show when={!form.submitting} fallback="Updating device">Update device</Show>
        </Button>
        <Button type="button" onClick={refreshForm} disabled={formDisabled()} variant="secondary">Refresh</Button>
      </div>
      <FormMessage form={form} />
    </Form>
  )
}

function LocationsField(props: { form: FormStore<any, undefined>, field: FieldStore<any, any> }) {
  const locations = createAsync(() => getListLocations())

  return (
    <Show when={locations()}>
      <SelectRoot<string>
        validationState={validationState(props.field.error)}
        value={props.field.value}
        onChange={setFormValue(props.form, props.field)}
        options={locations()!}
        placeholder="Location"
        itemComponent={props => (
          <SelectItem item={props.item}>
            {props.item.rawValue}
          </SelectItem>
        )}
        virtualized
        class="space-y-2"
      >
        <SelectLabel>Location</SelectLabel>
        <SelectTrigger>
          <SelectValue<string>>
            {state => state.selectedOption()}
          </SelectValue>
        </SelectTrigger>
        <SelectErrorMessage>{props.field.error}</SelectErrorMessage>
        <SelectPortal>
          <SelectContent>
            <SelectListBoxVirtual options={locations()!} />
          </SelectContent>
        </SelectPortal>
      </SelectRoot>
    </Show>
  )
}

function FeaturesField(props: { form: FormStore<any, undefined>, field: FieldStore<any, any> }) {
  const deviceFeatures = createAsync(() => getListDeviceFeatures())

  return (
    <Show when={deviceFeatures()}>
      <SelectRoot<ListDeviceFeaturesResp_Item>
        validationState={validationState(props.field.error)}
        value={deviceFeatures()?.filter(v => props.field.value?.includes(v.value))}
        onChange={(v) => setValue(props.form, props.field.name, v.map(v => v.value))}
        options={deviceFeatures()!}
        optionValue="value"
        optionTextValue="name"
        placeholder="Features"
        itemComponent={props => (
          <SelectItem item={props.item} >
            <div class="flex gap-2">
              {props.item.rawValue.name} <p class="text-muted-foreground">{props.item.rawValue.description}</p>
            </div>
          </SelectItem>
        )}
        multiple
        class="space-y-2"
      >
        <SelectLabel>Features</SelectLabel>
        <SelectTrigger>
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
        <SelectErrorMessage>{props.field.error}</SelectErrorMessage>
        <SelectPortal>
          <SelectContent>
            <SelectListbox />
          </SelectContent>
        </SelectPortal>
      </SelectRoot>
    </Show>
  )
}

export default AdminDevices
