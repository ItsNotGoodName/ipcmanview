import { action, createAsync, revalidate, useAction, useSubmission } from "@solidjs/router"
import { RiSystemCheckLine, RiSystemCloseLine } from "solid-icons/ri"
import { ErrorBoundary, For, ParentProps, Show, Suspense, } from "solid-js"
import { FormError, createForm, required, reset } from "@modular-forms/solid"

import { formatDate, parseDate, catchAsToast, throwAsFormError } from "~/lib/utils"
import { CardContent, CardHeader, CardRoot, CardTitle } from "~/ui/Card"
import { getListMyGroups, getProfile, } from "./Profile.data"
import { Button } from "~/ui/Button"
import { TableBody, TableCaption, TableCell, TableHead, TableHeader, TableRoot, TableRow } from "~/ui/Table"
import { useClient } from "~/providers/client"
import { Badge } from "~/ui/Badge"
import { RevokeMySessionReq } from "~/twirp/rpc"
import { Seperator } from "~/ui/Seperator"
import { FieldControl, FieldLabel, FieldMessage, FieldRoot, FormMessage } from "~/ui/Form"
import { Input } from "~/ui/Input"
import { Skeleton } from "~/ui/Skeleton"
import { getSession } from "~/providers/session"
import { PageError } from "~/ui/Page"
import { ConfirmButton } from "~/ui/Confirm"

const actionRevokeAllMySessions = action(() => useClient()
  .user.revokeAllMySessions({})
  .then(() => revalidate(getProfile.key))
  .catch(catchAsToast))

const actionRevokeMySession = action((input: RevokeMySessionReq) => useClient()
  .user.revokeMySession(input)
  .then(() => revalidate(getProfile.key))
  .catch(catchAsToast))

