import { Accessor, createSignal } from "solid-js"

export type CreateRowSelectorReturn<T> = {
  selected: Accessor<Array<T>>
  selections: Accessor<Array<boolean>>
  multiple: Accessor<boolean>
  indeterminate: Accessor<boolean>
  check: (id: T, value: boolean) => void
  checkAll: (value: boolean) => void
}

// TODO: stop copying arrays
export function createRowSelector<T>(rows: Accessor<Array<T>>): CreateRowSelectorReturn<T> {
  let state: T[] = []

  let msg: { id?: T, value: boolean, all?: boolean } | undefined
  const [checkMsg, setCheckMsg] = createSignal(0)

  // Rows that are selected
  const selected = () => {
    checkMsg()

    // Sync with rows
    state = state.filter(x => rows().includes(x));

    if (msg) {
      if (msg.all) {
        if (msg.value) {
          // Select all
          state = rows()
        } else {
          // Select none
          state = []
        }
      } else if (msg.id) {
        if (msg.value) {
          // Select by id
          state.push(msg.id)
        } else {
          // Remove by id
          state = state.filter(x => x != msg!.id);
        }
      }
    }
    msg = undefined

    return state
  }

  // All rows selections states
  const selections = () => {
    const data = []
    for (let i = 0; i < rows().length; i++) {
      data.push(selected().includes(rows()[i]))
    }
    return data
  }

  return {
    selected,
    selections,
    multiple: () => {
      return selected().length != 0 && selections().length == selected().length
    },
    indeterminate: () => {
      return selected().length != 0 && selections().length != selected().length
    },
    check: (id, value) => {
      msg = { id, value }
      setCheckMsg((prev) => prev + 1)
    },
    checkAll: (value) => {
      msg = { value, all: true }
      setCheckMsg((prev) => prev + 1)
    }
  }
}

