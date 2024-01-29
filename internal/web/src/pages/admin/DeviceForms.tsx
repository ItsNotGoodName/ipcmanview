import { FieldControl, FieldLabel, FieldMessage, FieldRoot, FormMessage } from "~/ui/Form";
import { Input } from "~/ui/Input";
import { createForm, required, reset } from "@modular-forms/solid";
import { getListDeviceFeatures, getListLocations } from "./data";
import { SelectHTML } from "~/ui/Select";
import { action, createAsync, useAction, useNavigate } from "@solidjs/router";
import { ErrorBoundary, For, Show, Suspense, createResource, createSignal } from "solid-js";
import { Button } from "~/ui/Button";
import { PageError } from "~/ui/Page";
import { setupForm, throwAsFormError } from "~/lib/utils";
import { useClient } from "~/providers/client";
import { CheckboxControl, CheckboxInput, CheckboxLabel, CheckboxRoot } from "~/ui/Checkbox";
import { Skeleton } from "~/ui/Skeleton";
import { LayoutNormal } from "~/ui/Layout";
import { Seperator } from "~/ui/Seperator";
import { getDevice } from "./DeviceForms.data";

type CreateDeviceForm = {
  name: string
  url: string
  username: string
  password: string
  location: string
  features: {
    array: string[]
  }
}

const actionCreateDevice = action((form: CreateDeviceForm) => useClient()
  .admin.createDevice({ ...form, features: form.features.array })
  .then((value => value.response.id))
  .catch(throwAsFormError))

export function AdminDevicesCreate() {
  const navigate = useNavigate()
  const [addMore, setAddMore] = createSignal(false)

  const [createDeviceForm, { Field, Form }] = createForm<CreateDeviceForm>();
  const createDeviceAction = useAction(actionCreateDevice)
  const submit = async (form: CreateDeviceForm) => {
    const id = await createDeviceAction(form)

    if (addMore()) {
      reset(createDeviceForm, {
        initialValues: {
          ...form,
          name: "",
          url: ""
        },
      })
    } else {
      navigate(`./${id}`)
    }
  }

  const locations = createAsync(getListLocations)
  const deviceFeatures = createAsync(getListDeviceFeatures)

  return (
    <LayoutNormal class="max-w-lg">
      <div class="text-xl">Create device</div>
      <Seperator />

      <ErrorBoundary fallback={(e: Error) => <PageError error={e} />}>
        <Suspense fallback={<Skeleton class="h-32" />}>
          <Form class="flex flex-col gap-4" onSubmit={submit}>
            <Field name="name">
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
            <Field name="url" validate={required("Please enter a URL.")}>
              {(field, props) => (
                <FieldRoot class="gap-1.5">
                  <FieldLabel field={field}>URL</FieldLabel>
                  <FieldControl field={field}>
                    <Input
                      {...props}
                      placeholder="URL"
                      value={field.value}
                    />
                  </FieldControl>
                  <FieldMessage field={field} />
                </FieldRoot>
              )}
            </Field>
            <Field name="username">
              {(field, props) => (
                <FieldRoot class="gap-1.5">
                  <FieldLabel field={field}>Username</FieldLabel>
                  <FieldControl field={field}>
                    <Input
                      {...props}
                      placeholder="Username"
                      value={field.value}
                    />
                  </FieldControl>
                  <FieldMessage field={field} />
                </FieldRoot>
              )}
            </Field>
            <Field name="password">
              {(field, props) => (
                <FieldRoot class="gap-1.5">
                  <FieldLabel field={field}>Password</FieldLabel>
                  <FieldControl field={field}>
                    <Input
                      {...props}
                      autocomplete="off"
                      placeholder="Password"
                      type="password"
                      value={field.value}
                    />
                  </FieldControl>
                  <FieldMessage field={field} />
                </FieldRoot>
              )}
            </Field>
            <Field name="location">
              {(field, props) => (
                <FieldRoot class="gap-1.5">
                  <FieldLabel field={field}>Location</FieldLabel>
                  <FieldControl field={field}>
                    <SelectHTML
                      {...props}
                      value={field.value}
                    >
                      <For each={locations()}>
                        {v =>
                          <option value={v}>{v}</option>
                        }
                      </For>
                    </SelectHTML>
                  </FieldControl>
                  <FieldMessage field={field} />
                </FieldRoot>
              )}
            </Field>
            <Field name="features.array" type="string[]">
              {(field, props) => (
                <FieldRoot class="gap-1.5">
                  <FieldLabel field={field}>Features</FieldLabel>
                  <FieldControl field={field}>
                    <SelectHTML
                      {...props}
                      class="h-32"
                      multiple
                    >
                      <option value="" hidden selected disabled>
                        Features
                      </option>
                      <For each={deviceFeatures()}>
                        {v =>
                          <option value={v.value} selected={field.value?.includes(v.value)}>
                            {v.name}
                          </option>
                        }
                      </For>
                    </SelectHTML>
                  </FieldControl>
                  <FieldMessage field={field} />
                </FieldRoot>
              )}
            </Field>
            <Button type="submit" disabled={createDeviceForm.submitting}>
              <Show when={!createDeviceForm.submitting} fallback={<>Creating device</>}>
                Create device
              </Show>
            </Button>
            <FormMessage form={createDeviceForm} />
            <CheckboxRoot checked={addMore()} onChange={setAddMore}>
              <CheckboxInput />
              <CheckboxControl />
              <CheckboxLabel>Add more</CheckboxLabel>
            </CheckboxRoot>
          </Form>
        </Suspense>
      </ErrorBoundary>
    </LayoutNormal>
  )
}

