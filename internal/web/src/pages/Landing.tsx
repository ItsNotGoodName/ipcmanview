import { createForm, required, reset } from "@modular-forms/solid";
import { A, action, redirect, revalidate, useAction } from "@solidjs/router";
import { ParentProps, Show } from "solid-js";

import { useClient } from "~/providers/client";
import { Button } from "~/ui/Button";
import { CardRoot } from "~/ui/Card";
import { FieldControl, FieldRoot, FieldLabel, FieldMessage, FormMessage } from "~/ui/Form";
import { Input } from "~/ui/Input";
import { linkVariants } from "~/ui/Link";
import { ThemeIcon } from "~/ui/ThemeIcon";
import { toggleTheme, useThemeTitle } from "~/ui/theme";
import { CheckboxControl, CheckboxErrorMessage, CheckboxInput, CheckboxLabel, CheckboxRoot } from "~/ui/Checkbox";
import { throwAsFormError } from "~/lib/utils";
import { toast } from "~/ui/Toast";
import { getSession } from "~/providers/session";

function Layout(props: ParentProps) {
  return (
    <div class="mx-auto flex max-w-xs flex-col gap-4 pt-10">
      {props.children}
    </div>
  )
}

function Header() {
  return (
    <div class="text-center text-2xl">IPCManView</div>
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

const actionSignIn = action((form: SignInForm) =>
  fetch("/v1/session", {
    credentials: "include",
    headers: [['Content-Type', 'application/json'], ['Accept', 'application/json']],
    method: "POST",
    body: JSON.stringify(form),
  }).then(async (resp) => {
    if (!resp.ok) {
      const json = await resp.json()
      throw new Error(json.message)
    }
    return revalidate(getSession.key)
  }).catch(throwAsFormError)
)

export function SignIn() {
  const [signInForm, { Field, Form }] = createForm<SignInForm>();
  const signIn = useAction(actionSignIn)

  return (
    <Layout>
      <Header />
      <CardRoot class="flex flex-col gap-4 p-4">
        <CardHeader>Sign in</CardHeader>
        <Form class="flex flex-col gap-4" onSubmit={signIn}>
          <Field name="usernameOrEmail" validate={required("Please enter your username or email.")}>
            {(field, props) => (
              <FieldRoot class="gap-1.5">
                <FieldLabel field={field}>Username or email</FieldLabel>
                <FieldControl field={field}>
                  <Input
                    {...props}
                    autocomplete="username"
                    placeholder="Username or email"
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
                <div class="flex items-center justify-between gap-2">
                  <FieldLabel field={field}>
                    Password
                  </FieldLabel>
                  <A href="/forgot" class={linkVariants()}>
                    Forgot password?
                  </A>
                </div>
                <FieldControl field={field}>
                  <Input
                    {...props}
                    autocomplete="current-password"
                    placeholder="Password"
                    type="password"
                    value={field.value}
                  />
                </FieldControl>
                <FieldMessage field={field} />
              </FieldRoot>
            )}
          </Field>
          <Field name="rememberMe" type="boolean">
            {(field, props) => (
              <CheckboxRoot validationState={field.error ? "invalid" : "valid"}>
                <CheckboxInput {...props} />
                <CheckboxControl />
                <CheckboxLabel>Remember me</CheckboxLabel>
                <CheckboxErrorMessage>{field.error}</CheckboxErrorMessage>
              </CheckboxRoot>
            )}
          </Field>
          <Button type="submit" disabled={signInForm.submitting}>
            <Show when={!signInForm.submitting} fallback={<>Signing in</>}>
              Sign in
            </Show>
          </Button>
          <FormMessage form={signInForm} />
        </Form>
      </CardRoot>
      <Footer>
        <A href="/signup" class={linkVariants()}>Sign up</A>
      </Footer>
    </Layout>
  )
}

type SignUpForm = {
  email: string
  username: string
  password: string
  confirmPassword: string
}

const actionSignUp = action((form: SignUpForm) => useClient()
  .auth.signUp(form)
  .then()
  .catch(throwAsFormError)
  .then(async () => { throw redirect("/signin") }))

export function Signup() {
  const [signUpForm, { Field, Form }] = createForm<SignUpForm>({
    validate: (form) => {
      if (form.password != form.confirmPassword) {
        return {
          confirmPassword: "Password does not match."
        }
      }
      return {}
    }
  });
  const signUp = useAction(actionSignUp)

  return (
    <Layout>
      <Header />
      <CardRoot class="flex flex-col gap-4 p-4">
        <CardHeader>Sign up</CardHeader>
        <Form class="flex flex-col gap-4" onSubmit={signUp}>
          <Field name="email" validate={required('Please enter your email.')}>
            {(field, props) => (
              <FieldRoot class="gap-1.5">
                <FieldLabel field={field}>Email</FieldLabel>
                <FieldControl field={field}>
                  <Input
                    {...props}
                    placeholder="Email"
                    type="email"
                    value={field.value}
                  />
                </FieldControl>
                <FieldMessage field={field} />
              </FieldRoot>
            )}
          </Field>
          <Field name="username" validate={required('Please enter a username.')}>
            {(field, props) => (
              <FieldRoot class="gap-1.5">
                <FieldLabel field={field}>Username</FieldLabel>
                <FieldControl field={field}>
                  <Input
                    {...props}
                    autocomplete="username"
                    placeholder="Username"
                    value={field.value}
                  />
                </FieldControl>
                <FieldMessage field={field} />
              </FieldRoot>
            )}
          </Field>
          <Field name="password" validate={required('Please enter a password.')}>
            {(field, props) => (
              <FieldRoot class="gap-1.5">
                <div class="flex items-center justify-between gap-2">
                  <FieldLabel field={field}>
                    Password
                  </FieldLabel>
                </div>
                <FieldControl field={field}>
                  <Input
                    {...props}
                    autocomplete="new-password"
                    placeholder="Password"
                    type="password"
                    value={field.value}
                  />
                </FieldControl>
                <FieldMessage field={field} />
              </FieldRoot>
            )}
          </Field>
          <Field name="confirmPassword" validate={required('Please confirm your password.')}>
            {(field, props) => (
              <FieldRoot class="gap-1.5">
                <div class="flex items-center justify-between gap-2">
                  <FieldLabel field={field}>
                    Confirm password
                  </FieldLabel>
                </div>
                <FieldControl field={field}>
                  <Input
                    {...props}
                    autocomplete="new-password"
                    placeholder="Confirm password"
                    type="password"
                    value={field.value}
                  />
                </FieldControl>
                <FieldMessage field={field} />
              </FieldRoot>
            )}
          </Field>
          <Button type="submit" disabled={signUpForm.submitting}>
            <Show when={!signUpForm.submitting} fallback={<>Signing up</>}>
              Sign up
            </Show>
          </Button>
          <FormMessage form={signUpForm} />
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

const actionForgot = action((form: ForgotForm) => useClient()
  .auth.forgotPassword(form)
  .then(() => { toast.success("Sent password reset email.") })
  .catch(throwAsFormError))

export function Forgot() {
  const [forgetForm, { Field, Form }] = createForm<ForgotForm>({ initialValues: { email: "" } });
  const forgotSubmit = useAction(actionForgot)

  return (
    <Layout>
      <Header />
      <CardRoot class="flex flex-col gap-4 p-4">
        <CardHeader>Forgot</CardHeader>
        <Form class="flex flex-col gap-4" onSubmit={(form) => forgotSubmit(form).then(() => reset(forgetForm))}>
          <Field name="email" validate={required('Please enter your email.')}>
            {(field, props) => (
              <FieldRoot class="gap-1.5">
                <FieldLabel field={field}>Email</FieldLabel>
                <FieldControl field={field}>
                  <Input
                    {...props}
                    placeholder="Email"
                    type="email"
                    value={field.value}
                  />
                </FieldControl>
                <FieldMessage field={field} />
              </FieldRoot>
            )}
          </Field>
          <Button type="submit" disabled={forgetForm.submitting}>
            <Show when={!forgetForm.submitting} fallback={<>Sending password reset email</>}>
              Send password reset email
            </Show>
          </Button>
          <FormMessage form={forgetForm} />
        </Form>
      </CardRoot>
      <Footer>
        <A href="/signin" class={linkVariants()}>Sign in</A>
      </Footer>
    </Layout>
  )
}
