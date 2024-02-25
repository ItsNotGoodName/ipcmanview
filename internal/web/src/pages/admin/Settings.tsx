import { createForm, reset } from "@modular-forms/solid";
import { ErrorBoundary, For, Show, Suspense, batch, createSignal, } from "solid-js";
import { Shared } from "~/components/Shared";
import { Button } from "~/ui/Button";
import { FieldDescription, FieldLabel, FieldMessage, FieldRoot, SwitchFieldRoot, FormMessage, fieldControlProps, CheckboxFieldRoot } from "~/ui/Form";
import { Input } from "~/ui/Input";
import { LayoutNormal } from "~/ui/Layout";
import { SwitchControl, SwitchDescription, SwitchErrorMessage, SwitchLabel } from "~/ui/Switch";
import { useClient } from "~/providers/client";
import { action, createAsync, revalidate, useAction, useSubmission } from "@solidjs/router";
import { getConfig } from "../data";
import { GetConfigResp, ListEventRulesResp_Item } from "~/twirp/rpc";
import { Skeleton } from "~/ui/Skeleton";
import { catchAsToast, createModal, createRowSelection, throwAsFormError } from "~/lib/utils";
import { PageError } from "~/ui/Page";
import { TableBody, TableCell, TableHead, TableHeader, TableRoot, TableRow } from "~/ui/Table";
import { CheckboxControl, CheckboxErrorMessage, CheckboxLabel, CheckboxRoot } from "~/ui/Checkbox";
import { getListEventRules } from "./data";
import { DialogContent, DialogHeader, DialogOverflow, DialogOverlay, DialogPortal, DialogRoot, DialogTitle } from "~/ui/Dialog";
import { RiDeviceSaveLine, RiSystemAddLine, RiSystemDeleteBinLine, RiSystemRefreshLine } from "solid-icons/ri";
import { Crud } from "~/components/Crud";
import { AlertDialogAction, AlertDialogCancel, AlertDialogDescription, AlertDialogFooter, AlertDialogHeader, AlertDialogModal, AlertDialogRoot, AlertDialogTitle } from "~/ui/AlertDialog";

type UpdateForm = {
  siteName: string
  enableSignUp: {
    boolean: boolean
  }
}

export function AdminSettings() {
  const config = createAsync(() => getConfig())
  const refetchConfig = () => revalidate(getConfig.key)

  const eventRules = createAsync(() => getListEventRules())
  const refetchEventRules = () => revalidate(getListEventRules.key)

  return (
    <LayoutNormal class="max-w-4xl">
      <Shared.Title>Settings</Shared.Title>
      <ErrorBoundary fallback={(e) => <PageError error={e} />}>
        <Suspense fallback={<Skeleton class="h-32" />}>
          <Show when={config()}>
            <UpdateSettingsForm config={config()!} refetchConfig={refetchConfig} />
          </Show>
        </Suspense>

        <Shared.Title>Event rules</Shared.Title>
        <Suspense fallback={<Skeleton class="h-32" />}>
          <Show when={eventRules()}>
            <Thing eventRules={eventRules()!} refetchEventRules={refetchEventRules} />
          </Show>
        </Suspense>
      </ErrorBoundary>
    </LayoutNormal >
  )
}

function UpdateSettingsForm(props: { config: GetConfigResp, refetchConfig: () => Promise<void> }) {
  const formInitialValues = (): UpdateForm => ({
    siteName: props.config.siteName,
    enableSignUp: {
      boolean: props.config.enableSignUp
    }
  })
  const [updateForm, { Field, Form }] = createForm<UpdateForm>({ initialValues: formInitialValues() });
  const formReset = () => props.refetchConfig().then(() => reset(updateForm, { initialValues: formInitialValues() }))
  const formSubmit = (form: UpdateForm) => useClient()
    .admin.updateConfig({
      enableSignUp: form.enableSignUp.boolean,
      siteName: form.siteName,
    })
    .then(formReset)
    .catch(throwAsFormError)

  return (
    <div class="flex justify-center">
      <Form class="flex w-full max-w-sm flex-col gap-4" onSubmit={formSubmit}>
        <Field name="siteName">
          {(field, props) => (
            <FieldRoot>
              <FieldLabel field={field}>Site name</FieldLabel>
              <Input
                {...props}
                {...fieldControlProps(field)}
                value={field.value}
              />
              <FieldDescription>Name of site.</FieldDescription>
              <FieldMessage field={field} />
            </FieldRoot>
          )}
        </Field>
        <Field name="enableSignUp.boolean" type="boolean">
          {(field, props) => (
            <SwitchFieldRoot
              form={updateForm}
              field={field}
              class="flex items-center justify-between gap-2"
            >
              <div>
                <SwitchLabel>Enable sign up</SwitchLabel>
                <SwitchDescription>Allow public sign up.</SwitchDescription>
                <SwitchErrorMessage>{field.error}</SwitchErrorMessage>
              </div>
              <SwitchControl inputProps={props} />
            </SwitchFieldRoot>
          )}
        </Field>
        <div class="flex flex-col gap-4 sm:flex-row-reverse">
          <Button type="submit" disabled={updateForm.submitting} class="sm:flex-1">
            <Show when={!updateForm.submitting} fallback="Updating settings">Update settings</Show>
          </Button>
          <Button type="button" onClick={formReset} variant="destructive" disabled={updateForm.submitting}>Reset</Button>
        </div>
        <FormMessage form={updateForm} />
      </Form>
    </div>
  )
}

