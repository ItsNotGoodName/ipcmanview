import { action, createAsync, revalidate, useAction } from "@solidjs/router";
import { CheckboxControl, CheckboxInput, CheckboxRoot } from "~/ui/Checkbox";
import { TableBody, TableCaption, TableCell, TableHead, TableHeader, TableRoot, TableRow } from "~/ui/Table";
import { getListGroups } from "./Home.data";
import { FieldControl, FieldLabel, FieldMessage, FieldRoot, FormMessage } from "~/ui/Form"
import { ListGroupsReq } from "~/twirp/rpc";
import { ErrorBoundary, For, Show, Suspense } from "solid-js";
import { createLoading, formatDate, parseDate, throwAsFormError } from "~/lib/utils";
import { Seperator } from "~/ui/Seperator";
import { DialogCloseButton, DialogContent, DialogHeader, DialogOverlay, DialogPortal, DialogRoot, DialogTitle, DialogTrigger } from "~/ui/Dialog";
import { As } from "@kobalte/core";
import { Button } from "~/ui/Button";
import { createForm, required, reset } from "@modular-forms/solid";
import { Input } from "~/ui/Input";
import { Textarea } from "~/ui/Textarea";
import { useClient } from "~/providers/client";
import { ConfirmButton } from "~/ui/Confirm";
import { PageError } from "~/ui/Page";
import { Skeleton } from "~/ui/Skeleton";

export function AdminHome() {
  const input: ListGroupsReq = {
    page: {
      page: 1,
      perPage: 100,
    }
  }
  const data = createAsync(() => getListGroups(input))
  const [dataRefreshing, refreshData] = createLoading(() => revalidate(getListGroups.key))

  return (
    <div class="mx-auto flex max-w-4xl flex-col p-4">
      <div class="flex flex-col gap-2">
        <ErrorBoundary fallback={(e: Error) => <PageError error={e} />}>
          <Suspense fallback={<Skeleton class="h32" />}>
            <div class="flex flex-col gap-2">
              <div class="text-xl">Groups</div>
              <Seperator />
            </div>
            <div class="flex flex-wrap gap-2">
              <DialogRoot>
                <DialogTrigger asChild>
                  <As component={Button}>Create</As>
                </DialogTrigger>
                <DialogPortal>
                  <DialogOverlay />
                  <DialogContent>
                    <DialogHeader>
                      <DialogCloseButton />
                      <DialogTitle>Create group</DialogTitle>
                    </DialogHeader>
                    <CreateGroupForm />
                  </DialogContent>
                </DialogPortal>
              </DialogRoot>
              <DialogRoot>
                <DialogTrigger asChild>
                  <As component={Button} variant="secondary">Update</As>
                </DialogTrigger>
                <DialogPortal>
                  <DialogOverlay />
                  <DialogContent>
                    <DialogHeader>
                      <DialogCloseButton />
                      <DialogTitle>Update group</DialogTitle>
                    </DialogHeader>
                    <UpdateGroupForm />
                  </DialogContent>
                </DialogPortal>
              </DialogRoot>
              <ConfirmButton variant="destructive" message="Are you sure you wish to delete 0 group?">
                Delete
              </ConfirmButton>
              <Button variant="outline" disabled={dataRefreshing()} onClick={refreshData}>
                Refresh
              </Button>
            </div>
            <form>
              <TableRoot>
                <TableCaption>{data()?.groups.length} / {data()?.pageResult?.totalItems.toString()} Groups</TableCaption>
                <TableHeader>
                  <TableRow>
                    <TableHead>
                      <CheckboxRoot>
                        <CheckboxControl />
                      </CheckboxRoot>
                    </TableHead>
                    <TableHead>Name</TableHead>
                    <TableHead>User Count</TableHead>
                    <TableHead>Created At</TableHead>
                    <TableHead>Updated At</TableHead>
                  </TableRow>
                </TableHeader>
                <TableBody>
                  <For each={data()?.groups}>{(group) =>
                    <TableRow>
                      <TableCell>
                        <input type="hidden" name={`group[${group.id}]`} />
                        <CheckboxRoot name={`group[${group.id}].selected`}>
                          <CheckboxInput />
                          <CheckboxControl />
                        </CheckboxRoot>
                      </TableCell>
                      <TableCell>{group.name}</TableCell>
                      <TableCell>{group.userCount.toString()}</TableCell>
                      <TableCell>{formatDate(parseDate(group.createdAtTime))}</TableCell>
                      <TableCell>{formatDate(parseDate(group.updatedAtTime))}</TableCell>
                    </TableRow>
                  }
                  </For>
                </TableBody>
              </TableRoot>
            </form>
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
  .then(() => revalidate(getListGroups.key))
  .catch(throwAsFormError)
)

function CreateGroupForm() {
  const [createGroupForm, { Field, Form }] = createForm<CreateGroupForm>({ initialValues: { name: "", description: "" } });
  const submit = useAction(actionCreateGroupForm)

  return (
    <Form class="flex flex-col gap-4" onSubmit={(form) => submit(form).then(() => reset(createGroupForm))}>
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
                placeholder="Description"
              >{field.value}</Textarea>
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
    </Form>
  )
}

type UpdateGroupForm = {
  name: string
  description: string
}

const actionUpdateGroupForm = action((form: UpdateGroupForm) => useClient()
  .admin.updateGroup(form)
  .then(() => revalidate(getListGroups.key))
  .catch(throwAsFormError)
)

function UpdateGroupForm() {
  const [updateGroupForm, { Field, Form }] = createForm<UpdateGroupForm>({ initialValues: { name: "", description: "" } });
  const submit = useAction(actionUpdateGroupForm)

  return (
    <Form class="flex flex-col gap-4" onSubmit={(form) => submit(form).then(() => reset(updateGroupForm))}>
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
                placeholder="Description"
              >{field.value}</Textarea>
            </FieldControl>
            <FieldMessage field={field} />
          </FieldRoot>
        )}
      </Field>
      <Button type="submit" disabled={updateGroupForm.submitting}>
        <Show when={updateGroupForm.submitting} fallback={<>Update group</>}>
          Updating group
        </Show>
      </Button>
      <FormMessage form={updateGroupForm} />
    </Form>
  )

}
