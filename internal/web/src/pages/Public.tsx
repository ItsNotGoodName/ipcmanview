import { createForm, required, reset } from "@modular-forms/solid";
import { A, createAsync, revalidate, useNavigate } from "@solidjs/router";
import { ParentProps, Show, } from "solid-js";
import { useClient } from "~/providers/client";
import { Button } from "~/ui/Button";
import { CardRoot } from "~/ui/Card";
import { FormMessage, } from "~/ui/Form";
import { linkVariants } from "~/ui/Link";
import { ThemeIcon } from "~/ui/ThemeIcon";
import { toggleTheme, useThemeTitle } from "~/ui/theme";
import { CheckboxControl, CheckboxErrorMessage, CheckboxLabel, CheckboxRoot } from "~/ui/Checkbox";
import { setFormValue, throwAsFormError, validationState } from "~/lib/utils";
import { toast } from "~/ui/Toast";
import { getSession } from "~/providers/session";
import { AlertDescription, AlertRoot, AlertTitle } from "~/ui/Alert";
import { getConfig } from "./data";
import { TextFieldErrorMessage, TextFieldInput, TextFieldLabel, TextFieldRoot } from "~/ui/TextField";

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

  const [form, { Field, Form }] = createForm<SignInForm>({
    initialValues: {
      usernameOrEmail: "",
      password: "",
      rememberMe: false
    }
  });
  const submitForm = (input: SignInForm) =>
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
      <Show when={session()?.valid && session()?.disabled}>
        <AlertRoot variant="destructive">
          <AlertTitle>Account disabled</AlertTitle>
          <AlertDescription>
            Your account "{session()?.username}" is disabled.
          </AlertDescription>
        </AlertRoot>
      </Show>
      <CardRoot class="flex flex-col gap-4 p-4">
        <CardHeader>Sign in</CardHeader>
        <Form onSubmit={submitForm} class="flex flex-col gap-4">
          <Field name="usernameOrEmail" validate={required("Please enter your username or email.")}>
            {(field, props) => (
              <TextFieldRoot
                validationState={validationState(field.error)}
                value={field.value}
                class="space-y-2"
              >
                <TextFieldLabel>Username or email</TextFieldLabel>
                <TextFieldInput
                  {...props}
                  placeholder="Username or email"
                  autocomplete="username"
                />
                <TextFieldErrorMessage>{field.error}</TextFieldErrorMessage>
              </TextFieldRoot>
            )}
          </Field>
          <Field name="password">
            {(field, props) => (
              <TextFieldRoot
                validationState={validationState(field.error)}
                value={field.value}
                class="space-y-2"
              >
                <div class="flex items-center justify-between gap-2">
                  <TextFieldLabel>Password</TextFieldLabel>
                  <A href="/forgot" class={linkVariants()}>Forgot password?</A>
                </div>
                <TextFieldInput
                  {...props}
                  autocomplete="current-password"
                  placeholder="Password"
                  type="password"
                />
                <TextFieldErrorMessage>{field.error}</TextFieldErrorMessage>
              </TextFieldRoot>
            )}
          </Field>
          <Field name="rememberMe" type="boolean">
            {(field) => (
              <CheckboxRoot
                validationState={validationState(field.error)}
                checked={field.value}
                onChange={setFormValue(form, field)}
                class="space-y-2"
              >
                <div class="flex items-center gap-2">
                  <CheckboxControl />
                  <CheckboxLabel>Remember me</CheckboxLabel>
                </div>
                <CheckboxErrorMessage>{field.error}</CheckboxErrorMessage>
              </CheckboxRoot>
            )}
          </Field>
          <Button type="submit" disabled={form.submitting}>
            <Show when={!form.submitting} fallback="Signing in">Sign in</Show>
          </Button>
          <FormMessage form={form} />
        </Form>
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
    initialValues: {
      email: "",
      username: "",
      password: "",
      confirmPassword: "",
    },
    validate: (input) => {
      if (input.password != input.confirmPassword) {
        return {
          confirmPassword: "Password does not match."
        }
      }
      return {}
    }
  });
  const submitForm = (input: SignUpForm) => useClient()
    .public.signUp(input)
    .then()
    .catch(throwAsFormError)
    .then(() => navigate("/signin"))

  return (
    <Layout>
      <Header>{config()?.siteName}</Header>
      <CardRoot class="flex flex-col gap-4 p-4">
        <CardHeader>Sign up</CardHeader>
        <Form onSubmit={submitForm} class="flex flex-col gap-4">
          <Field name="email" validate={required('Please enter your email.')}>
            {(field, props) => (
              <TextFieldRoot
                validationState={validationState(field.error)}
                value={field.value}
                class="space-y-2"
              >
                <TextFieldLabel>Email</TextFieldLabel>
                <TextFieldInput
                  {...props}
                  placeholder="Email"
                  type="email"
                  value={field.value}
                />
                <TextFieldErrorMessage>{field.error}</TextFieldErrorMessage>
              </TextFieldRoot>
            )}
          </Field>
          <Field name="username" validate={required('Please enter a username.')}>
            {(field, props) => (
              <TextFieldRoot
                validationState={validationState(field.error)}
                value={field.value}
                class="space-y-2"
              >
                <TextFieldLabel>Username</TextFieldLabel>
                <TextFieldInput
                  {...props}
                  autocomplete="username"
                  placeholder="Username"
                />
                <TextFieldErrorMessage>{field.error}</TextFieldErrorMessage>
              </TextFieldRoot>
            )}
          </Field>
          <Field name="password" validate={required('Please enter a password.')}>
            {(field, props) => (
              <TextFieldRoot
                validationState={validationState(field.error)}
                value={field.value}
                class="space-y-2"
              >
                <TextFieldLabel>Password</TextFieldLabel>
                <TextFieldInput
                  {...props}
                  autocomplete="new-password"
                  placeholder="Password"
                  type="password"
                />
                <TextFieldErrorMessage>{field.error}</TextFieldErrorMessage>
              </TextFieldRoot>
            )}
          </Field>
          <Field name="confirmPassword" validate={required('Please confirm your password.')}>
            {(field, props) => (
              <TextFieldRoot
                validationState={validationState(field.error)}
                value={field.value}
                class="space-y-2"
              >
                <TextFieldLabel>Confirm password</TextFieldLabel>
                <TextFieldInput
                  {...props}
                  autocomplete="new-password"
                  placeholder="Confirm password"
                  type="password"
                />
                <TextFieldErrorMessage>{field.error}</TextFieldErrorMessage>
              </TextFieldRoot>
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

  const [form, { Field, Form }] = createForm<ForgotForm>({
    initialValues: {
      email: "",
    }
  });
  const submitForm = (input: ForgotForm) => useClient()
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
        <Form onSubmit={submitForm} class="flex flex-col gap-4">
          <Field name="email" validate={required('Please enter your email.')}>
            {(field, props) => (
              <TextFieldRoot
                validationState={validationState(field.error)}
                value={field.value}
                class="space-y-2"
              >
                <TextFieldLabel>Email</TextFieldLabel>
                <TextFieldInput
                  {...props}
                  placeholder="Email"
                  type="email"
                />
                <TextFieldErrorMessage>{field.error}</TextFieldErrorMessage>
              </TextFieldRoot>
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

