import { createForm, reset } from "@modular-forms/solid";
import { ErrorBoundary, Show, Suspense, } from "solid-js";
import { Shared } from "~/components/Shared";
import { Button } from "~/ui/Button";
import { FormMessage } from "~/ui/Form";
import { LayoutNormal } from "~/ui/Layout";
import { SwitchControl, SwitchDescription, SwitchErrorMessage, SwitchLabel, SwitchRoot } from "~/ui/Switch";
import { useClient } from "~/providers/client";
import { createAsync, revalidate, } from "@solidjs/router";
import { getConfig } from "../data";
import { GetConfigResp } from "~/twirp/rpc";
import { Skeleton } from "~/ui/Skeleton";
import { createLoading, setFormValue, throwAsFormError, validationState } from "~/lib/utils";
import { PageError } from "~/ui/Page";
import { TextFieldDescription, TextFieldErrorMessage, TextFieldInput, TextFieldLabel, TextFieldRoot } from "~/ui/TextField";

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
  const [form, { Field, Form }] = createForm<UpdateForm>({ initialValues: formInitialValues() });
  const [refreshFormLoading, refreshForm] = createLoading(() => props.refetchConfig().then(() => reset(form, { initialValues: formInitialValues() })))
  const submitForm = (form: UpdateForm) => useClient()
    .admin.updateConfig({
      enableSignUp: form.enableSignUp.boolean,
      siteName: form.siteName,
    })
    .then(refreshForm)
    .catch(throwAsFormError)
  const formDisabled = () => refreshFormLoading() || form.submitting

  return (
    <div class="flex justify-center">
      <Form class="flex w-full max-w-sm flex-col gap-4" onSubmit={submitForm}>
        <Field name="siteName">
          {(field, props) => (
            <TextFieldRoot
              validationState={validationState(field.error)}
              value={field.value}
              class="space-y-2"
            >
              <TextFieldLabel>Site name</TextFieldLabel>
              <TextFieldInput
                {...props}
                placeholder="Site name"
              />
              <TextFieldDescription>Name of site.</TextFieldDescription>
              <TextFieldErrorMessage>{field.error}</TextFieldErrorMessage>
            </TextFieldRoot>
          )}
        </Field>
        <Field name="enableSignUp.boolean" type="boolean">
          {(field, props) => (
            <SwitchRoot
              validationState={validationState(field.error)}
              checked={field.value}
              onChange={setFormValue(form, field)}
              class="flex items-center justify-between gap-2"
            >
              <div>
                <SwitchLabel>Enable sign up</SwitchLabel>
                <SwitchDescription>Allow public to sign up.</SwitchDescription>
                <SwitchErrorMessage>{field.error}</SwitchErrorMessage>
              </div>
              <SwitchControl inputProps={props} />
            </SwitchRoot>
          )}
        </Field>
        <div class="flex flex-col gap-4 sm:flex-row-reverse">
          <Button type="submit" disabled={formDisabled()} class="sm:flex-1">
            <Show when={!form.submitting} fallback="Updating settings">Update settings</Show>
          </Button>
          <Button type="button" onClick={refreshForm} variant="secondary" disabled={formDisabled()}>Refresh</Button>
        </div>
        <FormMessage form={form} />
      </Form>
    </div>
  )
}

export default AdminSettings
