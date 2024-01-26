import { action, createAsync, revalidate, useAction, useNavigate, useSearchParams, useSubmission } from "@solidjs/router";
import { AlertDialogAction, AlertDialogCancel, AlertDialogContent, AlertDialogFooter, AlertDialogHeader, AlertDialogRoot, AlertDialogTitle, } from "~/ui/AlertDialog";
import { DropdownMenuArrow, DropdownMenuContent, DropdownMenuItem, DropdownMenuPortal, DropdownMenuRoot, DropdownMenuTrigger } from "~/ui/DropdownMenu";
import { AdminGroupsPageSearchParams, getAdminGroupsPage, getGroup } from "./Groups.data";
import { ErrorBoundary, For, Show, Suspense, batch, createSignal } from "solid-js";
import { RiArrowsArrowLeftSLine, RiArrowsArrowRightSLine, RiSystemLockLine, RiSystemMore2Line, } from "solid-icons/ri";
import { Button } from "~/ui/Button";
import { SelectContent, SelectItem, SelectListbox, SelectRoot, SelectTrigger, SelectValue } from "~/ui/Select";
import { catchAsToast, formatDate, parseDate, throwAsFormError } from "~/lib/utils";
import { encodeOrder, nextSort, parseOrder } from "~/lib/utils";
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
import { paginateOptions } from "~/lib/utils";
import { LayoutNormal } from "~/ui/Layout";
import { SetGroupDisableReq } from "~/twirp/rpc";
import { createRowSelector } from "~/lib/row";

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
  const [searchParams, setSearchParams] = useSearchParams<AdminGroupsPageSearchParams>()
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
  const rowSelector = createRowSelector(() => data()?.items.map(v => v.id) || [])

  // List
  const previousDisabled = () => data()?.pageResult?.previousPage == data()?.pageResult?.page
  const previous = () => !previousDisabled() && setSearchParams({ page: data()?.pageResult?.previousPage.toString() } as AdminGroupsPageSearchParams)
  const nextDisabled = () => data()?.pageResult?.nextPage == data()?.pageResult?.page
  const next = () => !nextDisabled() && setSearchParams({ page: data()?.pageResult?.nextPage.toString() } as AdminGroupsPageSearchParams)
  const toggleSort = (field: string) => {
    const sort = nextSort(data()?.sort, field)
    return setSearchParams({ sort: sort.field, order: encodeOrder(sort.order) } as AdminGroupsPageSearchParams)
  }

  // Create
  const [createFormOpen, setCreateFormOpen] = createSignal(false);

  // Update
  const [updateGroupFormID, setUpdateGroupFormID] = createSignal<bigint>(BigInt(0))

  // Delete
  const deleteGroupSubmission = useSubmission(actionDeleteGroup)
  const deleteGroupAction = useAction(actionDeleteGroup)

  const [deleteGroupSelection, setDeleteGroupSelection] = createSignal<{ name: string, id: bigint } | undefined>()
  const deleteGroupBySelection = () => deleteGroupAction([deleteGroupSelection()!.id])
    .then(() => setDeleteGroupSelection(undefined))

  const [deleteGroupRowSelector, setDeleteGroupRowSelector] = createSignal(false)
  const deleteGroupByRowSelector = () => deleteGroupAction(rowSelector.selected())
    .then(() => setDeleteGroupRowSelector(false))

  // Disable/Enable
  const setGroupDisableSubmission = useSubmission(actionSetGroupDisable)
  const setGroupDisableAction = useAction(actionSetGroupDisable)
  const setGroupDisableActionByRowSelector = (disable: boolean) => setGroupDisableAction({ items: rowSelector.selected().map(v => ({ id: v, disable })) })
    .then(() => rowSelector.checkAll(false))

  return (
    <LayoutNormal>
      <DialogRoot open={createFormOpen()} onOpenChange={setCreateFormOpen}>
        <DialogPortal>
          <DialogOverlay />
          <DialogContent>
            <DialogHeader>
              <DialogCloseButton />
              <DialogTitle>Create group</DialogTitle>
            </DialogHeader>
            <CreateGroupForm setOpen={setCreateFormOpen} />
          </DialogContent>
        </DialogPortal>
      </DialogRoot>

      <DialogRoot open={updateGroupFormID() != BigInt(0)} onOpenChange={() => setUpdateGroupFormID(BigInt(0))}>
        <DialogPortal>
          <DialogOverlay />
          <DialogContent>
            <DialogHeader>
              <DialogCloseButton />
              <DialogTitle>Update group</DialogTitle>
            </DialogHeader>
            <UpdateGroupForm setOpen={() => setUpdateGroupFormID(BigInt(0))} id={updateGroupFormID()} />
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
            <AlertDialogTitle>Are you sure you wish to delete {rowSelector.selected().length} groups?</AlertDialogTitle>
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
              onChange={(value) => value && setSearchParams({ page: 1, perPage: value })}
              options={paginateOptions}
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
                  <CheckboxRoot checked={rowSelector.multiple()} indeterminate={rowSelector.indeterminate()} onChange={(v) => rowSelector.checkAll(v)}>
                    <CheckboxControl />
                  </CheckboxRoot>
                </TableHead>
                <TableHead class="w-full">
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
                          <DropdownMenuItem onSelect={() => setCreateFormOpen(true)}>
                            Create
                          </DropdownMenuItem>
                          <DropdownMenuItem onSelect={() => setGroupDisableActionByRowSelector(true)} disabled={rowSelector.selected().length == 0}>
                            Disable
                          </DropdownMenuItem>
                          <DropdownMenuItem onSelect={() => setGroupDisableActionByRowSelector(false)} disabled={rowSelector.selected().length == 0}>
                            Enable
                          </DropdownMenuItem>
                          <DropdownMenuItem onSelect={() => setDeleteGroupRowSelector(true)} disabled={rowSelector.selected().length == 0}>
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
                {(group, index) => {
                  const navigateToGroup = () => navigate(`./${group.id}`)
                  const toggleGroupDisable = () => setGroupDisableAction({ items: [{ id: group.id, disable: !group.disabled }] })

                  return (
                    <TableRow class="">
                      <TableHead>
                        <CheckboxRoot checked={rowSelector.selections[index()]?.checked} onChange={(v) => rowSelector.check(group.id, v)}>
                          <CheckboxControl />
                        </CheckboxRoot>
                      </TableHead>
                      <TableCell onClick={navigateToGroup} class="cursor-pointer select-none">{group.name}</TableCell>
                      <TableCell onClick={navigateToGroup} class="text-nowrap cursor-pointer select-none whitespace-nowrap">{group.userCount.toString()}</TableCell>
                      <TableCell onClick={navigateToGroup} class="text-nowrap cursor-pointer select-none whitespace-nowrap">{formatDate(parseDate(group.createdAtTime))}</TableCell>
                      <TableCell class="py-0">
                        <div class="flex justify-end gap-2">
                          <Show when={group.disabled}>
                            <TooltipRoot>
                              <TooltipTrigger class="p-1">
                                <RiSystemLockLine class="h-5 w-5" />
                              </TooltipTrigger>
                              <TooltipContent>
                                Disabled since {formatDate(parseDate(group.disabledAtTime))}
                              </TooltipContent>
                            </TooltipRoot>
                          </Show>
                          <DropdownMenuRoot placement="bottom-end">
                            <DropdownMenuTrigger class="hover:bg-accent hover:text-accent-foreground rounded p-1" title="Actions">
                              <RiSystemMore2Line class="h-5 w-5" />
                            </DropdownMenuTrigger>
                            <DropdownMenuPortal>
                              <DropdownMenuContent>
                                <DropdownMenuItem disabled={setGroupDisableSubmission.pending} onSelect={toggleGroupDisable}>
                                  <Show when={group.disabled} fallback={<>Disable</>}>
                                    Enable
                                  </Show>
                                </DropdownMenuItem>
                                <DropdownMenuItem onSelect={() => setUpdateGroupFormID(group.id)}>
                                  Update
                                </DropdownMenuItem>
                                <DropdownMenuItem onSelect={() => setDeleteGroupSelection(group)}>
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
    .then(() => {
      batch(() => {
        props.setOpen(addMore())
        reset(createGroupForm)
      })
    })

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
  id: BigInt | any
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
  const submit = (form: UpdateGroupForm) => updateGroupFormAction(form).then(() => props.setOpen(false))
  const disabled = createAsync(async () => {
    if (props.id == BigInt(0))
      return false

    return getGroup(props.id).then(res => {
      if (updateGroupForm.submitted) {
        return false
      }

      reset(updateGroupForm, {
        initialValues: { ...res, ...res.model }
      })
      return false
    })
  })

  return (
    <ErrorBoundary fallback={(e: Error) => <PageError error={e} />}>
      <Suspense fallback={<Skeleton class="h-32" />}>
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
                    disabled={disabled()}
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
                    disabled={disabled()}
                    placeholder="Description"
                  >
                    {field.value}
                  </Textarea>
                </FieldControl>
                <FieldMessage field={field} />
              </FieldRoot>
            )}
          </Field>
          <Button type="submit" disabled={disabled() || updateGroupForm.submitting}>
            <Show when={!updateGroupForm.submitting} fallback={<>Updating group</>}>
              Update group
            </Show>
          </Button>
          <FormMessage form={updateGroupForm} />
        </Form>
      </Suspense>
    </ErrorBoundary>
  )
}
