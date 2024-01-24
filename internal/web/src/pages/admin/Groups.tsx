import { action, createAsync, revalidate, useAction, useNavigate, useSearchParams } from "@solidjs/router";
import { getListGroups } from "./Groups.data";
import { For, ParentProps, Show, Suspense, createSignal } from "solid-js";
import { RiArrowsArrowDownSLine, RiArrowsArrowLeftSLine, RiArrowsArrowRightSLine, RiArrowsArrowUpSLine } from "solid-icons/ri";
import { Button } from "~/ui/Button";
import { SelectContent, SelectItem, SelectListbox, SelectRoot, SelectTrigger, SelectValue } from "~/ui/Select";
import { cn, formatDate, parseDate, throwAsFormError } from "~/lib/utils";
import { Order } from "~/twirp/rpc";
import { encodeOrder, nextOrder, parseOrder } from "~/lib/order";
import { TableBody, TableCaption, TableCell, TableHead, TableHeader, TableRoot, TableRow } from "~/ui/Table";
import { Seperator } from "~/ui/Seperator";
import { useClient } from "~/providers/client";
import { createForm, required, reset } from "@modular-forms/solid";
import { FieldControl, FieldLabel, FieldMessage, FieldRoot, FormMessage } from "~/ui/Form";
import { Input } from "~/ui/Input";
import { Textarea } from "~/ui/Textarea";
import { DialogCloseButton, DialogContent, DialogHeader, DialogOverlay, DialogPortal, DialogRoot, DialogTitle, DialogTrigger } from "~/ui/Dialog";
import { As } from "@kobalte/core";
import { CheckboxControl, CheckboxLabel, CheckboxRoot } from "~/ui/Checkbox";
import { Skeleton } from "~/ui/Skeleton";

type SearchParams = {
  page: string
  perPage: string
  sort: string
  order: string
}

export function AdminGroups() {
  const [searchParams, setSearchParams] = useSearchParams<SearchParams>()
  const groups = createAsync(() => getListGroups({
    page: {
      page: Number(searchParams.page) || 1,
      perPage: Number(searchParams.perPage) || 10
    },
    sort: searchParams.sort || "",
    order: parseOrder(searchParams.order)
  }))
  const navigate = useNavigate()

  const [createFormOpen, setCreateFormOpen] = createSignal(false);

  const previousDisabled = () => groups()?.pageResult?.previousPage == groups()?.pageResult?.page
  const previous = () => !previousDisabled() && setSearchParams({ page: groups()?.pageResult?.previousPage.toString() } as SearchParams)
  const nextDisabled = () => groups()?.pageResult?.nextPage == groups()?.pageResult?.page
  const next = () => !nextDisabled() && setSearchParams({ page: groups()?.pageResult?.nextPage.toString() } as SearchParams)
  const toggleSort = (value: string) => {
    if (value == groups()?.sort) {
      const order = nextOrder(groups()?.order ?? Order.ORDER_UNSPECIFIED)

      if (order == Order.ORDER_UNSPECIFIED) {
        return setSearchParams({ sort: undefined, order: undefined })
      }

      return setSearchParams({ sort: value, order: encodeOrder(order) } as SearchParams)
    }

    return setSearchParams({ sort: value, order: encodeOrder(Order.DESC) } as SearchParams)
  }

  return (
    <div class="flex justify-center p-4">
      <div class="flex w-full max-w-4xl flex-col gap-2">
        <div class="flex items-center justify-between gap-2">
          <div class="text-xl">Groups</div>
          <DialogRoot open={createFormOpen()} onOpenChange={setCreateFormOpen}>
            <DialogTrigger asChild>
              <As component={Button} size="sm">Create</As>
            </DialogTrigger>
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
        </div>
        <Seperator />
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
                  <SortButton
                    name="name"
                    onClick={toggleSort}
                    sort={groups()?.sort}
                    order={groups()?.order}
                  >
                    Name
                  </SortButton>
                </TableHead>
                <TableHead>
                  <SortButton
                    name="userCount"
                    onClick={toggleSort}
                    sort={groups()?.sort}
                    order={groups()?.order}
                  >
                    Users
                  </SortButton>
                </TableHead>
                <TableHead>
                  <SortButton
                    name="createdAt"
                    onClick={toggleSort}
                    sort={groups()?.sort}
                    order={groups()?.order}
                  >
                    Created At
                  </SortButton>
                </TableHead>
              </tr>
            </TableHeader>
            <TableBody>
              <For each={groups()?.items}>
                {(group) =>
                  <TableRow onClick={() => navigate(`./${group.id}`)} class="cursor-pointer select-none">
                    <TableCell class="p-2">{group.name}</TableCell>
                    <TableCell class="text-nowrap whitespace-nowrap p-2">{group.userCount.toString()}</TableCell>
                    <TableCell class="text-nowrap whitespace-nowrap p-2">{formatDate(parseDate(group.createdAtTime))}</TableCell>
                  </TableRow>
                }
              </For>
            </TableBody>
            <TableCaption>
              <div>
                {groups()?.pageResult?.seenItems.toString()} / {groups()?.pageResult?.totalItems.toString()}
              </div>
              <div>
                Page {groups()?.pageResult?.page}
              </div>
            </TableCaption>
          </TableRoot>
        </Suspense>
      </div>
    </div >
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

type CreateGroupForm = {
  name: string
  description: string
}

const actionCreateGroupForm = action((form: CreateGroupForm) => useClient()
  .admin.createGroup(form)
  .then(() => revalidate(getListGroups.key))
  .catch(throwAsFormError)
)

function CreateGroupForm(props: { setOpen: (value: boolean) => void }) {
  const [createGroupForm, { Field, Form }] = createForm<CreateGroupForm>({ initialValues: { name: "", description: "" } });
  const createGroupFormAction = useAction(actionCreateGroupForm)
  const [keepOpen, setKeepOpen] = createSignal(false)
  const submit = (form: CreateGroupForm) => createGroupFormAction(form)
    .then(() => {
      props.setOpen(keepOpen())
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
      <CheckboxRoot checked={keepOpen()} onChange={setKeepOpen}>
        <CheckboxControl />
        <CheckboxLabel>Keep open</CheckboxLabel>
      </CheckboxRoot>
    </Form>
  )
}

function SortButton(props: ParentProps<{ onClick: (name: string) => void, name: string, sort?: string, order?: Order }>) {
  return <button
    onClick={[props.onClick, props.name]}
    class={cn("text-nowrap flex items-center whitespace-nowrap text-lg", props.name == props.sort && 'text-blue-500')}
  >
    {props.children}
    <Show when={props.sort == props.name && props.order == Order.ASC} fallback={
      <RiArrowsArrowDownSLine class="h-5 w-5" />
    }>
      <RiArrowsArrowUpSLine class="h-5 w-5" />
    </Show>
  </button>
}
