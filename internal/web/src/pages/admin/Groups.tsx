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

const actionDeleteGroup = action((id: bigint) => useClient()
  .admin.deleteGroup({ id })
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

  const previousDisabled = () => data()?.pageResult?.previousPage == data()?.pageResult?.page
  const previous = () => !previousDisabled() && setSearchParams({ page: data()?.pageResult?.previousPage.toString() } as AdminGroupsPageSearchParams)
  const nextDisabled = () => data()?.pageResult?.nextPage == data()?.pageResult?.page
  const next = () => !nextDisabled() && setSearchParams({ page: data()?.pageResult?.nextPage.toString() } as AdminGroupsPageSearchParams)
  const toggleSort = (field: string) => {
    const sort = nextSort(data()?.sort, field)
    return setSearchParams({ sort: sort.field, order: encodeOrder(sort.order) } as AdminGroupsPageSearchParams)
  }

  const [createFormOpen, setCreateFormOpen] = createSignal(false);

  const [updateGroupFormID, setUpdateGroupFormID] = createSignal<bigint>(BigInt(0))

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
                {(group) => {
                  const onClick = () => navigate(`./${group.id}`)

                  const [deleteGroupAlertOpen, setDeleteGroupAlertOpen] = createSignal(false)
                  const deleteGroupSubmission = useSubmission(actionDeleteGroup)
                  const deleteGroupAction = useAction(actionDeleteGroup)
                  const deleteGroup = () => deleteGroupAction(group.id).then(() => setDeleteGroupAlertOpen(false))


                  return (
                    <>
                      <AlertDialogRoot open={deleteGroupAlertOpen()} onOpenChange={setDeleteGroupAlertOpen}>
                        <AlertDialogContent>
                          <AlertDialogHeader>
                            <AlertDialogTitle>Are you sure you wish to delete {group.name}?</AlertDialogTitle>
                          </AlertDialogHeader>
                          <AlertDialogFooter>
                            <AlertDialogCancel>Cancel</AlertDialogCancel>
                            <AlertDialogAction variant="destructive" disabled={deleteGroupSubmission.pending} onClick={deleteGroup}>
                              Delete
                            </AlertDialogAction>
                          </AlertDialogFooter>
                        </AlertDialogContent>
                      </AlertDialogRoot>

                      <TableRow class="">
                        <TableCell onClick={onClick} class="cursor-pointer select-none">{group.name}</TableCell>
                        <TableCell onClick={onClick} class="text-nowrap cursor-pointer select-none whitespace-nowrap">{group.userCount.toString()}</TableCell>
                        <TableCell onClick={onClick} class="text-nowrap cursor-pointer select-none whitespace-nowrap">{formatDate(parseDate(group.createdAtTime))}</TableCell>
                        <TableCell class="py-0">
                          <div class="flex gap-2">
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
                            <div class="flex items-center justify-end">
                              <DropdownMenuRoot placement="bottom-end">
                                <DropdownMenuTrigger class="hover:bg-accent hover:text-accent-foreground rounded p-1" title="Actions">
                                  <RiSystemMore2Line class="h-5 w-5" />
                                </DropdownMenuTrigger>
                                <DropdownMenuPortal>
                                  <DropdownMenuContent>
                                    <DropdownMenuItem onSelect={() => setUpdateGroupFormID(group.id)}>
                                      Edit
                                    </DropdownMenuItem>
                                    <DropdownMenuItem onSelect={() => setDeleteGroupAlertOpen(true)}>
                                      Delete
                                    </DropdownMenuItem>
                                    <DropdownMenuArrow />
                                  </DropdownMenuContent>
                                </DropdownMenuPortal>
                              </DropdownMenuRoot>
                            </div>
                          </div>
                        </TableCell>
                      </TableRow>
                    </>
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
  .admin.createGroup(form)
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
        <Show when={createGroupForm.submitting} fallback={<>Create group</>}>
          Creating group
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
  .admin.updateGroup(form)
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
        initialValues: { ...res }
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
            <Show when={updateGroupForm.submitting} fallback={<>Update group</>}>
              Updating group
            </Show>
          </Button>
          <FormMessage form={updateGroupForm} />
        </Form>
      </Suspense>
    </ErrorBoundary>
  )
}
