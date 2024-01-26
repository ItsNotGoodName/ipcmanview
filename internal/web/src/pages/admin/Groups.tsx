import { action, createAsync, revalidate, useAction, useNavigate, useSearchParams, useSubmission } from "@solidjs/router";
import { AlertDialogAction, AlertDialogCancel, AlertDialogContent, AlertDialogFooter, AlertDialogHeader, AlertDialogRoot, AlertDialogTitle, } from "~/ui/AlertDialog";
import { DropdownMenuArrow, DropdownMenuContent, DropdownMenuItem, DropdownMenuPortal, DropdownMenuRoot, DropdownMenuTrigger } from "~/ui/DropdownMenu";
import { AdminGroupsPageSearchParams, getAdminGroupsPage, getGroup } from "./Groups.data";
import { ErrorBoundary, For, Show, Suspense, batch, createResource, createSignal } from "solid-js";
import { RiArrowsArrowLeftSLine, RiArrowsArrowRightSLine, RiSystemLockLine, RiSystemMore2Line, } from "solid-icons/ri";
import { Button } from "~/ui/Button";
import { SelectContent, SelectItem, SelectListbox, SelectRoot, SelectTrigger, SelectValue } from "~/ui/Select";
import { catchAsToast, createPagePagination, createRowSelection, createToggleSortField, formatDate, parseDate, syncForm as setupForm, throwAsFormError } from "~/lib/utils";
import { parseOrder } from "~/lib/utils";
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
import { SetGroupDisableReq } from "~/twirp/rpc";

const actionDeleteGroup = action((ids: bigint[]) => useClient()
  .admin.deleteGroup({ ids })
  .then(() => revalidate(getAdminGroupsPage.key))
  .catch(catchAsToast)
)

const actionSetGroupDisable = action((input: SetGroupDisableReq) => useClient()
  .admin.setGroupDisable(input)
  .then(() => revalidate(getAdminGroupsPage.key))
  .catch(catchAsToast)
)

