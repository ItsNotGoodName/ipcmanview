import { action, createAsync, revalidate, useAction, useNavigate, useSearchParams, useSubmission } from "@solidjs/router";
import { AlertDialogAction, AlertDialogCancel, AlertDialogModal, AlertDialogDescription, AlertDialogFooter, AlertDialogHeader, AlertDialogRoot, AlertDialogTitle, } from "~/ui/AlertDialog";
import { DropdownMenuArrow, DropdownMenuContent, DropdownMenuItem, DropdownMenuPortal, DropdownMenuRoot, } from "~/ui/DropdownMenu";
import { AdminGroupsPageSearchParams, getAdminGroupsPage } from "./Groups.data";
import { ErrorBoundary, For, Show, Suspense, createResource, createSignal } from "solid-js";
import { RiSystemLockLine, } from "solid-icons/ri";
import { Button } from "~/ui/Button";
import { catchAsToast, createPagePagination, createRowSelection, createToggleSortField, formatDate, parseDate, syncForm, throwAsFormError } from "~/lib/utils";
import { parseOrder } from "~/lib/utils";
import { TableBody, TableCaption, TableCell, TableHead, TableHeader, TableRoot, TableRow, } from "~/ui/Table";
import { useClient } from "~/providers/client";
import { createForm, required, reset } from "@modular-forms/solid";
import { FieldControl, FieldLabel, FieldMessage, FieldRoot, FormMessage } from "~/ui/Form";
import { Input } from "~/ui/Input";
import { Textarea } from "~/ui/Textarea";
import { DialogModal, DialogHeader, DialogContent, DialogOverlay, DialogPortal, DialogRoot, DialogTitle, } from "~/ui/Dialog";
import { CheckboxControl, CheckboxInput, CheckboxLabel, CheckboxRoot } from "~/ui/Checkbox";
import { Skeleton } from "~/ui/Skeleton";
import { PageError } from "~/ui/Page";
import { TooltipContent, TooltipRoot, TooltipTrigger } from "~/ui/Tooltip";
import { LayoutNormal } from "~/ui/Layout";
import { SetGroupDisableReq } from "~/twirp/rpc";
import { Crud } from "~/components/Crud";
import { Shared } from "~/components/Shared";

const actionDelete = action((ids: bigint[]) => useClient()
  .admin.deleteGroup({ ids })
  .then(() => revalidate(getAdminGroupsPage.key))
  .catch(catchAsToast))

