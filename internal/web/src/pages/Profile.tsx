import { action, createAsync, revalidate, useAction, useSubmission } from "@solidjs/router"
import { RiSystemCheckLine, RiSystemCloseLine } from "solid-icons/ri"
import { ComponentProps, ErrorBoundary, For, ParentProps, Show, Suspense, createSignal, resetErrorBoundaries, splitProps } from "solid-js"
import { FormError, createForm, required, reset } from "@modular-forms/solid"

import { formatDate, parseDate, createLoading, catchAsToast, throwAsFormError } from "~/lib/utils"
import { CardContent, CardHeader, CardRoot, CardTitle } from "~/ui/Card"
import { getProfile, getListGroup } from "./Profile.data"
import { AlertDescription, AlertRoot, AlertTitle } from "~/ui/Alert"
import { Button } from "~/ui/Button"
import { TableBody, TableCaption, TableCell, TableHead, TableHeader, TableRoot, TableRow } from "~/ui/Table"
import { useClient } from "~/providers/client"
import { PopoverArrow, PopoverCloseButton, PopoverContent, PopoverPortal, PopoverRoot, PopoverTrigger } from "~/ui/Popover"
import { As } from "@kobalte/core"
import { Badge } from "~/ui/Badge"
import { UserRevokeSessionReq } from "~/twirp/rpc"
import { Seperator } from "~/ui/Seperator"
import { FieldControl, FieldLabel, FieldMessage, FieldRoot, FormMessage } from "~/ui/Form"
import { Input } from "~/ui/Input"
import { Loading } from "~/ui/Loading"
import { Skeleton } from "~/ui/Skeleton"

export const actionRevokeAllSessions = action(() => useClient()
  .user.revokeAllSessions({})
  .then(() => revalidate(getProfile.key)))
export const actionRevokeSession = action((input: UserRevokeSessionReq) => useClient()
  .user.revokeSession(input)
  .then(() => revalidate(getProfile.key)))

export function Profile() {
  const data = createAsync(getProfile)
  const [loading, refreshData] = createLoading(() => revalidate(getProfile.key).then(resetErrorBoundaries))

  const revokeAllSessionsSubmission = useSubmission(actionRevokeAllSessions)
  const revokeAllSessionsAction = useAction(actionRevokeAllSessions)
  const revokeAllSessions = () => revokeAllSessionsAction().catch(catchAsToast)

  return (
    <div class="p-4">
      <ErrorBoundary fallback={(error) => (
        <AlertRoot variant="destructive">
          <AlertTitle>{error.message}</AlertTitle>
          <AlertDescription>
            <Button onClick={refreshData} disabled={loading()}>Retry</Button>
          </AlertDescription>
        </AlertRoot>
      )}>
        <Suspense fallback={<Loading />}>
          <div class="mx-auto flex max-w-4xl flex-col gap-4">
            <CardRoot>
              <CardHeader>
                <CardTitle>Profile</CardTitle>
              </CardHeader>
              <CardContent class="overflow-x-auto">
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
                      <td>{formatDate(parseDate(data()?.createdAt))}</td>
                    </tr>
                    <tr>
                      <td class="pr-2"><Badge class="w-full">Updated At</Badge></td>
                      <td>{formatDate(parseDate(data()?.updatedAt))}</td>
                    </tr>
                  </tbody>
                </table>
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
            <div class="flex">
              <ConfirmButton
                message="Are you sure you wish to revoke all sessions?"
                pending={revokeAllSessionsSubmission.pending}
                onYes={revokeAllSessions}
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
                      const revokeSessionSubmission = useSubmission(actionRevokeSession)
                      const revokeSessionAction = useAction(actionRevokeSession)
                      const revokeSession = (input: UserRevokeSessionReq) => revokeSessionAction(input).catch(catchAsToast)

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
                          <TableCell>{formatDate(parseDate(session.lastUsedAt))}</TableCell>
                          <TableCell>{formatDate(parseDate(session.createdAt))}</TableCell>
                          <TableCell>
                            <Show when={!session.current} fallback={
                              <Badge>Current</Badge>
                            }>
                              <ConfirmButton
                                message="Are you sure you wish to revoke this session?"
                                pending={revokeSessionSubmission.pending}
                                onYes={() => revokeSession({ sessionId: session.id })}
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
            <div class="flex flex-col gap-2">
              <div class="text-xl">Groups</div>
              <Seperator />
            </div>
            <GroupTable />
          </div>
        </Suspense>
      </ErrorBoundary>
    </div>
  )
}


