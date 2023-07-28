import { Component } from "solid-js";
import { createForm, FormError, required, ResponseData, SubmitHandler } from "@modular-forms/solid";
import { styled } from "@macaron-css/solid";
import { style } from "@macaron-css/core";
import { Button } from "~/ui/Button";
import { InputText } from "~/ui/InputText";
import { Card, CardBody, CardHeader, CardHeaderTitle } from "~/ui/Card";
import { LayoutCenter } from "~/ui/Layout";
import { ThemeSwitcher, ThemeSwitcherIcon } from "~/ui/ThemeSwitcher";
import { theme } from "~/ui/theme";
import { utility } from "~/ui/utility";
import { useAuthStore } from "~/providers/auth";
import { ErrorText } from "~/ui/ErrorText";
import { UserRegister } from "~/core/client.gen";
import { A, useNavigate } from "@solidjs/router";

const Center = styled("div", {
  base: {
    display: "flex",
    justifyContent: "center",
  },
});

const Stack = styled("div", {
  base: {
    ...utility.stack("4"),
  },
});

const themeSwitcherClass = style({
  display: "flex",
  alignItems: "center",
  borderRadius: theme.borderRadius,
  ":hover": {
    backgroundColor: theme.color.Surface2,
  },
});

type RegisterMutation = {
  [Property in keyof UserRegister]: Property;
};

export const Register: Component = () => {
  const [form, { Form, Field }] = createForm<RegisterMutation, ResponseData>({});
  const navigate = useNavigate();
  const auth = useAuthStore()

  const submit: SubmitHandler<RegisterMutation> = (values) => auth.register({ user: values }).then(() => {
    navigate("/login")
  }).catch((e: Error) => {
    throw new FormError<RegisterMutation>(e.message);
  });

  return (
    <LayoutCenter>
      <Card>
        <CardHeader>
          <CardHeaderTitle>Register</CardHeaderTitle>
          <ThemeSwitcher class={themeSwitcherClass}>
            <ThemeSwitcherIcon class={style({ ...utility.size("6") })} />
          </ThemeSwitcher>
        </CardHeader>
        <CardBody>
          <Form onSubmit={submit}>
            <Stack>
              <Field
                name="email"
                validate={[required("Please enter your email.")]}
              >
                {(field, props) => (
                  <InputText
                    {...props}
                    label="Email"
                    placeholder="Email"
                    disabled={form.submitting}
                    error={field.error}
                  />
                )}
              </Field>

              <Field
                name="username"
                validate={[required("Please enter a username.")]}
              >
                {(field, props) => (
                  <InputText
                    {...props}
                    label="Username"
                    placeholder="Username"
                    disabled={form.submitting}
                    error={field.error}
                  />
                )}
              </Field>

              <Field
                name="password"
                validate={[required("Please enter a password.")]}
              >
                {(field, props) => (
                  <InputText
                    {...props}
                    label="Password"
                    type="password"
                    placeholder="Password"
                    disabled={form.submitting}
                    error={field.error}
                  />
                )}
              </Field>

              <Field
                name="passwordConfirm"
                validate={[required("Please confirm password.")]}
              >
                {(field, props) => (
                  <InputText
                    {...props}
                    label="Password confirm"
                    type="password"
                    placeholder="Password confirm"
                    disabled={form.submitting}
                    error={field.error}
                  />
                )}
              </Field>

              <Button type="submit" disabled={form.submitting}>
                Register
              </Button>

              <ErrorText>{form.response.message}</ErrorText>
            </Stack>
          </Form>
        </CardBody>
      </Card>
      <Center>
        <A href="/login">Login</A>
      </Center>
    </LayoutCenter>
  );
};

