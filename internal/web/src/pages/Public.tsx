import { createForm, required, reset } from "@modular-forms/solid";
import { A, createAsync, revalidate, useNavigate } from "@solidjs/router";
import { ParentProps, Show, } from "solid-js";

import { useClient } from "~/providers/client";
import { Button } from "~/ui/Button";
import { CardRoot } from "~/ui/Card";
import { FieldRoot, FieldLabel, FieldMessage, FormMessage, CheckboxFieldRoot, fieldControlProps } from "~/ui/Form";
import { Input } from "~/ui/Input";
import { linkVariants } from "~/ui/Link";
import { ThemeIcon } from "~/ui/ThemeIcon";
import { toggleTheme, useThemeTitle } from "~/ui/theme";
import { CheckboxControl, CheckboxErrorMessage, CheckboxLabel } from "~/ui/Checkbox";
import { throwAsFormError } from "~/lib/utils";
import { toast } from "~/ui/Toast";
import { getSession } from "~/providers/session";
import { AlertDescription, AlertRoot, AlertTitle } from "~/ui/Alert";
import { getConfig } from "./data";

function Layout(props: ParentProps) {
  return (
    <div class="mx-auto flex max-w-sm flex-col gap-4 p-4 pt-10">
      {props.children}
    </div>
  )
}

function Header(props: ParentProps) {
  return (
    <div class="text-center ">
      <div class="text-2xl">IPCManView</div>
      <div class="text-muted-foreground text-sm">
        {props.children}
      </div>
    </div>
  )
}

function CardHeader(props: ParentProps) {
  return (
    <div class="flex items-center justify-between gap-2">
      <div class="flex-1"></div>
      <div class="text-xl">
        {props.children}
      </div>
      <div class="flex flex-1 items-center justify-end">
        <Button size='icon' variant='ghost' onClick={toggleTheme} title={useThemeTitle()}>
          <ThemeIcon class="h-6 w-6" />
        </Button>
      </div>
    </div>
  )
}

function Footer(props: ParentProps) {
  return (
    <CardRoot class="p-4">
      <div class="flex items-center justify-center">
        {props.children}
      </div>
    </CardRoot>
  )
}

type SignInForm = {
  usernameOrEmail: string
  password: string
  rememberMe: boolean
}

export function SignIn() {
  const navigate = useNavigate()

  const config = createAsync(() => getConfig())
  const session = createAsync(() => getSession())

  const [form, { Field, Form }] = createForm<SignInForm>();
  const submit = (input: SignInForm) =>
    fetch("/v1/session", {
      credentials: "include",
      headers: [['Content-Type', 'application/json'], ['Accept', 'application/json']],
      method: "POST",
      body: JSON.stringify(input),
    }).then(async (resp) => {
      if (!resp.ok) {
        const json = await resp.json()
        throw new Error(json.message)
      }

      await revalidate(getSession.key)
      navigate('/', { replace: true })
    }).catch(throwAsFormError)

  return (
    <Layout>
      <Header>{config()?.siteName}</Header>
      <CardRoot class="flex flex-col gap-4 p-4">
        <CardHeader>Sign in</CardHeader>
        <Form class="flex flex-col gap-4" onSubmit={submit}>
          <Field name="usernameOrEmail" validate={required("Please enter your username or email.")}>
            {(field, props) => (
              <FieldRoot>
                <FieldLabel field={field}>Username or email</FieldLabel>
                <Input
                  {...props}
                  {...fieldControlProps(field)}
                  autocomplete="username"
                  placeholder="Username or email"
                  value={field.value}
                />
                <FieldMessage field={field} />
              </FieldRoot>
            )}
          </Field>
          <Field name="password">
            {(field, props) => (
              <FieldRoot>
                <div class="flex items-center justify-between gap-2">
                  <FieldLabel field={field}>
                    Password
                  </FieldLabel>
                  <A href="/forgot" class={linkVariants()}>
                    Forgot password?
                  </A>
                </div>
                <Input
                  {...props}
                  {...fieldControlProps(field)}
                  autocomplete="current-password"
                  placeholder="Password"
                  type="password"
                  value={field.value}
                />
                <FieldMessage field={field} />
              </FieldRoot>
            )}
          </Field>
          <Field name="rememberMe" type="boolean">
            {(field, props) => (
              <CheckboxFieldRoot form={form} field={field} class="space-y-2">
                <div class="flex items-center gap-2">
                  <CheckboxControl inputProps={props} />
                  <CheckboxLabel>Remember me</CheckboxLabel>
                </div>
                <CheckboxErrorMessage>{field.error}</CheckboxErrorMessage>
              </CheckboxFieldRoot>
            )}
          </Field>
          <Button type="submit" disabled={form.submitting}>
            <Show when={!form.submitting} fallback="Signing in">Sign in</Show>
          </Button>
          <FormMessage form={form} />
        </Form>
        <Show when={session()?.valid && session()?.disabled}>
          <AlertRoot variant="destructive">
            <AlertTitle>Account disabled</AlertTitle>
            <AlertDescription>
              Your account "{session()?.username}" is disabled.
            </AlertDescription>
          </AlertRoot>
        </Show>
      </CardRoot>
      <Show when={config()?.enableSignUp}>
        <Footer>
          <A href="/signup" class={linkVariants()}>Sign up</A>
        </Footer>
      </Show>
    </Layout>
  )
}