function GroupTable() {
  const data = createAsync(getListGroup)

  return (
    <ErrorBoundary fallback={(error: Error) =>
      <AlertRoot variant="destructive">
        <AlertTitle>{error.message}</AlertTitle>
      </AlertRoot>
    }>
      <Suspense fallback={<Skeleton class="w-full h-32" />}>
        <TableRoot>
          <TableCaption>{data()?.groups.length} Groups(s)</TableCaption>
          <TableHeader>
            <TableRow>
              <TableHead>Name</TableHead>
              <TableHead>Description</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            <For each={data()?.groups}>
              {(group) =>
                <TableRow>
                  <TableCell>{group.name}</TableCell>
                  <TableCell>{group.description}</TableCell>
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

const actionChangeUsername = action((form: ChangeUsernameForm) => useClient()
  .user.updateUsername(form)
  .then(() => revalidate(getProfile.key))
  .catch(throwAsFormError))

function ChangeUsernameForm() {
  const [changeUsernameForm, { Field, Form }] = createForm<ChangeUsernameForm>({ initialValues: { newUsername: "" } });
  const changeUsername = useAction(actionChangeUsername)

  return (
    <Form class="flex w-full max-w-xs flex-col gap-4" onSubmit={(form) => changeUsername(form).then(() => reset(changeUsernameForm))}>
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

const actionChangePassword = action((form: ChangePasswordForm) => {
  if (form.newPassword != form.confirmPassword) {
    throw new FormError<ChangePasswordForm>("", { confirmPassword: "Password does not match." })
  }
  return useClient()
    .user.updatePassword(form)
    .then(() => revalidate(getProfile.key))
    .catch(throwAsFormError)
})

function ChangePasswordForm() {
  const [changePasswordForm, { Field, Form }] = createForm<ChangePasswordForm>({ initialValues: { oldPassword: "", newPassword: "", confirmPassword: "" } });
  const changePassword = useAction(actionChangePassword)

  return (
    <Form class="flex w-full max-w-xs flex-col gap-4" onSubmit={(form) => changePassword(form).then(() => reset(changePasswordForm))}>
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

type ConfirmButtonProps = {
  message: string,
  pending: boolean,
  onYes: () => Promise<unknown>
} & Omit<ComponentProps<typeof Button>, "disabled">

function ConfirmButton(props: ConfirmButtonProps) {
  const [_, rest] = splitProps(props, ["message", "pending", "onYes"])
  const [open, setOpen] = createSignal(false);
  const onYes = () => props.onYes().then(() => setOpen(false))

  return (
    <PopoverRoot open={open()} onOpenChange={setOpen}>
      <PopoverTrigger asChild>
        <As component={Button} disabled={props.pending} {...rest}>
          {props.children}
        </As>
      </PopoverTrigger>
      <PopoverPortal>
        <PopoverContent class="flex flex-col gap-2">
          <PopoverArrow />
          <div>{props.message}</div>
          <div class="flex gap-4">
            <PopoverCloseButton asChild>
              <As component={Button} size="sm" disabled={props.pending}>No</As>
            </PopoverCloseButton>
            <Button
              onClick={onYes}
              disabled={props.pending}
              size="sm"
              variant={props.variant}
            >
              Yes
            </Button>
          </div>
        </PopoverContent>
      </PopoverPortal>
    </PopoverRoot>
  )
}

function Center(props: ParentProps) {
  return (
    <div class="flex justify-center">
      {props.children}
    </div>
  )
}

