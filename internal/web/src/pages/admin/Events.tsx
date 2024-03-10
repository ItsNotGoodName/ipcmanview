import Humanize from "humanize-plus"
import { createForm, reset } from "@modular-forms/solid";
import { ErrorBoundary, For, Show, Suspense, batch, createSignal, } from "solid-js";
import { Shared } from "~/components/Shared";
import { Button } from "~/ui/Button";
import { FormMessage } from "~/ui/Form";
import { LayoutNormal } from "~/ui/Layout";
import { useClient } from "~/providers/client";
import { action, createAsync, revalidate, useAction, useSubmission } from "@solidjs/router";
import { ListEventRulesResp_Item, UpdateEventRuleReq_Item } from "~/twirp/rpc";
import { Skeleton } from "~/ui/Skeleton";
import { catchAsToast, createModal, createRowSelection, setFormValue, throwAsFormError, validationState } from "~/lib/utils";
import { PageError } from "~/ui/Page";
import { TableBody, TableCell, TableHead, TableHeader, TableRoot, TableRow } from "~/ui/Table";
import { CheckboxControl, CheckboxErrorMessage, CheckboxLabel, CheckboxRoot } from "~/ui/Checkbox";
import { getListEventRules } from "./data";
import { DialogContent, DialogHeader, DialogOverflow, DialogOverlay, DialogPortal, DialogRoot, DialogTitle } from "~/ui/Dialog";
import { RiDeviceSaveLine, RiSystemAddLine, RiSystemDeleteBinLine, RiSystemRefreshLine } from "solid-icons/ri";
import { AlertDialogAction, AlertDialogCancel, AlertDialogDescription, AlertDialogFooter, AlertDialogHeader, AlertDialogModal, AlertDialogRoot, AlertDialogTitle } from "~/ui/AlertDialog";
import { createStore } from "solid-js/store";
import { TextFieldDescription, TextFieldErrorMessage, TextFieldInput, TextFieldLabel, TextFieldRoot } from "~/ui/TextField";
import { getAdminEventsPage } from "./Events.data";

const actionDeleteEvents = action(() => useClient()
  .admin.deleteEvents({})
  .then()
  .catch(catchAsToast))

export function AdminEvents() {
  const data = createAsync(() => getAdminEventsPage())
  const eventRules = createAsync(() => getListEventRules())
  const refetchEventRules = () => revalidate(getListEventRules.key)

  const deleteEventsModal = createModal(0)
  const deleteEventsSubmission = useSubmission(actionDeleteEvents)
  const deleteEventsAction = useAction(actionDeleteEvents)
  const submitDeleteEvents = () => deleteEventsAction().then(deleteEventsModal.setClose)

  return (
    <LayoutNormal class="max-w-4xl">
      <AlertDialogRoot open={deleteEventsModal.open()} onOpenChange={() => deleteEventsModal.setClose()}>
        <AlertDialogModal>
          <AlertDialogHeader>
            <AlertDialogTitle>Are you sure you wish to delete {deleteEventsModal.value()} {Humanize.pluralize(deleteEventsModal.value(), "event")}?</AlertDialogTitle>
          </AlertDialogHeader>
          <AlertDialogFooter>
            <AlertDialogCancel>Cancel</AlertDialogCancel>
            <AlertDialogAction variant="destructive" disabled={deleteEventsSubmission.pending} onClick={submitDeleteEvents}>
              Delete
            </AlertDialogAction>
          </AlertDialogFooter>
        </AlertDialogModal>
      </AlertDialogRoot>

      <Shared.Title>Events</Shared.Title>
      <ErrorBoundary fallback={(e) => <PageError error={e} />}>
        <div class="flex gap-2 flex-wrap">
          <Button
            onClick={() => deleteEventsModal.setValue(Number(data()?.eventCount))}
            disabled={deleteEventsSubmission.pending}
            variant="destructive"
          >
            Delete {data()?.eventCount} {Humanize.pluralize(deleteEventsModal.value(), "event")}
          </Button>
        </div>
        <Shared.Title>Rules</Shared.Title>
        <Suspense fallback={<Skeleton class="h-32" />}>
          <Show when={eventRules()}>
            <EventRulesTable eventRules={eventRules()!} refetchEventRules={refetchEventRules} />
          </Show>
        </Suspense>
      </ErrorBoundary>
    </LayoutNormal>
  )
}

const actionDeleteEventRule = action((ids: string[]) => useClient()
  .admin.deleteEventRules({ ids })
  .then(() => true)
  .catch(catchAsToast))

const actionUpdateEventRule = action((items: UpdateEventRuleReq_Item[]) => useClient()
  .admin.updateEventRule({ items })
  .then(() => true)
  .catch(catchAsToast))

