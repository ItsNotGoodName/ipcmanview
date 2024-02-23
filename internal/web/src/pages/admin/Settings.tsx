import { createForm, reset } from "@modular-forms/solid";
import { createAsync, revalidate } from "@solidjs/router";
import { Show, createEffect } from "solid-js";
import { Shared } from "~/components/Shared";
import { Button } from "~/ui/Button";
import { FieldControl, FieldDescription, FieldLabel, FieldMessage, FieldRoot, FieldSwitchRoot, FormMessage } from "~/ui/Form";
import { Input } from "~/ui/Input";
import { LayoutNormal } from "~/ui/Layout";
import { getConfig } from "../data";
import { SwitchControl, SwitchDescription, SwitchErrorMessage, SwitchLabel } from "~/ui/Switch";
import { useClient } from "~/providers/client";

type UpdateForm = {
  siteName: string
  enableSignUp: {
    boolean: boolean
  }
}

export function AdminSettings() {
  const [updateForm, { Field, Form }] = createForm<UpdateForm>();
  const config = createAsync(() => getConfig())
  createEffect(() => {
    const cfg = config()
    if (!cfg || updateForm.touched) return
    reset(updateForm, {
      initialValues: {
        siteName: cfg.siteName,
        enableSignUp: {
          boolean: cfg.enableSignUp
        }
      }
    })
  })

  const formRefresh = () => revalidate(getConfig.key)
  const formDisabled = () => updateForm.submitting
  const formSubmit = (form: UpdateForm) =>
    useClient().admin.updateConfig({
      enableSignUp: form.enableSignUp.boolean,
      siteName: form.siteName,
    }).then(formRefresh)

  return (
    <LayoutNormal class="max-w-4xl">
      <Shared.Title>Settings</Shared.Title>
      <div class="flex justify-center">
        <Form class="flex w-full max-w-sm flex-col gap-4" onSubmit={formSubmit}>
          <Field name="siteName">
            {(field, props) => (
              <FieldRoot class="gap-1.5">
                <FieldLabel field={field}>Site name</FieldLabel>
                <FieldControl field={field}>
                  <Input
                    {...props}
                    value={field.value}
                  />
                </FieldControl>
                <FieldDescription>Name of site.</FieldDescription>
                <FieldMessage field={field} />
              </FieldRoot>
            )}
          </Field>
          <Field name="enableSignUp.boolean" type="boolean">
            {(field) => (
              <FieldSwitchRoot
                form={updateForm}
                field={field}
                class="flex items-center justify-between gap-2"
              >
                <div>
                  <SwitchLabel>Enable sign up</SwitchLabel>
                  <SwitchDescription>Allow public sign up.</SwitchDescription>
                  <SwitchErrorMessage>{field.error}</SwitchErrorMessage>
                </div>
                <SwitchControl />
              </FieldSwitchRoot>
            )}
          </Field>
          <div class="flex gap-4">
            <Button type="button" onClick={formRefresh} variant="destructive">Refresh</Button>
            <Button type="submit" disabled={formDisabled()} class="flex-1">
              <Show when={!updateForm.submitting} fallback="Updating settings">Update settings</Show>
            </Button>
          </div>
          <FormMessage form={updateForm} />
        </Form>
      </div>
    </LayoutNormal>
  )
}