export function Profile() {
  const data = createAsync(getProfile)

  const revokeAllMySessionsSubmission = useSubmission(actionRevokeAllMySessions)
  const revokeAllMySessions = useAction(actionRevokeAllMySessions)

  return (
    <div class="p-4">
      <ErrorBoundary fallback={(e: Error) => <PageError error={e} />}>
        <div class="mx-auto flex max-w-4xl flex-col gap-4">

          <CardRoot>
            <CardHeader>
              <CardTitle>Profile</CardTitle>
            </CardHeader>
            <CardContent class="overflow-x-auto">
              <Suspense fallback={<Skeleton class="h-32" />}>
                <table>
                  <tbody>
                    <tr>
                      <td class="pr-2"><Badge class="flex w-full justify-center">Username</Badge></td>
                      <td>{data()?.username}</td>
                    </tr>
                    <tr>
                      <td class="pr-2"><Badge class="flex w-full justify-center">Email</Badge></td>
                      <td>{data()?.email}</td>
                    </tr>
                    <tr>
                      <td class="pr-2"><Badge class="flex w-full justify-center">Admin</Badge></td>
                      <td>
                        <Show when={data()?.admin} fallback={<RiSystemCloseLine class="h-6 w-6 text-red-500" />}>
                          <RiSystemCheckLine class="h-6 w-6 text-green-500" />
                        </Show>
                      </td>
                    </tr>
                    <tr>
                      <td class="pr-2"><Badge class="flex w-full justify-center">Created At</Badge></td>
                      <td>{formatDate(parseDate(data()?.createdAtTime))}</td>
                    </tr>
                    <tr>
                      <td class="pr-2"><Badge class="w-full">Updated At</Badge></td>
                      <td>{formatDate(parseDate(data()?.updatedAtTime))}</td>
                    </tr>
                  </tbody>
                </table>
              </Suspense>
            </CardContent>
          </CardRoot>

          <div class="flex flex-col gap-2">
            <div class="text-xl">Change username</div>
            <Seperator />
          </div>
          <Center>
            <ChangeUsernameForm />
          </Center>

          <div class="flex flex-col gap-2">
            <div class="text-xl">Change password</div>
            <Seperator />
          </div>
          <Center>
            <ChangePasswordForm />
          </Center>

          <div class="flex flex-col gap-2">
            <div class="text-xl">Sessions</div>
            <Seperator />
          </div>
          <Suspense fallback={<Skeleton class="h-32" />}>
            <div class="flex">
              <ConfirmButton
                message="Are you sure you wish to revoke all sessions?"
                pending={revokeAllMySessionsSubmission.pending}
                onYes={revokeAllMySessions}
                variant="destructive"
              >
                Revoke all sessions
              </ConfirmButton>
            </div>
            <TableRoot>
              <TableCaption>{data()?.sessions.length} Session(s)</TableCaption>
              <TableHeader>
                <TableRow>
                  <TableHead>Active</TableHead>
                  <TableHead>User Agent</TableHead>
                  <TableHead>IP</TableHead>
                  <TableHead>Last IP</TableHead>
                  <TableHead>Last Used At</TableHead>
                  <TableHead>Created At</TableHead>
                  <TableHead></TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                <For each={data()?.sessions}>
                  {
                    (session) => {
                      const revokeMySessionSubmission = useSubmission(actionRevokeMySession)
                      const revokeMySession = useAction(actionRevokeMySession)

                      return (
                        <TableRow>
                          <TableCell>
                            <Show when={session.active} fallback={<div class="mx-auto h-4 w-4 rounded-full bg-gray-500" title="Inactive" />}>
                              <div class="mx-auto h-4 w-4 rounded-full bg-green-500" title="Active" />
                            </Show>
                          </TableCell>
                          <TableCell>{session.userAgent}</TableCell>
                          <TableCell>{session.ip}</TableCell>
                          <TableCell>{session.lastIp}</TableCell>
                          <TableCell>{formatDate(parseDate(session.lastUsedAtTime))}</TableCell>
                          <TableCell>{formatDate(parseDate(session.createdAtTime))}</TableCell>
                          <TableCell>
                            <Show when={!session.current} fallback={
                              <Badge>Current</Badge>
                            }>
                              <ConfirmButton
                                message="Are you sure you wish to revoke this session?"
                                pending={revokeMySessionSubmission.pending}
                                onYes={() => revokeMySession({ sessionId: session.id })}
                                variant="destructive"
                                size="sm"
                              >
                                Revoke
                              </ConfirmButton>
                            </Show>
                          </TableCell>
                        </TableRow>
                      )
                    }
                  }
                </For>
              </TableBody>
            </TableRoot>
          </Suspense>

          <div class="flex flex-col gap-2">
            <div class="text-xl">Groups</div>
            <Seperator />
          </div>
          <GroupTable />

        </div>
      </ErrorBoundary>
    </div>
  )
}

function GroupTable() {
  const data = createAsync(getListMyGroups)

  return (
    <ErrorBoundary fallback={(e: Error) => <PageError error={e} />}>
      <Suspense fallback={<Skeleton class="h-32" />}>
        <TableRoot>
          <TableCaption>{data()?.groups.length} Groups(s)</TableCaption>
          <TableHeader>
            <TableRow>
              <TableHead>Name</TableHead>
              <TableHead>Description</TableHead>
              <TableHead>Joined At</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            <For each={data()?.groups}>
              {(group) =>
                <TableRow>
                  <TableCell>{group.name}</TableCell>
                  <TableCell>{group.description}</TableCell>
                  <TableCell>{formatDate(parseDate(group.joinedAtTime))}</TableCell>
                </TableRow>
              }
            </For>
          </TableBody>
        </TableRoot>
      </Suspense>
    </ErrorBoundary>
  )
}

type ChangeUsernameForm = {
  newUsername: string
}

const actionUpdateMyUsername = action((form: ChangeUsernameForm) => useClient()
  .user.updateMyUsername(form)
  .then(() => revalidate([getProfile.key, getSession.key]))
  .catch(throwAsFormError))