const actionSetDisable = action((input: SetGroupDisableReq) => useClient()
  .admin.setGroupDisable(input)
  .then(() => revalidate(getAdminGroupsPage.key))
  .catch(catchAsToast))

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
  const [openCreateForm, setOpenCreateForm] = createSignal(false);

  // Update
  const [openUpdateForm, setOpenUpdateForm] = createSignal<bigint>(BigInt(0))

  // Delete
  const deleteGroupSubmission = useSubmission(actionDelete)
  const deleteGroupAction = useAction(actionDelete)
  // Single
  const [openDeleteConfirm, setOpenDeleteConfirm] = createSignal<{ name: string, id: bigint } | undefined>()
  const deleteSubmit = () => deleteGroupAction([openDeleteConfirm()!.id])
    .then(() => setOpenDeleteConfirm(undefined))
  // Multiple
  const [openDeleteMultipleConfirm, setDeleteMultipleConfirm] = createSignal(false)
  const deleteMultipleSubmit = () => deleteGroupAction(rowSelection.selections())
    .then(() => setDeleteMultipleConfirm(false))

  // Disable
  const setDisableSubmission = useSubmission(actionSetDisable)
  const setDisableAction = useAction(actionSetDisable)
  const setDisableMultipleSubmit = (disable: boolean) => setDisableAction({ items: rowSelection.selections().map(v => ({ id: v, disable })) })
    .then(() => rowSelection.setAll(false))

  return (
    <LayoutNormal class="max-w-4xl">
      <DialogRoot open={openCreateForm()} onOpenChange={setOpenCreateForm}>
        <DialogPortal>
          <DialogOverlay />
          <DialogModal>
            <DialogHeader>
              <DialogTitle>Create group</DialogTitle>
            </DialogHeader>
            <DialogContent>
              <CreateForm close={() => setOpenCreateForm(false)} />
            </DialogContent>
          </DialogModal>
        </DialogPortal>
      </DialogRoot>

      <DialogRoot open={openUpdateForm() != BigInt(0)} onOpenChange={() => setOpenUpdateForm(BigInt(0))}>
        <DialogPortal>
          <DialogOverlay />
          <DialogModal>
            <DialogHeader>
              <DialogTitle>Update group</DialogTitle>
            </DialogHeader>
            <DialogContent>
              <UpdateForm close={() => setOpenUpdateForm(BigInt(0))} id={openUpdateForm()} />
            </DialogContent>
          </DialogModal>
        </DialogPortal>
      </DialogRoot>

      <AlertDialogRoot open={openDeleteConfirm() != undefined} onOpenChange={() => setOpenDeleteConfirm(undefined)}>
        <AlertDialogModal>
          <AlertDialogHeader>
            <AlertDialogTitle>Are you sure you wish to delete {openDeleteConfirm()?.name}?</AlertDialogTitle>
          </AlertDialogHeader>
          <AlertDialogFooter>
            <AlertDialogCancel>Cancel</AlertDialogCancel>
            <AlertDialogAction variant="destructive" disabled={deleteGroupSubmission.pending} onClick={deleteSubmit}>
              Delete
            </AlertDialogAction>
          </AlertDialogFooter>
        </AlertDialogModal>
      </AlertDialogRoot>

      <AlertDialogRoot open={openDeleteMultipleConfirm()} onOpenChange={setDeleteMultipleConfirm}>
        <AlertDialogModal>
          <AlertDialogHeader>
            <AlertDialogTitle>Are you sure you wish to delete {rowSelection.selections().length} groups?</AlertDialogTitle>
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
            <AlertDialogAction variant="destructive" disabled={deleteGroupSubmission.pending} onClick={deleteMultipleSubmit}>
              Delete
            </AlertDialogAction>
          </AlertDialogFooter>
        </AlertDialogModal>
      </AlertDialogRoot>

      <Shared.Title>Groups</Shared.Title>

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
                    name="userCount"
                    onClick={toggleSort}
                    sort={data()?.sort}
                  >
                    Users
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
                        <DropdownMenuItem onSelect={() => setOpenCreateForm(true)}>
                          Create
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
                          disabled={rowSelection.selections().length == 0 || deleteGroupSubmission.pending}
                          onSelect={() => setDeleteMultipleConfirm(true)}
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
                  const toggleGroupDisable = () => setDisableAction({ items: [{ id: item.id, disable: !item.disabled }] })

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
                      <TableCell class="cursor-pointer select-none" onClick={onClick}>{item.userCount.toString()}</TableCell>
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
                              <DropdownMenuItem onSelect={() => setOpenUpdateForm(item.id)}>
                                Edit
                              </DropdownMenuItem>
                              <DropdownMenuItem
                                disabled={setDisableSubmission.pending}
                                onSelect={toggleGroupDisable}
                              >
                                <Show when={item.disabled} fallback={<>Disable</>}>
                                  Enable
                                </Show>
                              </DropdownMenuItem>
                              <DropdownMenuItem
                                disabled={deleteGroupSubmission.pending}
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

type CreateForm = {
  name: string
  description: string
}

const actionCreateForm = action((data: CreateForm) => useClient()
  .admin.createGroup(data)
  .then(() => revalidate(getAdminGroupsPage.key))
  .catch(throwAsFormError))

function CreateForm(props: { close: () => void }) {
  const [addMore, setAddMore] = createSignal(false)

  const [form, { Field, Form }] = createForm<CreateForm>({ initialValues: { name: "", description: "" } });
  const action = useAction(actionCreateForm)
  const submit = async (data: CreateForm) => {
    await action(data)
    if (addMore()) {
      reset(form)
    } else {
      props.close()
    }
  }

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
      <Button type="submit" disabled={form.submitting}>
        <Show when={!form.submitting} fallback={<>Creating group</>}>
          Create group
        </Show>
      </Button>
      <FormMessage form={form} />
      <CheckboxRoot checked={addMore()} onChange={setAddMore}>
        <CheckboxInput />
        <CheckboxControl />
        <CheckboxLabel>Add more</CheckboxLabel>
      </CheckboxRoot>
    </Form>
  )
}

type UpdateForm = {
  id: any
  name: string
  description: string
}

const actionUpdateForm = action((data: UpdateForm) => useClient()
  .admin.updateGroup(data)
  .then(() => revalidate(getAdminGroupsPage.key))
  .catch(throwAsFormError))

function UpdateForm(props: { close: () => void, id: bigint }) {
  const [form, { Field, Form }] = createForm<UpdateForm>();
  const action = useAction(actionUpdateForm)
  const submit = (data: UpdateForm) => action(data)
    .then(() => props.close())

  const [data] = createResource(() => props.id != BigInt(0),
    () => useClient().admin.getGroup({ id: props.id })
      .then((data) => data.response satisfies UpdateForm))
  const disabled = syncForm(form, data)

  return (
    <ErrorBoundary fallback={(e) => <PageError error={e} />}>
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
                  disabled={disabled()}
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
                  disabled={disabled()}
                >
                  {field.value}
                </Textarea>
              </FieldControl>
              <FieldMessage field={field} />
            </FieldRoot>
          )}
        </Field>
        <Button type="submit" disabled={disabled() || form.submitting}>
          <Show when={!form.submitting} fallback={<>Updating group</>}>
            Update group
          </Show>
        </Button>
        <FormMessage form={form} />
      </Form>
    </ErrorBoundary>
  )
}
