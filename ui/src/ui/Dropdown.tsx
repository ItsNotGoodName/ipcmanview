import { keyframes } from "@macaron-css/core";
import { styled } from "@macaron-css/solid";
import {
  Accessor,
  Component,
  createEffect,
  createSignal,
  JSX,
  onCleanup,
  ParentComponent,
  Show,
} from "solid-js";
import { theme } from "./theme";
import { utility } from "./utility";

const TheDropdown = styled("details", {});

export type DropdownProps = {
  children: (props: {
    open: Accessor<boolean>;
    close: () => void;
  }) => JSX.Element;
};

export const Dropdown: Component<DropdownProps> = (props) => {
  const [open, setOpen] = createSignal(false);
  let det: HTMLDetailsElement;
  const close = () => (det.open = false);

  const onClick = (ev: MouseEvent) => {
    if (!det.contains(ev.target as Node)) {
      det.open = false;
    }
  };

  createEffect(() => {
    if (open()) {
      document.addEventListener("click", onClick);
    } else {
      document.removeEventListener("click", onClick);
    }
  });

  onCleanup(() => {
    document.removeEventListener("click", onClick);
  });

  return (
    <TheDropdown
      ref={det!}
      onToggle={() => {
        setOpen(det.open);
      }}
    >
      <props.children open={open} close={close} />
    </TheDropdown>
  );
};

export const DropdownSummary = styled("summary", {
  base: {
    height: "100%",
    cursor: "pointer",
    selectors: {
      ["&::marker"]: {
        content: "",
      },
    },
  },
});

export const DropdownButton = styled("summary", {
  base: {
    height: "100%",
    whiteSpace: "nowrap",
    borderRadius: theme.borderRadius,
    cursor: "pointer",
    selectors: {
      ["&::marker"]: {
        content: "",
      },
    },
  },
  variants: {
    size: {
      small: {
        padding: `${theme.space["0.5"]} ${theme.space[2]}`,
      },
      medium: {
        padding: `${theme.space[2]} ${theme.space[2]}`,
      },
      large: {
        padding: `${theme.space[4]} ${theme.space[4]}`,
      },
    },
    color: {
      primary: {
        background: theme.color.Mauve,
        selectors: {
          ["&:hover"]: {
            background: theme.color.Mauve2,
          },
          [`${TheDropdown}[open] &`]: {
            background: theme.color.Mauve2,
          },
        },
        color: theme.color.Crust,
      },
      secondary: {
        background: theme.color.Subtext0,
        selectors: {
          ["&:hover"]: {
            background: theme.color.Subtext1,
          },
          [`${TheDropdown}[open] &`]: {
            background: theme.color.Subtext1,
          },
        },
        color: theme.color.Crust,
      },
      success: {
        background: theme.color.Green,
        selectors: {
          ["&:hover"]: {
            background: theme.color.Green2,
          },
          [`${TheDropdown}[open] &`]: {
            background: theme.color.Green2,
          },
        },
        color: theme.color.Crust,
      },
      danger: {
        background: theme.color.Red,
        selectors: {
          ["&:hover"]: {
            background: theme.color.Red2,
          },
          [`${TheDropdown}[open] &`]: {
            background: theme.color.Red2,
          },
        },
        color: theme.color.Crust,
      },
    },
  },
  defaultVariants: {
    size: "medium",
    color: "primary",
  },
});

const appearAnimation = keyframes({
  from: { transform: "scale(95%)", opacity: 0 },
  to: { transform: "scale(100%)", opacity: 1 },
});

const TheDropdownEnd = styled("div", {
  base: {
    display: "flex",
    flexDirection: "row-reverse",
  },
});

const TheDropdownContent = styled("div", {
  base: {
    ...utility.shadowXl,
    zIndex: 10,
    position: "absolute",
    width: theme.space[32],
    borderRadius: theme.borderRadius,
    backgroundColor: theme.color.Surface0,
    border: `${theme.space.px} solid ${theme.color.Overlay0}`,
    marginTop: theme.space[1],
    selectors: {
      [`${TheDropdown}[open] &`]: {
        animation: `${appearAnimation} 0.1s`,
      },
      [`${TheDropdownEnd} &`]: {},
    },
  },
});

type DropdownContentProps = {
  end?: boolean;
} & JSX.HTMLAttributes<HTMLDivElement>;

export const DropdownContent: ParentComponent<DropdownContentProps> = (
  props
) => (
  <Show
    when={props.end}
    fallback={
      <TheDropdownContent {...props}>{props.children}</TheDropdownContent>
    }
  >
    <TheDropdownEnd>
      <TheDropdownContent {...props}>{props.children}</TheDropdownContent>
    </TheDropdownEnd>
  </Show>
);

