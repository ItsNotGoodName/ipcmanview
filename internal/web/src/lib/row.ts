import { Accessor, createEffect } from "solid-js"
import { createStore } from "solid-js/store"

export type CreateRowSelectorReturn<T> = {
  selected: Accessor<Array<T>>
  selections: Array<{ id: T, checked: boolean }>
  multiple: Accessor<boolean>
  indeterminate: Accessor<boolean>
  check: (id: T, value: boolean) => void
  checkAll: (value: boolean) => void
}

export function createRowSelector<T>(rows: Accessor<Array<T>>): CreateRowSelectorReturn<T> {
  const [selections, setSelections] = createStore<Array<{ id: T, checked: boolean }>>([])
  createEffect(() =>
    setSelections((prev) => rows().map(v => ({ id: v, checked: prev.find(p => p.id == v)?.checked || false }))))

  const selected = () => selections.filter(v => v.checked == true).map(v => v.id)

  return {
    selected,
    selections,
    multiple: () => {
      const length = selected().length
      return length != 0 && length == selections.length
    },
    indeterminate: () => {
      const length = selected().length
      return length != 0 && length != selections.length
    },
    check: (id, value) => {
      setSelections(
        (todo) => todo.id === id,
        "checked",
        value,
      );
    },
    checkAll: (value) => {
      setSelections(
        () => true,
        "checked",
        value,
      );
    }
  }
}

