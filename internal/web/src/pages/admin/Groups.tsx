import { action, createAsync, revalidate, useAction, useNavigate, useSearchParams, useSubmission } from "@solidjs/router";
import { DropdownMenuArrow, DropdownMenuContent, DropdownMenuItem, DropdownMenuPortal, DropdownMenuRoot, DropdownMenuTrigger } from "~/ui/DropdownMenu";
import { AdminGroupsPageSearchParams, getAdminGroupsPage } from "./Groups.data";
import { ErrorBoundary, For, Show, Suspense, createSignal } from "solid-js";
import { RiArrowsArrowLeftSLine, RiArrowsArrowRightSLine, RiSystemDeleteBinLine, RiSystemLockLine, RiSystemMore2Line, } from "solid-icons/ri";
import { Button } from "~/ui/Button";
import { SelectContent, SelectItem, SelectListbox, SelectRoot, SelectTrigger, SelectValue } from "~/ui/Select";
import { catchAsToast, formatDate, parseDate, throwAsFormError } from "~/lib/utils";
import { Order, } from "~/twirp/rpc";
import { encodeOrder, nextOrder, parseOrder } from "~/lib/order";
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
import { ConfirmButton } from "~/ui/Confirm";

const actionDeleteGroup = action((id: bigint) => useClient()
  .admin.deleteGroup({ id })
  .then(() => revalidate(getAdminGroupsPage.key))
  .catch(catchAsToast)
)

export function AdminGroups() {
  const navigate = useNavigate()
  const [searchParams, setSearchParams] = useSearchParams<AdminGroupsPageSearchParams>()
  const groups = createAsync(() => getAdminGroupsPage({
    page: {
      page: Number(searchParams.page) || 1,
      perPage: Number(searchParams.perPage) || 10
    },
    sort: {
      field: searchParams.sort || "",
      order: parseOrder(searchParams.order)
    },
  }))

  const previousDisabled = () => groups()?.pageResult?.previousPage == groups()?.pageResult?.page
  const previous = () => !previousDisabled() && setSearchParams({ page: groups()?.pageResult?.previousPage.toString() } as AdminGroupsPageSearchParams)
  const nextDisabled = () => groups()?.pageResult?.nextPage == groups()?.pageResult?.page
  const next = () => !nextDisabled() && setSearchParams({ page: groups()?.pageResult?.nextPage.toString() } as AdminGroupsPageSearchParams)
  const toggleSort = (value: string) => {
    if (value == groups()?.sort?.field) {
      const order = nextOrder(groups()?.sort?.order ?? Order.ORDER_UNSPECIFIED)

      if (order == Order.ORDER_UNSPECIFIED) {
        return setSearchParams({ sort: undefined, order: undefined })
      }

      return setSearchParams({ sort: value, order: encodeOrder(order) } as AdminGroupsPageSearchParams)
    }

    return setSearchParams({ sort: value, order: encodeOrder(Order.DESC) } as AdminGroupsPageSearchParams)
  }

  const [createFormOpen, setCreateFormOpen] = createSignal(false);

  return (
    <div class="flex justify-center p-4">
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

      <div class="flex w-full max-w-4xl flex-col gap-2">
        <div class="text-xl">Groups</div>
        <Seperator />
        <ErrorBoundary fallback={(e: Error) => <PageError error={e} />}>
          <Suspense fallback={<Skeleton class="h-32" />}>
            <div class="flex justify-between gap-2">
              <SelectRoot
                class="w-20"
                value={groups()?.pageResult?.perPage}
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
                  <TableHead class="w-full">
                    <TableSortButton
                      name="name"
                      onClick={toggleSort}
                      sort={groups()?.sort}
                    >
                      Name
                    </TableSortButton>
                  </TableHead>
                  <TableHead>
                    <TableSortButton
                      name="userCount"
                      onClick={toggleSort}
                      sort={groups()?.sort}
                    >
                      Users
                    </TableSortButton>
                  </TableHead>
                  <TableHead>
                    <TableSortButton
                      name="createdAt"
                      onClick={toggleSort}
                      sort={groups()?.sort}
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
                <For each={groups()?.items}>
                  {(group) => {
                    const onClick = () => navigate(`./${group.id}`)

                    const deleteGroupSubmission = useSubmission(actionDeleteGroup)
                    const deleteGroupAction = useAction(actionDeleteGroup)
                    const deleteGroup = () => deleteGroupAction(group.id)

                    return (
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
                            <ConfirmButton
                              class="bg-destructive text-destructive-foreground rounded p-1 disabled:opacity-50"
                              message={`Are you sure wish to delete ${group.name}?`}
                              disabled={deleteGroupSubmission.pending}
                              onYes={deleteGroup}
                              title="Delete"
                            >
                              <RiSystemDeleteBinLine class="h-5 w-5" />
                            </ConfirmButton>
                          </div>
                        </TableCell>
                      </TableRow>
                    )
                  }}
                </For>
              </TableBody>
              <TableCaption>
                <TableMetadata pageResult={groups()?.pageResult} />
              </TableCaption>
            </TableRoot>
          </Suspense>
        </ErrorBoundary>
      </div>
    </div>
  )
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
  const [createGroupForm, { Field, Form }] = createForm<CreateGroupForm>({ initialValues: { name: "", description: "" } });
  const createGroupFormAction = useAction(actionCreateGroupForm)
  const [addMore, setAddMore] = createSignal(false)
  const submit = (form: CreateGroupForm) => createGroupFormAction(form)
    .then(() => {
      props.setOpen(addMore())
      reset(createGroupForm)
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

// type UpdateGroupForm = {
//   name: string
//   description: string
// }
//
// const actionUpdateGroupForm = action((form: UpdateGroupForm) => useClient()
//   .admin.updateGroup(form)
//   .then(() => revalidate(getListGroups.key))
//   .catch(throwAsFormError)
// )
//
// function UpdateGroupForm() {
//   const [updateGroupForm, { Field, Form }] = createForm<UpdateGroupForm>({ initialValues: { name: "", description: "" } });
//   const submit = useAction(actionUpdateGroupForm)
//
//   return (
//     <Form class="flex flex-col gap-4" onSubmit={(form) => submit(form).then(() => reset(updateGroupForm))}>
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
//                 placeholder="Description"
//               >
//                 {field.value}
//               </Textarea>
//             </FieldControl>
//             <FieldMessage field={field} />
//           </FieldRoot>
//         )}
//       </Field>
//       <Button type="submit" disabled={updateGroupForm.submitting}>
//         <Show when={updateGroupForm.submitting} fallback={<>Update group</>}>
//           Updating group
//         </Show>
//       </Button>
//       <FormMessage form={updateGroupForm} />
//     </Form>
//   )
// }
