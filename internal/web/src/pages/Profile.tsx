import Humanize from "humanize-plus"
import { action, createAsync, revalidate, useAction, useSubmission } from "@solidjs/router"
import { RiSystemCheckLine, RiSystemCloseLine } from "solid-icons/ri"
import { ErrorBoundary, For, ParentProps, Show, Suspense, createSignal, } from "solid-js"
import { createForm, required, reset } from "@modular-forms/solid"

import { formatDate, parseDate, catchAsToast, throwAsFormError, createModal } from "~/lib/utils"
import { CardRoot, } from "~/ui/Card"
import { getProfilePage } from "./Profile.data"
import { Button } from "~/ui/Button"
import { TableBody, TableCaption, TableCell, TableHead, TableHeader, TableRoot, TableRow } from "~/ui/Table"
import { useClient } from "~/providers/client"
import { Badge } from "~/ui/Badge"
import { FieldLabel, FieldMessage, FieldRoot, FormMessage, fieldControlProps } from "~/ui/Form"
import { Input } from "~/ui/Input"
import { Skeleton } from "~/ui/Skeleton"
import { getSession } from "~/providers/session"
import { PageError } from "~/ui/Page"
import { LayoutNormal } from "~/ui/Layout"
import { AlertDialogAction, AlertDialogCancel, AlertDialogModal, AlertDialogFooter, AlertDialogHeader, AlertDialogRoot, AlertDialogTitle } from "~/ui/AlertDialog"
import { Shared } from "~/components/Shared"

function Center(props: ParentProps) {
  return (
    <div class="flex justify-center">
      {props.children}
    </div>
  )
}

const actionRevokeAllMySessions = action(() => useClient()
  .user.revokeAllMySessions({})
  .then(() => revalidate(getProfilePage.key))
  .catch(catchAsToast))

const actionRevokeMySession = action((sessionId: bigint) => useClient()
  .user.revokeMySession({ sessionId })
  .then(() => revalidate(getProfilePage.key))
  .catch(catchAsToast))

export function Profile() {
  const data = createAsync(() => getProfilePage())

  const [revokeAllMySessionsModal, setRevokeAllMySessionsModal] = createSignal(false)
  const revokeAllMySessionsSubmission = useSubmission(actionRevokeAllMySessions)
  const revokeAllMySessionsAction = useAction(actionRevokeAllMySessions)
  const submitRevokeAllMySessions = () => revokeAllMySessionsAction()
    .then(() => setRevokeAllMySessionsModal(false))

  const revokeMySessionModal = createModal(BigInt(0))
  const revokeMySessionSubmission = useSubmission(actionRevokeMySession)
  const revokeMySessionAction = useAction(actionRevokeMySession)
  const submitRevokeMySession = () => revokeMySessionAction(revokeMySessionModal.value())
    .then(revokeMySessionModal.setClose)

  return (
    <LayoutNormal class="max-w-4xl">
      <ErrorBoundary fallback={(e) => <PageError error={e} />}>
        <AlertDialogRoot open={revokeAllMySessionsModal()} onOpenChange={setRevokeAllMySessionsModal}>
          <AlertDialogModal>
            <AlertDialogHeader>
              <AlertDialogTitle>Are you sure you wish to revoke all sessions?</AlertDialogTitle>
            </AlertDialogHeader>
            <AlertDialogFooter>
              <AlertDialogCancel>Cancel</AlertDialogCancel>
              <AlertDialogAction disabled={revokeAllMySessionsSubmission.pending} onClick={submitRevokeAllMySessions} variant="destructive">
                Delete
              </AlertDialogAction>
            </AlertDialogFooter>
          </AlertDialogModal>
        </AlertDialogRoot>

        <AlertDialogRoot open={revokeMySessionModal.open()} onOpenChange={revokeMySessionModal.setClose}>
          <AlertDialogModal>
            <AlertDialogHeader>
              <AlertDialogTitle>Are you sure you wish to revoke this session?</AlertDialogTitle>
            </AlertDialogHeader>
            <AlertDialogFooter>
              <AlertDialogCancel>Cancel</AlertDialogCancel>
              <AlertDialogAction disabled={revokeMySessionSubmission.pending} onClick={submitRevokeMySession} variant="destructive">
                Delete
              </AlertDialogAction>
            </AlertDialogFooter>
          </AlertDialogModal>
        </AlertDialogRoot>

        <Shared.Title>Profile</Shared.Title>

        <CardRoot class="overflow-x-auto p-4">
          <Suspense fallback={<Skeleton class="h-32" />}>
            <table>
              <tbody>
                <tr>
                  <th class="pr-2">Username</th>
                  <td>{data()?.username}</td>
                </tr>
                <tr>
                  <th class="pr-2">Email</th>
                  <td>{data()?.email}</td>
                </tr>
                <tr>
                  <th class="pr-2">Admin</th>
                  <td>
                    <Show when={data()?.admin} fallback={<RiSystemCloseLine class="h-6 w-6 text-red-500" />}>
                      <RiSystemCheckLine class="h-6 w-6 text-green-500" />
                    </Show>
                  </td>
                </tr>
                <tr>
                  <th class="pr-2">Created At</th>
                  <td>{formatDate(parseDate(data()?.createdAtTime))}</td>
                </tr>
                <tr>
                  <th class="pr-2">Updated At</th>
                  <td>{formatDate(parseDate(data()?.updatedAtTime))}</td>
                </tr>
              </tbody>
            </table>
          </Suspense>
        </CardRoot>

        <Shared.Title>Change username</Shared.Title>
        <ChangeUsernameForm />

        <Shared.Title>Change password</Shared.Title>
        <ChangePasswordForm />

        <Shared.Title>Sessions</Shared.Title>
        <div class="flex">
          <Button onClick={() => setRevokeAllMySessionsModal(true)} variant="destructive">
            Revoke all sessions
          </Button>
        </div>
        <Suspense fallback={<Skeleton class="h-32" />}>
          <TableRoot>
            <TableCaption>{data()?.sessions.length} {Humanize.pluralize(data()?.sessions.length || 0, "Session")}</TableCaption>
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
                {(session) => (
                  <TableRow>
                    <TableCell>
                      <Show when={session.active} fallback={<div title="Inactive" class="mx-auto h-4 w-4 rounded-full bg-gray-500" />}>
                        <div title="Active" class="mx-auto h-4 w-4 rounded-full bg-green-500" />
                      </Show>
                    </TableCell>
                    <TableCell>{session.userAgent}</TableCell>
                    <TableCell>{session.ip}</TableCell>
                    <TableCell>{session.lastIp}</TableCell>
                    <TableCell>{formatDate(parseDate(session.lastUsedAtTime))}</TableCell>
                    <TableCell>{formatDate(parseDate(session.createdAtTime))}</TableCell>
                    <TableCell class="py-0">
                      <Show when={!session.current} fallback={<Badge>Current</Badge>}>
                        <Button onClick={() => revokeMySessionModal.setValue(session.id)} variant="destructive" size="sm">
                          Revoke
                        </Button>
                      </Show>
                    </TableCell>
                  </TableRow>
                )}
              </For>
            </TableBody>
          </TableRoot>
        </Suspense>

        <Shared.Title>Groups</Shared.Title>
        <Suspense fallback={<Skeleton class="h-32" />}>
          <TableRoot>
            <TableCaption>{data()?.groups.length} {Humanize.pluralize(data()?.groups.length || 0, "Group")}</TableCaption>
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
    </LayoutNormal>
  )
}

