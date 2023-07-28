import { Component } from "solid-js";
import { createForm, required, ResponseData } from "@modular-forms/solid";
import { styled } from "@macaron-css/solid";
import { style } from "@macaron-css/core";
import { Button } from "~/ui/Button";
import { InputText } from "~/ui/InputText";
import { Card, CardBody, CardHeader, CardHeaderTitle } from "~/ui/Card";
import { LayoutCenter } from "~/ui/Layout";
import { ThemeSwitcher, ThemeSwitcherIcon } from "~/ui/ThemeSwitcher";
import { theme } from "~/ui/theme";
import { utility } from "~/ui/utility";

const Center = styled("div", {
  base: {
    display: "flex",
    justifyContent: "center",
  },
});

const Right = styled("div", {
  base: {
    display: "flex",
    justifyContent: "end",
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

type LoginMutation = {
  usernameOrEmail: string;
  password: string;
};

export const Login: Component = () => {
  const [form, { Form, Field }] = createForm<LoginMutation, ResponseData>({});

  return (
    <LayoutCenter>
      <Card>
        <CardHeader>
          <CardHeaderTitle>IPCMango</CardHeaderTitle>
          <ThemeSwitcher class={themeSwitcherClass}>
            <ThemeSwitcherIcon class={style({ ...utility.size("6") })} />
          </ThemeSwitcher>
        </CardHeader>
        <CardBody>
          <Form onSubmit={() => alert("Submitted")}>
            <Stack>
              <Field
                name="usernameOrEmail"
                validate={[required("Please enter your username or email.")]}
              >
                {(field, props) => (
                  <InputText
                    {...props}
                    label="Username or Email"
                    placeholder="Username or Email"
                    autocomplete="username"
                    disabled={form.submitting}
                    error={field.error}
                  />
                )}
              </Field>

              <Field
                name="password"
                validate={[required("Please enter your password.")]}
              >
                {(field, props) => (
                  <InputText
                    {...props}
                    label="Password"
                    type="password"
                    placeholder="Password"
                    autocomplete="current-password"
                    disabled={form.submitting}
                    error={field.error}
                  />
                )}
              </Field>

              <Right>
                <a href="#">Forgot Password?</a>
              </Right>

              <Button type="submit" disabled={form.submitting}>
                Log in
              </Button>
            </Stack>
          </Form>
        </CardBody>
      </Card>
      <Center>
        <a href="#">Admin Panel</a>
      </Center>
    </LayoutCenter>
  );
};