type UpdateDeviceForm = {
  id: any
} & CreateDeviceForm

const actionUpdateDevice = action((form: UpdateDeviceForm) => useClient()
  .admin.updateDevice({ ...form, features: form.features.array })
  .then()
  .catch(throwAsFormError))

export function AdminDevicesIDUpdate(props: any) {
  const [updateDeviceForm, { Field, Form }] = createForm<UpdateDeviceForm>();
  const updateDeviceAction = useAction(actionUpdateDevice)
  const submit = (form: UpdateDeviceForm) => updateDeviceAction(form)

  const [form] = createResource(() => getDevice(props.params.id)
    .then((data) => setupForm(updateDeviceForm, { ...data, features: { array: data?.features || [] } })))

  const locations = createAsync(getListLocations)
  const deviceFeatures = createAsync(getListDeviceFeatures)

  return (
    <LayoutNormal class="max-w-lg">
      <div class="text-xl">Update device</div>
      <Seperator />

      <ErrorBoundary fallback={(e: Error) => <PageError error={e} />}>
        <Show when={!form.error} fallback={<PageError error={form.error} />}>
          <Form class="flex flex-col gap-4" onSubmit={(form) => submit(form)}>
            <Field name="id" type="number">
              {(field, props) => <input {...props} type="hidden" value={field.value} />}
            </Field>
            <Field name="name">
              {(field, props) => (
                <FieldRoot class="gap-1.5">
                  <FieldLabel field={field}>Name</FieldLabel>
                  <FieldControl field={field}>
                    <Input
                      {...props}
                      placeholder="Name"
                      value={field.value}
                      disabled={form.loading}
                    />
                  </FieldControl>
                  <FieldMessage field={field} />
                </FieldRoot>
              )}
            </Field>
            <Field name="url" validate={required("Please enter a URL.")}>
              {(field, props) => (
                <FieldRoot class="gap-1.5">
                  <FieldLabel field={field}>URL</FieldLabel>
                  <FieldControl field={field}>
                    <Input
                      {...props}
                      placeholder="URL"
                      value={field.value}
                      disabled={form.loading}
                    />
                  </FieldControl>
                  <FieldMessage field={field} />
                </FieldRoot>
              )}
            </Field>
            <Field name="username">
              {(field, props) => (
                <FieldRoot class="gap-1.5">
                  <FieldLabel field={field}>Username</FieldLabel>
                  <FieldControl field={field}>
                    <Input
                      {...props}
                      placeholder="Username"
                      value={field.value}
                      disabled={form.loading}
                    />
                  </FieldControl>
                  <FieldMessage field={field} />
                </FieldRoot>
              )}
            </Field>
            <Field name="password">
              {(field, props) => (
                <FieldRoot class="gap-1.5">
                  <FieldLabel field={field}>Password</FieldLabel>
                  <FieldControl field={field}>
                    <Input
                      {...props}
                      autocomplete="off"
                      placeholder="Password"
                      type="password"
                      value={field.value}
                      disabled={form.loading}
                    />
                  </FieldControl>
                  <FieldMessage field={field} />
                </FieldRoot>
              )}
            </Field>
            <Field name="location">
              {(field, props) => (
                <FieldRoot class="gap-1.5">
                  <FieldLabel field={field}>Location</FieldLabel>
                  <FieldControl field={field}>
                    <SelectHTML
                      {...props}
                      value={field.value}
                      disabled={form.loading}
                    >
                      <For each={locations()}>
                        {v =>
                          <option value={v}>{v}</option>
                        }
                      </For>
                    </SelectHTML>
                  </FieldControl>
                  <FieldMessage field={field} />
                </FieldRoot>
              )}
            </Field>
            <Field name="features.array" type="string[]">
              {(field, props) => (
                <FieldRoot class="gap-1.5">
                  <FieldLabel field={field}>Features</FieldLabel>
                  <FieldControl field={field}>
                    <SelectHTML
                      {...props}
                      class="h-32"
                      multiple
                      disabled={form.loading}
                    >
                      <option value="" hidden selected disabled>
                        Features
                      </option>
                      <For each={deviceFeatures()}>
                        {v =>
                          <option value={v.value} selected={field.value?.includes(v.value)}>
                            {v.name}
                          </option>
                        }
                      </For>
                    </SelectHTML>
                  </FieldControl>
                  <FieldMessage field={field} />
                </FieldRoot>
              )}
            </Field>
            <Button type="submit" disabled={form.loading || updateDeviceForm.submitting}>
              <Show when={!updateDeviceForm.submitting} fallback={<>Updating device</>}>
                Update device
              </Show>
            </Button>
            <FormMessage form={updateDeviceForm} />
          </Form>
        </Show>
      </ErrorBoundary>
    </LayoutNormal>
  )
}
