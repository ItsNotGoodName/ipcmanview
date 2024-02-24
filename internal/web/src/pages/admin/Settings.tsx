import { createForm, reset } from "@modular-forms/solid";
import { ErrorBoundary, Show, Suspense, } from "solid-js";
import { Shared } from "~/components/Shared";
import { Button } from "~/ui/Button";
import { FieldControl, FieldDescription, FieldLabel, FieldMessage, FieldRoot, SwitchFieldRoot, FormMessage } from "~/ui/Form";
import { Input } from "~/ui/Input";
import { LayoutNormal } from "~/ui/Layout";
import { SwitchControl, SwitchDescription, SwitchErrorMessage, SwitchLabel } from "~/ui/Switch";
import { useClient } from "~/providers/client";
import { createAsync, revalidate } from "@solidjs/router";
import { getConfig } from "../data";
import { GetConfigResp } from "~/twirp/rpc";
import { Skeleton } from "~/ui/Skeleton";
import { throwAsFormError } from "~/lib/utils";
import { PageError } from "~/ui/Page";

type UpdateForm = {
  siteName: string
  enableSignUp: {
    boolean: boolean
  }
}

export function AdminSettings() {
  const config = createAsync(() => getConfig())
  const refetchConfig = () => revalidate(getConfig.key)

  return (
    <LayoutNormal class="max-w-4xl">
      <Shared.Title>Settings</Shared.Title>
      <ErrorBoundary fallback={(e) => <PageError error={e} />}>
        <Suspense fallback={<Skeleton class="h-32" />}>
          <Show when={config()}>
            <UpdateSettingsForm config={config()!} refetchConfig={refetchConfig} />
          </Show>
        </Suspense>
      </ErrorBoundary>
    </LayoutNormal>
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
          <Button type="submit" disabled={updateForm.submitting} class="flex-1">
            <Show when={!updateForm.submitting} fallback="Updating settings">Update settings</Show>
          </Button>
          <Button type="button" onClick={formReset} variant="destructive" disabled={updateForm.submitting}>Reset</Button>
        </div>
        <FormMessage form={updateForm} />
      </Form>
    </div>
  )
}