export function AdminGroups() {
  const navigate = useNavigate()
  const [searchParams] = useSearchParams<AdminGroupsPageSearchParams>()
  const data = createAsync(() => getAdminGroupsPage({
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
  const [createFormDialog, setCreateFormDialog] = createSignal(false);

  // Update
  const [updateGroupFormDialog, setUpdateGroupFormDialog] = createSignal<bigint>(BigInt(0))

  // Delete
  const deleteGroupSubmission = useSubmission(actionDeleteGroup)
  const deleteGroupAction = useAction(actionDeleteGroup)
  // Single
  const [deleteGroupSelection, setDeleteGroupSelection] = createSignal<{ name: string, id: bigint } | undefined>()
  const deleteGroupBySelection = () => deleteGroupAction([deleteGroupSelection()!.id])
    .then(() => setDeleteGroupSelection(undefined))
  // Multiple
  const [deleteGroupRowSelector, setDeleteGroupRowSelector] = createSignal(false)
  const deleteGroupByRowSelector = () => deleteGroupAction(rowSelection.selections())
    .then(() => setDeleteGroupRowSelector(false))

  // Disable/Enable
  const setGroupDisableSubmission = useSubmission(actionSetGroupDisable)
  const setGroupDisable = useAction(actionSetGroupDisable)
  const setGroupDisableByRowSelector = (disable: boolean) => setGroupDisable({ items: rowSelection.selections().map(v => ({ id: v, disable })) })
    .then(() => rowSelection.setAll(false))
  const setGroupDisableDisabled = (disable: boolean) => {
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
            <CreateGroupForm setOpen={setCreateFormDialog} />
          </DialogContent>
        </DialogPortal>
      </DialogRoot>

      <DialogRoot open={updateGroupFormDialog() != BigInt(0)} onOpenChange={() => setUpdateGroupFormDialog(BigInt(0))}>
        <DialogPortal>
          <DialogOverlay />
          <DialogContent>
            <DialogHeader>
              <DialogCloseButton />
              <DialogTitle>Update group</DialogTitle>
            </DialogHeader>
            <UpdateGroupForm setOpen={() => setUpdateGroupFormDialog(BigInt(0))} id={updateGroupFormDialog()} />
          </DialogContent>
        </DialogPortal>
      </DialogRoot>

      <AlertDialogRoot open={deleteGroupSelection() != undefined} onOpenChange={() => setDeleteGroupSelection(undefined)}>
        <AlertDialogContent>
          <AlertDialogHeader>
            <AlertDialogTitle>Are you sure you wish to delete {deleteGroupSelection()?.name}?</AlertDialogTitle>
          </AlertDialogHeader>
          <AlertDialogFooter>
            <AlertDialogCancel>Cancel</AlertDialogCancel>
            <AlertDialogAction variant="destructive" disabled={deleteGroupSubmission.pending} onClick={deleteGroupBySelection}>
              Delete
            </AlertDialogAction>
          </AlertDialogFooter>
        </AlertDialogContent>
      </AlertDialogRoot>

      <AlertDialogRoot open={deleteGroupRowSelector()} onOpenChange={setDeleteGroupRowSelector}>
        <AlertDialogContent>
          <AlertDialogHeader>
            <AlertDialogTitle>Are you sure you wish to delete {rowSelection.selections().length} groups?</AlertDialogTitle>
          </AlertDialogHeader>
          <AlertDialogFooter>
            <AlertDialogCancel>Cancel</AlertDialogCancel>
            <AlertDialogAction variant="destructive" disabled={deleteGroupSubmission.pending} onClick={deleteGroupByRowSelector}>
              Delete
            </AlertDialogAction>
          </AlertDialogFooter>
        </AlertDialogContent>
      </AlertDialogRoot>

      <div class="text-xl">Groups</div>
      <Seperator />

      <ErrorBoundary fallback={(e: Error) => <PageError error={e} />}>
        <Suspense fallback={<Skeleton class="h-32" />}>
          <div class="flex justify-between gap-2">
            <SelectRoot
              class="w-20"
              value={data()?.pageResult?.perPage}
              onChange={pagination.setPerPage}
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
                            onSelect={() => setGroupDisableByRowSelector(true)}
                            disabled={setGroupDisableDisabled(true)}
                          >
                            Disable
                          </DropdownMenuItem>
                          <DropdownMenuItem
                            onSelect={() => setGroupDisableByRowSelector(false)}
                            disabled={setGroupDisableDisabled(false)}
                          >
                            Enable
                          </DropdownMenuItem>
                          <DropdownMenuItem
                            onSelect={() => setDeleteGroupRowSelector(true)}
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
                  const toggleGroupDisable = () => setGroupDisable({ items: [{ id: item.id, disable: !item.disabled }] })

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
                                <DropdownMenuItem onSelect={() => setUpdateGroupFormDialog(item.id)}>
                                  Edit
                                </DropdownMenuItem>
                                <DropdownMenuItem disabled={setGroupDisableSubmission.pending} onSelect={toggleGroupDisable}>
                                  <Show when={item.disabled} fallback={<>Disable</>}>
                                    Enable
                                  </Show>
                                </DropdownMenuItem>
                                <DropdownMenuItem onSelect={() => setDeleteGroupSelection(item)}>
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

type CreateGroupForm = {
  name: string
  description: string
}

const actionCreateGroupForm = action((form: CreateGroupForm) => useClient()
  .admin.createGroup({ model: form })
  .then(() => revalidate(getAdminGroupsPage.key))
  .catch(throwAsFormError)
)

function CreateGroupForm(props: { setOpen: (value: boolean) => void }) {
  const [addMore, setAddMore] = createSignal(false)

  const [createGroupForm, { Field, Form }] = createForm<CreateGroupForm>({ initialValues: { name: "", description: "" } });
  const createGroupFormAction = useAction(actionCreateGroupForm)
  const submit = (form: CreateGroupForm) => createGroupFormAction(form)
    .then(() => batch(() => {
      props.setOpen(addMore())
      reset(createGroupForm)
    }))

  return (
    <Form class="flex flex-col gap-4" onSubmit={submit}>
      <input class="hidden" type="text" name="username" autocomplete="username" />
      <Field name="name" validate={required("Please enter a name.")}>
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
      <Field name="description">
        {(field, props) => (
          <FieldRoot class="gap-1.5">
            <FieldLabel field={field}>Description</FieldLabel>
            <FieldControl field={field}>
              <Textarea
                {...props}
                value={field.value}
                placeholder="Description"
              />
            </FieldControl>
            <FieldMessage field={field} />
          </FieldRoot>
        )}
      </Field>
      <Button type="submit" disabled={createGroupForm.submitting}>
        <Show when={!createGroupForm.submitting} fallback={<>Creating group</>}>
          Create group
        </Show>
      </Button>
      <FormMessage form={createGroupForm} />
      <CheckboxRoot checked={addMore()} onChange={setAddMore}>
        <CheckboxInput />
        <CheckboxControl />
        <CheckboxLabel>Add more</CheckboxLabel>
      </CheckboxRoot>
    </Form>
  )
}

type UpdateGroupForm = {
  id: any
  name: string
  description: string
}

const actionUpdateGroupForm = action((form: UpdateGroupForm) => useClient()
  .admin.updateGroup({ id: form.id, model: form })
  .then(() => revalidate(getAdminGroupsPage.key))
  .catch(throwAsFormError)
)

function UpdateGroupForm(props: { setOpen: (value: boolean) => void, id: bigint }) {
  const [updateGroupForm, { Field, Form }] = createForm<UpdateGroupForm>();
  const updateGroupFormAction = useAction(actionUpdateGroupForm)
  const submit = (form: UpdateGroupForm) => updateGroupFormAction(form)
    .then(() => props.setOpen(false))
  const [form] = createResource(() => getGroup(props.id)
    .then((data) => setupForm(updateGroupForm, { ...data, ...data.model })))

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
        <Button type="submit" disabled={form.loading || updateGroupForm.submitting}>
          <Show when={!updateGroupForm.submitting} fallback={<>Updating group</>}>
            Update group
          </Show>
        </Button>
        <FormMessage form={updateGroupForm} />
      </Form>
    </Show>
  )
}
