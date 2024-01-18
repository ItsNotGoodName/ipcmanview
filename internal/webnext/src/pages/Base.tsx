import { FormError, createForm, required, reset } from "@modular-forms/solid";
import { RpcError } from "@protobuf-ts/runtime-rpc";
import { A, useNavigate } from "@solidjs/router";
import { ParentProps, Show } from "solid-js";

import { useAuth } from "~/providers/auth";
import { useClient } from "~/providers/client";
import { Button } from "~/ui/Button";
import { CardRoot } from "~/ui/Card";
import { FieldControl, FieldRoot, FieldLabel, FieldMessage, FormMessage } from "~/ui/Form";
import { Input } from "~/ui/Input";
import { linkVariants } from "~/ui/Link";
import { ThemeIcon } from "~/ui/ThemeIcon";
import { toast } from "~/ui/Toast";
import { toggleTheme, useThemeTitle } from "~/ui/theme";

function Header() {
  return (
    <div class="text-2xl text-center">IPCManView</div>
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
      <div class="flex items-center justify-between gap-2">
        <a href="/" class={linkVariants()}>Management</a>
        {props.children}
      </div>
    </CardRoot>
  )
}

type SigninForm = {
  usernameOrEmail: string
  password: string
}

function createSigninSubmit() {
  const auth = useAuth()
  const client = useClient()
  return async (form: SigninForm) => {
    try {
      const resp = await client.auth.signIn(form)
      auth.setSession(resp.response.token)
    } catch (e) {
      if (e instanceof RpcError)
        // @ts-ignore
        throw new FormError(e.message, e.meta ?? {})
      if (e instanceof Error)
        throw new FormError(e.message)
      throw new FormError("Unknown error has occured.")
    }
  }
}

export function Signin() {
  const [signinForm, { Field, Form }] = createForm<SigninForm>();
  const signinSubmit = createSigninSubmit()

  return (
    <div class="mx-auto flex max-w-xs flex-col gap-4 pt-10">
      <Header />
      <CardRoot class="flex flex-col gap-4 p-4">
        <CardHeader>Sign in</CardHeader>
        <Form class="flex flex-col gap-4" onSubmit={signinSubmit}>
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
                    required
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
                    value={field.value}
                    type="password"
                    required
                  />
                </FieldControl>
                <FieldMessage field={field} />
              </FieldRoot>
            )}
          </Field>
          <Button type="submit" disabled={signinForm.submitting}>
            <Show when={signinForm.submitting} fallback={<>Sign in</>} >
              Signing in
            </Show>
          </Button>
          <FormMessage form={signinForm} />
        </Form>
      </CardRoot>
      <Footer>
        <A href="/signup" class={linkVariants()}>Sign up</A>
      </Footer >
    </div>
  )
}

type SignupForm = {
  email: string
  username: string
  password: string
  confirmPassword: string
}

function createSignupSubmit() {
  const navigate = useNavigate()
  const client = useClient()
  return async (form: SignupForm) => {
    if (form.password != form.confirmPassword) {
      throw new FormError<SignupForm>("", { confirmPassword: "Password does not match." })
    }

    try {
      await client.auth.signUp(form)
    } catch (e) {
      if (e instanceof RpcError)
        // @ts-ignore
        throw new FormError(e.message, e.meta ?? {})
      if (e instanceof Error)
        throw new FormError(e.message)
      throw new FormError("Unknown error has occured.")
    }

    navigate('/signin')
  }
}

export function Signup() {
  const [signupForm, { Field, Form }] = createForm<SignupForm>();
  const signupSubmit = createSignupSubmit()

  return (
    <div class="mx-auto flex max-w-xs flex-col gap-4 pt-10">
      <Header />
      <CardRoot class="flex flex-col gap-4 p-4">
        <CardHeader>Sign up</CardHeader>
        <Form class="flex flex-col gap-4" onSubmit={signupSubmit}>
          <Field name="email" validate={required('Please enter your email.')}>
            {(field, props) => (
              <FieldRoot class="gap-1.5">
                <FieldLabel field={field}>Email</FieldLabel>
                <FieldControl field={field}>
                  <Input
                    {...props}
                    type="email"
                    placeholder="Email"
                    value={field.value}
                    required
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
                    required
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
                    value={field.value}
                    type="password"
                    required
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
                    value={field.value}
                    type="password"
                    required
                  />
                </FieldControl>
                <FieldMessage field={field} />
              </FieldRoot>
            )}
          </Field>
          <Button type="submit" disabled={signupForm.submitting}>
            <Show when={signupForm.submitting} fallback={<>Sign up</>} >
              Signing up
            </Show>
          </Button>
          <FormMessage form={signupForm} />
        </Form>
      </CardRoot>
      <Footer>
        <A href="/signin" class={linkVariants()}>Sign in</A>
      </Footer >
    </div >
  )
}

type ForgotForm = {
  email: string
}

function createForgotSubmit() {
  const client = useClient()
  return async (form: ForgotForm) => {
    try {
      await client.auth.resetPassword(form)
    } catch (e) {
      if (e instanceof RpcError)
        // @ts-ignore
        throw new FormError(e.message, e.meta ?? {})
      if (e instanceof Error)
        throw new FormError(e.message)
      throw new FormError("Unknown error has occured.")
    }

    toast.success("Sent password reset email.")
  }
}

export function Forgot() {
  const [forgetForm, { Field, Form }] = createForm<ForgotForm>({ initialValues: { email: "" } });
  const forgotSubmit = createForgotSubmit()


  return (
    <div class="mx-auto flex max-w-xs flex-col gap-4 pt-10">
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
                    type="email"
                    placeholder="Email"
                    value={field.value}
                  />
                </FieldControl>
                <FieldMessage field={field} />
              </FieldRoot>
            )}
          </Field>
          <Button type="submit" disabled={forgetForm.submitting}>
            <Show when={forgetForm.submitting} fallback={<>Send password reset email</>} >
              Sending password reset email
            </Show>
          </Button>
          <FormMessage form={forgetForm} />
        </Form>
      </CardRoot>
      <Footer>
        <A href="/signin" class={linkVariants()}>Sign in</A>
      </Footer>
    </div>
  )
}