const actionDeleteEventRule = action((ids: bigint[]) => useClient()
  .admin.deleteEventRules({ ids })
  .then(() => revalidate(getListEventRules.key))
  .catch(catchAsToast))

function Thing(props: { eventRules: ListEventRulesResp_Item[], refetchEventRules: () => void }) {
  const rowSelection = createRowSelection(() => props.eventRules?.map(v => v.id) || [])

  // Create
  const [createFormModal, setCreateFormModal] = createSignal(false)

  // Delete
  const deleteSubmission = useSubmission(actionDeleteEventRule)
  const deleteAction = useAction(actionDeleteEventRule)
  // Single
  const deleteModal = createModal({ name: "", id: BigInt(0) })
  const deleteSubmit = () =>
    deleteAction([deleteModal.value().id])
      .then(deleteModal.close)
  // Multiple
  const [deleteMultipleModal, setDeleteMultipleModal] = createSignal(false)
  const deleteMultipleSubmit = () =>
    deleteAction(rowSelection.selections())
      .then(() => batch(() => {
        rowSelection.setAll(false)
        setDeleteMultipleModal(false)
      }))

  return (
    <div class="flex flex-col gap-2">
      <DialogRoot open={createFormModal()} onOpenChange={setCreateFormModal}>
        <DialogPortal>
          <DialogOverlay />
          <DialogContent>
            <DialogHeader>
              <DialogTitle>Create event rule</DialogTitle>
            </DialogHeader>
            <DialogOverflow>
              <CreateEventRuleForm onSubmit={() => setCreateFormModal(false)} />
            </DialogOverflow>
          </DialogContent>
        </DialogPortal>
      </DialogRoot>

      <AlertDialogRoot open={deleteModal.open()} onOpenChange={deleteModal.close}>
        <AlertDialogModal>
          <AlertDialogHeader>
            <AlertDialogTitle>Are you sure you wish to delete {deleteModal.value().name}?</AlertDialogTitle>
          </AlertDialogHeader>
          <AlertDialogFooter>
            <AlertDialogCancel>Cancel</AlertDialogCancel>
            <AlertDialogAction variant="destructive" disabled={deleteSubmission.pending} onClick={deleteSubmit}>
              Delete
            </AlertDialogAction>
          </AlertDialogFooter>
        </AlertDialogModal>
      </AlertDialogRoot>

      <AlertDialogRoot open={deleteMultipleModal()} onOpenChange={setDeleteMultipleModal}>
        <AlertDialogModal>
          <AlertDialogHeader>
            <AlertDialogTitle>Are you sure you wish to delete {rowSelection.selections().length} devices?</AlertDialogTitle>
            <AlertDialogDescription>
              <ul>
                <For each={props.eventRules}>
                  {(e, index) =>
                    <Show when={rowSelection.rows[index()]?.checked}>
                      <li>{e.code}</li>
                    </Show>
                  }
                </For>
              </ul>
            </AlertDialogDescription>
          </AlertDialogHeader>
          <AlertDialogFooter>
            <AlertDialogCancel>Cancel</AlertDialogCancel>
            <AlertDialogAction variant="destructive" disabled={deleteSubmission.pending} onClick={deleteMultipleSubmit}>
              Delete
            </AlertDialogAction>
          </AlertDialogFooter>
        </AlertDialogModal>
      </AlertDialogRoot>

      <div class="flex justify-end gap-2">
        <Button size="icon" title="Create" variant="secondary" onClick={() => setCreateFormModal(true)}>
          <RiSystemAddLine class="size-5" />
        </Button>
        <Button size="icon" title="Update">
          <RiDeviceSaveLine class="size-5" />
        </Button>
        <Button size="icon" title="Delete" variant="destructive"
          disabled={rowSelection.selections().length == 0 || deleteSubmission.pending}
          onClick={() => setDeleteMultipleModal(true)}
        >
          <RiSystemDeleteBinLine class="size-5" />
        </Button>
      </div>

      <TableRoot>
        <TableHeader>
          <TableRow>
            <TableHead>
              <CheckboxRoot
                indeterminate={rowSelection.indeterminate()}
                checked={rowSelection.multiple()}
                onChange={(checked) => rowSelection.setAll(checked)}
              >
                <CheckboxControl />
              </CheckboxRoot>
            </TableHead>
            <TableHead>Code</TableHead>
            <TableHead>DB</TableHead>
            <TableHead>Live</TableHead>
            <TableHead>MQTT</TableHead>
            <Crud.LastTableHead>
              <Button size="icon" title="Refresh" variant="ghost" onClick={props.refetchEventRules}>
                <RiSystemRefreshLine class="size-5" />
              </Button>
            </Crud.LastTableHead>
          </TableRow>
        </TableHeader>
        <TableBody>
          <For each={props.eventRules}>
            {(item, index) =>
              <TableRow>
                <TableCell>
                  <Show when={item.code} fallback={
                    <CheckboxRoot checked={false} disabled>
                      <CheckboxControl />
                    </CheckboxRoot>
                  }>
                    <CheckboxRoot
                      checked={rowSelection.rows[index()]?.checked}
                      onChange={(checked) => rowSelection.set(item.id, checked)}
                    >
                      <CheckboxControl />
                    </CheckboxRoot>
                  </Show>
                </TableCell>
                <TableCell class="min-w-32 w-full py-0">
                  <Show when={item.code}>
                    <Input value={item.code} />
                  </Show>
                </TableCell>
                <TableCell>
                  <CheckboxRoot
                    checked={!item.ignoreDb}
                  >
                    <CheckboxControl />
                  </CheckboxRoot>
                </TableCell>
                <TableCell>
                  <CheckboxRoot
                    checked={!item.ignoreLive}
                  >
                    <CheckboxControl />
                  </CheckboxRoot>
                </TableCell>
                <TableCell>
                  <CheckboxRoot
                    checked={!item.ignoreMqtt}
                  >
                    <CheckboxControl />
                  </CheckboxRoot>
                </TableCell>
                <TableCell />
              </TableRow>
            }
          </For>
        </TableBody>
      </TableRoot>
    </div>
  )
}


