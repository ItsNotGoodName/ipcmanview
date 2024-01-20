import { JSX, Show, createUniqueId, splitProps, createContext, useContext, } from "solid-js";
import { FieldStore, FormStore } from "@modular-forms/solid";

import { cn } from "~/lib/utils"
import { Label, LabelProps } from "./Label"

type FieldContextValue = {
  id: string
}

const FormItemContext = createContext<FieldContextValue>(
  {} as FieldContextValue
)

export function FieldRoot(props: Omit<JSX.HTMLAttributes<HTMLDivElement>, "id">) {
  const [_, rest] = splitProps(props, ["class"])
  const id = createUniqueId()

  return (
    <FormItemContext.Provider value={{ id }}>
      <div class={cn("space-y-2", props.class)} {...rest} />
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

export function FieldControl(props: JSX.HTMLAttributes<HTMLDivElement> & { field: FieldStore<any, any> }) {
  const [_, rest] = splitProps(props, ["field"])
  const { formFieldId, formDescriptionId, formMessageId } = useField()

  return (
    <div
      id={formFieldId}
      aria-describedby={
        !props.field.error
          ? `${formDescriptionId}`
          : `${formDescriptionId} ${formMessageId}`
      }
      aria-invalid={!!props.field.error}
      {...rest}
    />
  )
}

export function FieldDescription(props: JSX.HTMLAttributes<HTMLParagraphElement>) {
  const [_, rest] = splitProps(props, ["class"])
  const { formDescriptionId } = useField()

  return (
    <p
      id={formDescriptionId}
      class={cn("text-sm text-muted-foreground", props.class)}
      {...rest}
    />
  )
}

export function FieldMessage(props: JSX.HTMLAttributes<HTMLParagraphElement> & { field: FieldStore<any, any> }) {
  const [_, rest] = splitProps(props, ["class", "field", "children"])
  const { formMessageId } = useField()
  const body = () => props.field.error ? props.field.error : props.children

  return (
    <Show when={body()}>
      <p
        id={formMessageId}
        class={cn("text-sm font-medium text-destructive", props.class)}
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
        class={cn("text-sm font-medium text-destructive", props.class)}
        {...rest}
      >
        {body()}
      </p>
    </Show>
  )
}
