import { Component } from "solid-js";
import {
  A,
  Route,
  Routes,
  useLocation,
} from "@solidjs/router";
import {
  RiBuildingsHome5Line,
  RiDeviceServerLine,
  RiUserFacesAdminLine,
} from "solid-icons/ri";
import { styled } from "@macaron-css/solid";
import { CSSProperties, style } from "@macaron-css/core";

import { minScreen, theme } from "~/ui/theme";
import { Menu, menuChildClass } from "~/ui/Menu";

const Root = styled("div", {
  base: {
    display: "flex",
    flexDirection: "column",
    height: "100%",
  },
});

const Header = styled("div", {
  base: {
    overflowX: "auto",
    overflowY: "hidden",
    display: "flex",
    height: theme.space[11],
    gap: theme.space[2],
    justifyContent: "space-between",
    padding: theme.space[2],
    background: theme.color.Crust,
    borderBottom: `${theme.space.px} solid ${theme.color.Overlay0}`,
  },
});

const HeaderText = styled("div", {
  base: {
    display: "flex",
    alignItems: "center",
    overflow: "hidden",
  },
});

const HeaderTextContent = styled("div", {
  base: {
    ...utility.textLine(),
    fontSize: "x-large",
  },
});

const HeaderEnd = styled("div", {
  base: {
    display: "flex",
    gap: theme.space[2],
  },
});

const Content = styled("div", {
  base: {
    display: "flex",
    flexDirection: "column",
    height: "100%",
    overflow: "hidden",
    "@media": {
      [minScreen.md]: {
        flexDirection: "row",
      },
    },
  },
});

const ContentNav = styled("div", {
  base: {
    display: "flex",
    gap: theme.space[1],
    padding: theme.space[2],
    background: theme.color.Mantle,
    justifyContent: "space-between",
    height: theme.space[11],
    borderBottom: `${theme.space.px} solid ${theme.color.Overlay0}`,
    overflowX: "auto",
    "@media": {
      [minScreen.md]: {
        width: theme.space[11],
        flexDirection: "column",
        height: "100%",
        borderRight: `${theme.space.px} solid ${theme.color.Overlay0}`,
        borderBottom: "none",
      },
    },
  },
});

const ContentNavStart = styled("div", {
  base: {
    display: "flex",
    gap: "inherit",
    flexDirection: "inherit",
  },
});

const ContentBody = styled("div", {
  base: {
    width: "100%",
    height: "100%",
    overflow: "auto",
  },
});

const chipClass = style({
  display: "flex",
  alignItems: "center",
  justifyContent: "center",
  color: theme.color.Text,
  ":hover": {
    color: theme.color.Mauve,
    borderColor: theme.color.Mauve,
    fill: theme.color.Mauve,
  },
});

const activeChipClass = style({
  color: theme.color.Mauve,
  borderColor: theme.color.Mauve,
  fill: theme.color.Mauve,
});

const chipChildStyle = {
  height: theme.space[7],
  width: theme.space[7],
} as CSSProperties;

const iconClass = style({
  ...chipChildStyle,
});

const avatarClass = style({
  ...chipChildStyle,
  border: `${theme.space.px} solid`,
  borderRadius: "100%",
  backgroundColor: theme.color.Surface0,
  textTransform: "uppercase",
  padding: theme.space["0.5"],
  display: "flex",
  justifyContent: "center",
  alignItems: "center",
  userSelect: "none",
});

import { Dropdown, DropdownSummary, DropdownContent } from "~/ui/Dropdown";
import { ThemeSwitcher, ThemeSwitcherIcon } from "~/ui/ThemeSwitcher";
import { utility } from "~/ui/utility";

import { Home } from "~/pages/Home.tsx"

export const Application: Component = () => {
  return (
    <Root>
      <div>
        <Header>
          <HeaderText>
            <HeaderTextContent>IPCMango</HeaderTextContent>
          </HeaderText>
          <HeaderEnd>
            <ThemeSwitcher class={chipClass}>
              <ThemeSwitcherIcon class={iconClass} />
            </ThemeSwitcher>
            <Dropdown>
              {(props) => {
                const location = useLocation();
                const active = () =>
                  props.open() || location.pathname == "/profile";

                return (
                  <>
                    <DropdownSummary
                      class={chipClass}
                      classList={{ [activeChipClass]: active() }}
                      title="User"
                    >
                      <div class={avatarClass}>
                        GD
                      </div>
                    </DropdownSummary>
                    <DropdownContent end={true}>
                      <Menu>
                        <A
                          href="/profile"
                          onclick={props.close}
                          class={menuChildClass}
                        >
                          Profile
                        </A>
                        <button class={menuChildClass}>
                          Log out
                        </button>
                      </Menu>
                    </DropdownContent>
                  </>
                );
              }}
            </Dropdown>
          </HeaderEnd>
        </Header>
      </div>
      <Content>
        <div>
          <ContentNav>
            <ContentNavStart>
              <A
                href="/"
                title="Home"
                end
                class={chipClass}
                activeClass={activeChipClass}
              >
                <RiBuildingsHome5Line class={iconClass} />
              </A>
              <A
                href="/stations"
                title="Stations"
                class={chipClass}
                activeClass={activeChipClass}
              >
                <RiDeviceServerLine class={iconClass} />
              </A>
            </ContentNavStart>
            <a href="#" title="Admin Panel" class={chipClass}>
              <RiUserFacesAdminLine class={iconClass} />
            </a>
          </ContentNav>
        </div>
        <ContentBody>
          <Routes>
            <Route path="/" component={Home} />
          </Routes>
        </ContentBody>
      </Content>
    </Root>
  );
};

