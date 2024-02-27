// # Changes
// - Replace React Hook Form with Modular Forms
//
// # URLs
// https://ui.shadcn.com/docs/components/form
import { JSX, Show, createUniqueId, splitProps, createContext, useContext, ParentProps, } from "solid-js";
import { FieldStore, FormStore, setValue } from "@modular-forms/solid";
import { Checkbox } from "@kobalte/core";

import { cn } from "~/lib/utils"
import { Label, LabelProps } from "./Label"

type FieldContextValue = {
  id: string
}

const FormItemContext = createContext<FieldContextValue>(
  {} as FieldContextValue
)

export function CheckboxFieldRoot(props: Checkbox.CheckboxRootProps & { field: FieldStore<any, any>, form: FormStore<any, any> }) {
  const [_, rest] = splitProps(props, ["field", "form"])
  return <Checkbox.Root
    validationState={props.field.error ? "invalid" : "valid"}
    checked={props.field.value}
    onChange={(value) => setValue(props.form, props.field.name, value)}
    {...rest}
  />
}

export function FieldRoot2(props: ParentProps) {
  const id = createUniqueId()

  return (
    <FormItemContext.Provider value={{ id }}>
      {props.children}
    </FormItemContext.Provider>
  )
}

export function FieldRoot(props: Omit<JSX.HTMLAttributes<HTMLDivElement>, "id">) {
  const [_, rest] = splitProps(props, ["class"])
  const id = createUniqueId()

  return (
    <FormItemContext.Provider value={{ id }}>
      <div {...rest} />
    </FormItemContext.Provider>
  )
}

function useField() {
  const itemContext = useContext(FormItemContext)
  if (!itemContext) throw new Error("useField should be used within <FieldRoot>");

  const { id } = itemContext

  return {
    id,
    formFieldId: `${id}-form-field`,
    formDescriptionId: `${id}-form-field-description`,
    formMessageId: `${id}-form-field-message`,
  }
}

export function FieldLabel(props: Omit<LabelProps, "for"> & { field: FieldStore<any, any> }) {
  const [_, rest] = splitProps(props, ["class", "field"])
  const { formFieldId } = useField()

  return (
    <Label
      class={cn(props.field.error && "text-destructive", props.class)}
      for={formFieldId}
      {...rest}
    />
  )
}

export function fieldControlProps(field: FieldStore<any, any>) {
  const { formFieldId, formDescriptionId, formMessageId } = useField()
  return {
    id: formFieldId,
    "aria-describedby": !field.error ? `${formDescriptionId}` : `${formDescriptionId} ${formMessageId}`,
    "aria-invalid": !!field.error
  }
}

export function FieldDescription(props: JSX.HTMLAttributes<HTMLParagraphElement>) {
  const [_, rest] = splitProps(props, ["class"])
  const { formDescriptionId } = useField()

  return (
    <p
      id={formDescriptionId}
      class={cn("text-muted-foreground text-sm", props.class)}
      {...rest}
    />
  )
}

export function FieldErrorMessage(props: JSX.HTMLAttributes<HTMLParagraphElement> & { field: FieldStore<any, any> }) {
  const [_, rest] = splitProps(props, ["class", "field", "children"])
  const { formMessageId } = useField()
  const body = () => props.field.error ? props.field.error : props.children

  return (
    <Show when={body()}>
      <p
        id={formMessageId}
        class={cn("text-destructive text-sm font-medium", props.class)}
        {...rest}
      >
        {body()}
      </p>
    </Show>
  )
}

export function FormMessage(props: JSX.HTMLAttributes<HTMLParagraphElement> & { form: FormStore<any, any> }) {
  const [_, rest] = splitProps(props, ["class", "form", "children"])
  const body = () => props.form.response.message ? props.form.response.message : props.children

  return (
    <Show when={body()}>
      <p
        class={cn("text-destructive text-sm font-medium", props.class)}
        {...rest}
      >
        {body()}
      </p>
    </Show>
  )
}