function EventRulesTable(props: { eventRules: ListEventRulesResp_Item[], refetchEventRules: () => Promise<void> }) {
  const [rows, setRows] = createStore<(ListEventRulesResp_Item & { _dirty: boolean })[]>(props.eventRules.map(v => ({ ...v, _dirty: false })))
  const resetRows = () => setRows(props.eventRules.map(v => ({ ...v, _dirty: false })))
  const rowsDirty = () => rows.filter(v => v._dirty)

  const rowSelection = createRowSelection(() => rows?.map(v => ({ id: v.id, disabled: v.code == "" })) || [])

  const submitRefresh = () => props.refetchEventRules().then(resetRows)

  // Create
  const [createFormModal, setCreateFormModal] = createSignal(false)

  // Update
  const updateSubmission = useSubmission(actionUpdateEventRule)
  const updateAction = useAction(actionUpdateEventRule)
  const submitUpdate = () => updateAction(rowsDirty())
    .then((value) => value === true && resetRows())

  // Delete
  const deleteSubmission = useSubmission(actionDeleteEventRule)
  const deleteAction = useAction(actionDeleteEventRule)
  // Single
  const deleteModal = createModal({ name: "", id: "" })
  const deleteSubmit = () => deleteAction([deleteModal.value().id])
    .then((value) => value === true &&
      batch(() => {
        deleteModal.setClose()
        resetRows()
      }))
  // Multiple
  const deleteMultipleModal = createModal<ListEventRulesResp_Item[]>([])
  const deleteMultipleSubmit = () => deleteAction(rowSelection.selections())
    .then((value) => value === true &&
      batch(() => {
        deleteMultipleModal.setClose()
        resetRows()
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
              <CreateEventRuleForm onSubmit={submitRefresh} onClose={() => setCreateFormModal(false)} />
            </DialogOverflow>
          </DialogContent>
        </DialogPortal>
      </DialogRoot>

      <AlertDialogRoot open={deleteModal.open()} onOpenChange={deleteModal.setClose}>
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

      <AlertDialogRoot open={deleteMultipleModal.open()} onOpenChange={deleteMultipleModal.setClose}>
        <AlertDialogModal>
          <AlertDialogHeader>
            <AlertDialogTitle>Are you sure you wish to delete {deleteMultipleModal.value().length} event {Humanize.pluralize(deleteMultipleModal.value().length, "rule")}?</AlertDialogTitle>
            <AlertDialogDescription>
              <ul>
                <For each={deleteMultipleModal.value()}>
                  {(e) =>
                    <li>{e.code}</li>
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
        <Button
          size="icon"
          variant="outline"
          onClick={() => setCreateFormModal(true)}
        >
          <RiSystemAddLine class="size-5" />
        </Button>
        <Button
          size="icon"
          title="Refresh"
          variant="secondary"
          onClick={submitRefresh}
        >
          <RiSystemRefreshLine class="size-5" />
        </Button>
        <Button
          size="icon"
          title="Update"
          disabled={!rows.some(v => v._dirty) || updateSubmission.pending}
          onClick={submitUpdate}
        >
          <RiDeviceSaveLine class="size-5" />
        </Button>
        <Button
          size="icon"
          variant="destructive"
          title="Delete"
          disabled={!rowSelection.multiple() || deleteSubmission.pending}
          onClick={() => deleteMultipleModal.setValue(rows.filter(v => rowSelection.selections().includes(v.id)))}
        >
          <RiSystemDeleteBinLine class="size-5" />
        </Button>
      </div>

      <TableRoot>
        <TableHeader>
          <TableRow>
            <TableHead>
              <CheckboxRoot
                indeterminate={rowSelection.multiple()}
                checked={rowSelection.all()}
                onChange={rowSelection.setAll}
              >
                <CheckboxControl />
              </CheckboxRoot>
            </TableHead>
            <TableHead>Code</TableHead>
            <TableHead>
              <button onClick={() => {
                const value = rows.some(v => v.ignoreDb)
                setRows(() => true, (v) => ({ ...v, _dirty: true, ignoreDb: !value }))
              }}>
                DB
              </button>
            </TableHead>
            <TableHead>
              <button onClick={() => {
                const value = rows.some(v => v.ignoreLive)
                setRows(() => true, (v) => ({ ...v, _dirty: true, ignoreLive: !value }))
              }}>
                Live
              </button>
            </TableHead>
            <TableHead>
              <button onClick={() => {
                const value = rows.some(v => v.ignoreMqtt)
                setRows(() => true, (v) => ({ ...v, _dirty: true, ignoreMqtt: !value }))
              }}>
                MQTT
              </button>
            </TableHead>
          </TableRow>
        </TableHeader>
        <TableBody>
          <For each={rows}>
            {(item, index) =>
              <TableRow>
                <TableCell>
                  <CheckboxRoot
                    disabled={!item.code}
                    checked={rowSelection.items[index()]?.checked}
                    onChange={(checked) => rowSelection.set(item.id, checked)}
                  >
                    <CheckboxControl />
                  </CheckboxRoot>
                </TableCell>
                <Show when={item.code} fallback={<TableCell class="w-full">All</TableCell>} >
                  <td class="w-full min-w-32 py-0 align-middle">
                    <TextFieldRoot
                      value={item.code}
                      onChange={(value) => setRows(
                        (todo) => todo.id === item.id,
                        (v) => ({ ...v, _dirty: true, code: value })
                      )}
                    >
                      <TextFieldInput />
                    </TextFieldRoot>
                  </td>
                </Show>
                <TableCell>
                  <CheckboxRoot
                    checked={!item.ignoreDb}
                    onChange={(value) => setRows(
                      (todo) => todo.id === item.id,
                      (v) => ({ ...v, _dirty: true, ignoreDb: !value })
                    )}
                  >
                    <CheckboxControl />
                  </CheckboxRoot>
                </TableCell>
                <TableCell>
                  <CheckboxRoot
                    checked={!item.ignoreLive}
                    onChange={(value) => setRows(
                      (todo) => todo.id === item.id,
                      (v) => ({ ...v, _dirty: true, ignoreLive: !value })
                    )}
                  >
                    <CheckboxControl />
                  </CheckboxRoot>
                </TableCell>
                <TableCell>
                  <CheckboxRoot
                    checked={!item.ignoreMqtt}
                    onChange={(value) => setRows(
                      (todo) => todo.id === item.id,
                      (v) => ({ ...v, _dirty: true, ignoreMqtt: !value })
                    )}
                  >
                    <CheckboxControl />
                  </CheckboxRoot>
                </TableCell>
              </TableRow>
            }
          </For>
        </TableBody>
      </TableRoot>
    </div >
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

function CreateEventRuleForm(props: { onSubmit: () => void, onClose: () => void }) {
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
      .then(props.onSubmit)
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
          props.onClose()
        }
      })
  }

  return (
    <ErrorBoundary fallback={(e) => <PageError error={e} />}>
      <Suspense fallback={<Skeleton class="h-32" />}>
        <Form class="flex flex-col gap-4" onSubmit={submit}>
          <Field name="code">
            {(field, props) => (
              <TextFieldRoot
                validationState={validationState(field.error)}
                value={field.value}
                class="space-y-2"
              >
                <TextFieldLabel>Code</TextFieldLabel>
                <TextFieldInput
                  {...props}
                  placeholder="Code"
                />
                <TextFieldDescription>Match by code.</TextFieldDescription>
                <TextFieldErrorMessage>{field.error}</TextFieldErrorMessage>
              </TextFieldRoot>
            )}
          </Field>
          <div class="flex flex-wrap gap-4">
            <Field name="db.boolean" type="boolean">
              {(field) => (
                <CheckboxRoot
                  validationState={validationState(field.error)}
                  checked={field.value}
                  onChange={setFormValue(form, field)}
                  class="flex items-center gap-2"
                >
                  <CheckboxControl />
                  <CheckboxLabel>DB</CheckboxLabel>
                  <CheckboxErrorMessage>{field.error}</CheckboxErrorMessage>
                </CheckboxRoot>
              )}
            </Field>
            <Field name="live.boolean" type="boolean">
              {(field) => (
                <CheckboxRoot
                  validationState={validationState(field.error)}
                  checked={field.value}
                  onChange={setFormValue(form, field)}
                  class="flex items-center gap-2"
                >
                  <CheckboxControl />
                  <CheckboxLabel>Live</CheckboxLabel>
                  <CheckboxErrorMessage>{field.error}</CheckboxErrorMessage>
                </CheckboxRoot>
              )}
            </Field>
            <Field name="mqtt.boolean" type="boolean">
              {(field) => (
                <CheckboxRoot
                  validationState={validationState(field.error)}
                  checked={field.value}
                  onChange={setFormValue(form, field)}
                  class="flex items-center gap-2"
                >
                  <CheckboxControl />
                  <CheckboxLabel>MQTT</CheckboxLabel>
                  <CheckboxErrorMessage>{field.error}</CheckboxErrorMessage>
                </CheckboxRoot>
              )}
            </Field>
          </div>
          <Button type="submit" disabled={form.submitting}>
            <Show when={!form.submitting} fallback="Creating event rule">Create event rule</Show>
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

export default AdminEvents