type SignUpForm = {
  email: string
  username: string
  password: string
  confirmPassword: string
}

export function SignUp() {
  const navigate = useNavigate()

  const config = createAsync(() => getConfig())

  const [form, { Field, Form }] = createForm<SignUpForm>({
    validate: (input) => {
      if (input.password != input.confirmPassword) {
        return {
          confirmPassword: "Password does not match."
        }
      }
      return {}
    }
  });
  const submit = (input: SignUpForm) => useClient()
    .public.signUp(input)
    .then()
    .catch(throwAsFormError)
    .then(() => navigate("/signin"))

  return (
    <Layout>
      <Header>{config()?.siteName}</Header>
      <CardRoot class="flex flex-col gap-4 p-4">
        <CardHeader>Sign up</CardHeader>
        <Form class="flex flex-col gap-4" onSubmit={submit}>
          <Field name="email" validate={required('Please enter your email.')}>
            {(field, props) => (
              <FieldRoot>
                <FieldLabel field={field}>Email</FieldLabel>
                <Input
                  {...props}
                  {...fieldControlProps(field)}
                  placeholder="Email"
                  type="email"
                  value={field.value}
                />
                <FieldMessage field={field} />
              </FieldRoot>
            )}
          </Field>
          <Field name="username" validate={required('Please enter a username.')}>
            {(field, props) => (
              <FieldRoot>
                <FieldLabel field={field}>Username</FieldLabel>
                <Input
                  {...props}
                  {...fieldControlProps(field)}
                  autocomplete="username"
                  placeholder="Username"
                  value={field.value}
                />
                <FieldMessage field={field} />
              </FieldRoot>
            )}
          </Field>
          <Field name="password" validate={required('Please enter a password.')}>
            {(field, props) => (
              <FieldRoot>
                <div class="flex items-center justify-between gap-2">
                  <FieldLabel field={field}>
                    Password
                  </FieldLabel>
                </div>
                <Input
                  {...props}
                  {...fieldControlProps(field)}
                  autocomplete="new-password"
                  placeholder="Password"
                  type="password"
                  value={field.value}
                />
                <FieldMessage field={field} />
              </FieldRoot>
            )}
          </Field>
          <Field name="confirmPassword" validate={required('Please confirm your password.')}>
            {(field, props) => (
              <FieldRoot>
                <div class="flex items-center justify-between gap-2">
                  <FieldLabel field={field}>
                    Confirm password
                  </FieldLabel>
                </div>
                <Input
                  {...props}
                  {...fieldControlProps(field)}
                  autocomplete="new-password"
                  placeholder="Confirm password"
                  type="password"
                  value={field.value}
                />
                <FieldMessage field={field} />
              </FieldRoot>
            )}
          </Field>
          <Button type="submit" disabled={form.submitting}>
            <Show when={!form.submitting} fallback="Signing up">Sign up</Show>
          </Button>
          <FormMessage form={form} />
        </Form>
      </CardRoot>
      <Footer>
        <A href="/signin" class={linkVariants()}>Sign in</A>
      </Footer>
    </Layout>
  )
}

type ForgotForm = {
  email: string
}

export function Forgot() {
  const config = createAsync(() => getConfig())

  const [form, { Field, Form }] = createForm<ForgotForm>({ initialValues: { email: "" } });
  const submit = (input: ForgotForm) => useClient()
    .public.forgotPassword(input)
    .then(() => {
      toast.success("Sent password reset email.")
      reset(form)
    })
    .catch(throwAsFormError)

  return (
    <Layout>
      <Header>{config()?.siteName}</Header>
      <CardRoot class="flex flex-col gap-4 p-4">
        <CardHeader>Forgot</CardHeader>
        <Form class="flex flex-col gap-4" onSubmit={submit}>
          <Field name="email" validate={required('Please enter your email.')}>
            {(field, props) => (
              <FieldRoot>
                <FieldLabel field={field}>Email</FieldLabel>
                <Input
                  {...props}
                  {...fieldControlProps(field)}
                  placeholder="Email"
                  type="email"
                  value={field.value}
                />
                <FieldMessage field={field} />
              </FieldRoot>
            )}
          </Field>
          <Button type="submit" disabled={form.submitting}>
            <Show when={!form.submitting} fallback="Sending password reset email">Send password reset email</Show>
          </Button>
          <FormMessage form={form} />
        </Form>
      </CardRoot>
      <Footer>
        <A href="/signin" class={linkVariants()}>Sign in</A>
      </Footer>
    </Layout>
  )
}