function ChangeUsernameForm() {
  const [changeUsernameForm, { Field, Form }] = createForm<ChangeUsernameForm>({ initialValues: { newUsername: "" } });
  const submit = useAction(actionUpdateMyUsername)

  return (
    <Form class="flex w-full max-w-xs flex-col gap-4" onSubmit={(form) => submit(form).then(() => reset(changeUsernameForm))}>
      <Field name="newUsername" validate={required("Please enter a new username.")}>
        {(field, props) => (
          <FieldRoot class="gap-1.5">
            <FieldLabel field={field}>New username</FieldLabel>
            <FieldControl field={field}>
              <Input
                {...props}
                placeholder="New username"
                value={field.value}
              />
            </FieldControl>
            <FieldMessage field={field} />
          </FieldRoot>
        )}
      </Field>
      <Button type="submit" disabled={changeUsernameForm.submitting}>
        <Show when={changeUsernameForm.submitting} fallback={<>Update username</>}>
          Updating username
        </Show>
      </Button>
      <FormMessage form={changeUsernameForm} />
    </Form>
  )
}

type ChangePasswordForm = {
  oldPassword: string
  newPassword: string
  confirmPassword: string
}

const actionUpdateMyPassword = action((form: ChangePasswordForm) => {
  if (form.newPassword != form.confirmPassword) {
    throw new FormError<ChangePasswordForm>("", { confirmPassword: "Password does not match." })
  }
  return useClient()
    .user.updateMyPassword(form)
    .then(() => revalidate(getProfile.key))
    .catch(throwAsFormError)
})

function ChangePasswordForm() {
  const [changePasswordForm, { Field, Form }] = createForm<ChangePasswordForm>({ initialValues: { oldPassword: "", newPassword: "", confirmPassword: "" } });
  const submit = useAction(actionUpdateMyPassword)

  return (
    <Form class="flex w-full max-w-xs flex-col gap-4" onSubmit={(form) => submit(form).then(() => reset(changePasswordForm))}>
      <input class="hidden" type="text" name="username" autocomplete="username" />
      <Field name="oldPassword" validate={required("Please enter your old password.")}>
        {(field, props) => (
          <FieldRoot class="gap-1.5">
            <FieldLabel field={field}>Old password</FieldLabel>
            <FieldControl field={field}>
              <Input
                {...props}
                autocomplete="current-password"
                placeholder="Old password"
                type="password"
                value={field.value}
              />
            </FieldControl>
            <FieldMessage field={field} />
          </FieldRoot>
        )}
      </Field>
      <Field name="newPassword" validate={required("Please enter a new password.")}>
        {(field, props) => (
          <FieldRoot class="gap-1.5">
            <FieldLabel field={field}>New password</FieldLabel>
            <FieldControl field={field}>
              <Input
                {...props}
                autocomplete="new-password"
                placeholder="New password"
                type="password"
                value={field.value}
              />
            </FieldControl>
            <FieldMessage field={field} />
          </FieldRoot>
        )}
      </Field>
      <Field name="confirmPassword">
        {(field, props) => (
          <FieldRoot class="gap-1.5">
            <FieldLabel field={field}>Confirm new password</FieldLabel>
            <FieldControl field={field}>
              <Input
                {...props}
                autocomplete="new-password"
                placeholder="Confirm new password"
                type="password"
                value={field.value}
              />
            </FieldControl>
            <FieldMessage field={field} />
          </FieldRoot>
        )}
      </Field>
      <Button type="submit" disabled={changePasswordForm.submitting}>
        <Show when={changePasswordForm.submitting} fallback={<>Update password</>}>
          Updating password
        </Show>
      </Button>
      <FormMessage form={changePasswordForm} />
    </Form>
  )
}

function Center(props: ParentProps) {
  return (
    <div class="flex justify-center">
      {props.children}
    </div>
  )
}