type ChangeUsernameForm = {
  newUsername: string
}

function ChangeUsernameForm() {
  const [form, { Field, Form }] = createForm<ChangeUsernameForm>({
    initialValues: {
      newUsername: ""
    }
  });
  const submit = (input: ChangeUsernameForm) => useClient()
    .user.updateMyUsername(input)
    .then(() => revalidate([getProfilePage.key, getSession.key]))
    .then(() => reset(form))
    .catch(throwAsFormError)

  return (
    <Center>
      <Form onSubmit={submit} class="flex w-full max-w-sm flex-col gap-4">
        <Field name="newUsername" validate={required("Please enter a new username.")}>
          {(field, props) => (
            <FieldRoot>
              <FieldLabel field={field}>New username</FieldLabel>
              <Input
                {...props}
                {...fieldControlProps(field)}
                placeholder="New username"
                value={field.value}
              />
              <FieldMessage field={field} />
            </FieldRoot>
          )}
        </Field>
        <Button type="submit" disabled={form.submitting}>
          <Show when={!form.submitting} fallback="Updating username">Update username</Show>
        </Button>
        <FormMessage form={form} />
      </Form>
    </Center>
  )
}

type ChangePasswordForm = {
  oldPassword: string
  newPassword: string
  confirmPassword: string
}

function ChangePasswordForm() {
  const [form, { Field, Form }] = createForm<ChangePasswordForm>({
    initialValues: {
      oldPassword: "",
      newPassword: "",
      confirmPassword: "",
    },
    validate: (input) => {
      if (input.newPassword != input.confirmPassword) {
        return { confirmPassword: "Password does not match." }
      }
      return {}
    }
  });
  const submit = (input: ChangePasswordForm) => useClient()
    .user.updateMyPassword(input)
    .then(() => revalidate(getProfilePage.key))
    .then(() => reset(form))
    .catch(throwAsFormError)

  return (
    <Center>
      <Form onSubmit={submit} class="flex w-full max-w-sm flex-col gap-4">
        <input class="hidden" type="text" name="username" autocomplete="username" />
        <Field name="oldPassword" validate={required("Please enter your old password.")}>
          {(field, props) => (
            <FieldRoot>
              <FieldLabel field={field}>Old password</FieldLabel>
              <Input
                {...props}
                {...fieldControlProps(field)}
                autocomplete="current-password"
                placeholder="Old password"
                type="password"
                value={field.value}
              />
              <FieldMessage field={field} />
            </FieldRoot>
          )}
        </Field>
        <Field name="newPassword" validate={required("Please enter a new password.")}>
          {(field, props) => (
            <FieldRoot>
              <FieldLabel field={field}>New password</FieldLabel>
              <Input
                {...props}
                {...fieldControlProps(field)}
                autocomplete="new-password"
                placeholder="New password"
                type="password"
                value={field.value}
              />
              <FieldMessage field={field} />
            </FieldRoot>
          )}
        </Field>
        <Field name="confirmPassword">
          {(field, props) => (
            <FieldRoot>
              <FieldLabel field={field}>Confirm new password</FieldLabel>
              <Input
                {...props}
                {...fieldControlProps(field)}
                autocomplete="new-password"
                placeholder="Confirm new password"
                type="password"
                value={field.value}
              />
              <FieldMessage field={field} />
            </FieldRoot>
          )}
        </Field>
        <Button type="submit" disabled={form.submitting}>
          <Show when={form.submitting} fallback="Update password">Updating password</Show>
        </Button>
        <FormMessage form={form} />
      </Form>
    </Center>
  )
}

export default Profile