type CreateEventRuleForm = {
  code: string
  db: {
    boolean: boolean
  }
  live: {
    boolean: boolean
  }
  mqtt: {
    boolean: boolean
  }
}

function CreateEventRuleForm(props: { onSubmit?: () => void }) {
  const [addMore, setAddMore] = createSignal(false)

  const [form, { Field, Form }] = createForm<CreateEventRuleForm>();
  const submit = async (data: CreateEventRuleForm) => {
    await useClient()
      .admin.createEventRule({
        ...data,
        ignoreDb: !data.db.boolean,
        ignoreLive: !data.live.boolean,
        ignoreMqtt: !data.mqtt.boolean,
      })
      .then(() => revalidate(getListEventRules.key))
      .catch(throwAsFormError)
      .then(() => {
        if (addMore()) {
          reset(form, {
            initialValues: {
              ...data,
              code: ""
            }
          })
        } else {
          props.onSubmit && props.onSubmit()
        }
      })
  }

  return (
    <ErrorBoundary fallback={(e) => <PageError error={e} />}>
      <Suspense fallback={<Skeleton class="h-32" />}>
        <Form class="flex flex-col gap-4" onSubmit={submit}>
          <Field name="code">
            {(field, props) => (
              <FieldRoot>
                <FieldLabel field={field}>Code</FieldLabel>
                <Input
                  {...props}
                  {...fieldControlProps(field)}
                  placeholder="Code"
                  value={field.value}
                />
                <FieldMessage field={field} />
              </FieldRoot>
            )}
          </Field>
          <div class="flex flex-wrap gap-4">
            <Field name="db.boolean" type="boolean">
              {(field, props) => (
                <CheckboxFieldRoot form={form} field={field} class="flex items-center gap-2">
                  <CheckboxControl inputProps={props} />
                  <CheckboxLabel>DB</CheckboxLabel>
                  <CheckboxErrorMessage>{field.error}</CheckboxErrorMessage>
                </CheckboxFieldRoot>
              )}
            </Field>
            <Field name="live.boolean" type="boolean">
              {(field, props) => (
                <CheckboxFieldRoot form={form} field={field} class="flex items-center gap-2">
                  <CheckboxControl inputProps={props} />
                  <CheckboxLabel>Live</CheckboxLabel>
                  <CheckboxErrorMessage>{field.error}</CheckboxErrorMessage>
                </CheckboxFieldRoot>
              )}
            </Field>
            <Field name="mqtt.boolean" type="boolean">
              {(field, props) => (
                <CheckboxFieldRoot form={form} field={field} class="flex items-center gap-2">
                  <CheckboxControl inputProps={props} />
                  <CheckboxLabel>MQTT</CheckboxLabel>
                  <CheckboxErrorMessage>{field.error}</CheckboxErrorMessage>
                </CheckboxFieldRoot>
              )}
            </Field>
          </div>
          <Button type="submit" disabled={form.submitting}>
            <Show when={!form.submitting} fallback="Creating device">Create device</Show>
          </Button>
          <FormMessage form={form} />
          <CheckboxRoot checked={addMore()} onChange={setAddMore} class="flex items-center gap-2">
            <CheckboxControl />
            <CheckboxLabel>Add more</CheckboxLabel>
          </CheckboxRoot>
        </Form>
      </Suspense>
    </ErrorBoundary>
  )
}
