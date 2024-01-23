import { As } from "@kobalte/core";
import { AlertDescription, AlertRoot, AlertTitle } from "~/ui/Alert";
import { Button } from "~/ui/Button";
import { DropdownMenuArrow, DropdownMenuCheckboxItem, DropdownMenuCheckboxItemIndicator, DropdownMenuContent, DropdownMenuGroup, DropdownMenuGroupLabel, DropdownMenuItem, DropdownMenuPortal, DropdownMenuRadioGroup, DropdownMenuRadioItem, DropdownMenuRadioItemIndicator, DropdownMenuRoot, DropdownMenuSeparator, DropdownMenuShortcut, DropdownMenuSub, DropdownMenuSubContent, DropdownMenuSubTrigger, DropdownMenuSubTriggerIndicator, DropdownMenuTrigger } from "~/ui/DropdownMenu";
import { Input } from "~/ui/Input";
import { Seperator } from "~/ui/Seperator";
import { Textarea } from "~/ui/Textarea";
import { Label } from "~/ui/Label";
import { SwitchControl, SwitchDescription, SwitchErrorMessage, SwitchInput, SwitchLabel, SwitchRoot } from "~/ui/Switch";
import { toggleTheme } from "~/ui/theme";
import { CardRoot, CardContent, CardDescription, CardFooter, CardHeader, CardTitle } from "~/ui/Card";
import { For, createEffect, createSignal, onCleanup, onMount } from "solid-js";
import { Badge } from "~/ui/Badge";
import { CheckboxControl, CheckboxDescription, CheckboxErrorMessage, CheckboxInput, CheckboxLabel, CheckboxRoot } from "~/ui/Checkbox";
import { PopoverArrow, PopoverCloseButton, PopoverCloseIcon, PopoverContent, PopoverDescription, PopoverPortal, PopoverRoot, PopoverTitle, PopoverTrigger } from "~/ui/Popover";
import { DialogCloseButton, DialogContent, DialogDescription, DialogFooter, DialogHeader, DialogOverlay, DialogPortal, DialogRoot, DialogTitle, DialogTrigger } from "~/ui/Dialog";
import { TableBody, TableCaption, TableCell, TableHead, TableHeader, TableRoot, TableRow } from "~/ui/Table";
import { ToastCloseButton, ToastContent, ToastDescription, ToastProgressFill, ToastProgressTrack, ToastTitle, toast } from "~/ui/Toast";
import { Skeleton } from "~/ui/Skeleton";
import { ThemeIcon } from "~/ui/ThemeIcon";
import { PaginationContent, PaginationEllipsis, PaginationItem, PaginationItems, PaginationLink, PaginationNext, PaginationPrevious, PaginationRoot } from "~/ui/Pagination";
import { SelectContent, SelectItem, SelectListbox, SelectRoot, SelectTrigger, SelectValue } from "~/ui/Select";
import { ConfirmButton } from "~/ui/Confirm";
import { MenubarCheckboxItem, MenubarContent, MenubarGroup, MenubarGroupLabel, MenubarItem, MenubarMenu, MenubarRadioGroup, MenubarRadioItem, MenubarRoot, MenubarSeparator, MenubarShortcut, MenubarSub, MenubarSubContent, MenubarSubTrigger, MenubarTrigger } from "~/ui/Menubar";
import { TabsContent, TabsList, TabsRoot, TabsTrigger } from "~/ui/Tabs";
import { RiMapRocketLine } from "solid-icons/ri";
import { AvatarFallback, AvatarImage, AvatarRoot } from "~/ui/Avatar";
import { ProgressFill, ProgressLabel, ProgressRoot, ProgressTrack, ProgressValueLabel } from "~/ui/Progress";

export function Ui() {
  const showToast = () => {
    toast.custom(() =>
      <ToastContent>
        <ToastCloseButton />
        <ToastTitle>Title</ToastTitle>
        <ToastDescription>Description</ToastDescription>
        <ToastProgressTrack>
          <ToastProgressFill />
        </ToastProgressTrack>
      </ToastContent>
    )
    toast.show("Hello World")
  }

  const [progress, setProgress] = createSignal(0)
  const timer = setInterval(() => setProgress((prev) => (prev + 10) % 100), 200)
  onCleanup(() => clearInterval(timer))

  return (
    <div class="flex flex-col gap-4 p-4">
      <Button onClick={toggleTheme} size="icon">
        <ThemeIcon class="h-6 w-6" />
      </Button>
      <AlertRoot>
        <RiMapRocketLine class="h-4 w-4" />
        <AlertTitle>Alert Title</AlertTitle>
        <AlertDescription>Alert Description</AlertDescription>
      </AlertRoot>
      <div class="flex flex-col gap-1.5">
        <Label for="input">Label</Label>
        <Input id="input" type="text" placeholder="Input" />
      </div>
      <Textarea placeholder="Textarea"></Textarea>
      <div>
        <div>Top Seperator</div>
        <Seperator />
        <div>Bottom Seperator</div>
      </div>
      <div class="flex justify-between">
        <div>Left Seperator</div>
        <Seperator orientation="vertical" />
        <div>Right Seperator</div>
      </div>
      <SwitchRoot class="flex gap-2">
        <SwitchLabel>Switch</SwitchLabel>
        <SwitchDescription />
        <SwitchErrorMessage />
        <SwitchInput />
        <SwitchControl />
      </SwitchRoot>
      <DropdownMenuRoot>
        <DropdownMenuTrigger asChild>
          <As component={Button}>
            DropdownMenu
          </As>
        </DropdownMenuTrigger>
        <DropdownMenuPortal>
          <DropdownMenuContent>
            <DropdownMenuItem>
              Commit <DropdownMenuShortcut>⌘+K</DropdownMenuShortcut>
            </DropdownMenuItem>
            <DropdownMenuItem>
              Push <DropdownMenuShortcut>⇧+⌘+K</DropdownMenuShortcut>
            </DropdownMenuItem>
            <DropdownMenuItem disabled>
              Update Project <DropdownMenuShortcut>⌘+T</DropdownMenuShortcut>
            </DropdownMenuItem>
            <DropdownMenuSub overlap gutter={4} shift={-8}>
              <DropdownMenuSubTrigger>
                GitHub
                <DropdownMenuSubTriggerIndicator />
              </DropdownMenuSubTrigger>
              <DropdownMenuPortal>
                <DropdownMenuSubContent>
                  <DropdownMenuItem>
                    Create Pull Request…
                  </DropdownMenuItem>
                  <DropdownMenuItem>
                    View Pull Requests
                  </DropdownMenuItem>
                  <DropdownMenuItem>
                    Sync Fork
                  </DropdownMenuItem>
                  <DropdownMenuSeparator />
                  <DropdownMenuItem>
                    Open on GitHub
                  </DropdownMenuItem>
                </DropdownMenuSubContent>
              </DropdownMenuPortal>
            </DropdownMenuSub>
            <DropdownMenuSeparator />
            <DropdownMenuCheckboxItem>
              <DropdownMenuCheckboxItemIndicator />
              Show Git Log
            </DropdownMenuCheckboxItem>
            <DropdownMenuCheckboxItem>
              <DropdownMenuCheckboxItemIndicator />
              Show History
            </DropdownMenuCheckboxItem>
            <DropdownMenuSeparator />
            <DropdownMenuGroup>
              <DropdownMenuGroupLabel>
                Branches
              </DropdownMenuGroupLabel>
              <DropdownMenuRadioGroup>
                <DropdownMenuRadioItem value="main">
                  <DropdownMenuRadioItemIndicator />
                  main
                </DropdownMenuRadioItem>
                <DropdownMenuRadioItem value="develop">
                  <DropdownMenuRadioItemIndicator />
                  develop
                </DropdownMenuRadioItem>
              </DropdownMenuRadioGroup>
            </DropdownMenuGroup>
            <DropdownMenuArrow />
          </DropdownMenuContent>
        </DropdownMenuPortal>
      </DropdownMenuRoot>
      <CardRoot>
        <CardHeader>
          <CardTitle>Card Title</CardTitle>
          <CardDescription>Card Description</CardDescription>
        </CardHeader>
        <CardContent>
          Card Content
        </CardContent>
        <CardFooter>Card Footer</CardFooter>
      </CardRoot>
      <div class="flex gap-4">
        <For each={["default", "secondary", "destructive", "outline"]}>
          {variant =>
            <Badge variant={variant as any}>{variant}</Badge>
          }
        </For>
      </div>
      <CheckboxRoot validationState="invalid">
        <CheckboxInput />
        <CheckboxControl />
        <CheckboxLabel>Checkbox Label</CheckboxLabel>
        <CheckboxDescription>Checkbox Description</CheckboxDescription>
        <CheckboxErrorMessage>Checkbox Error Message</CheckboxErrorMessage>
      </CheckboxRoot>
      <PopoverRoot>
        <PopoverTrigger asChild>
          <As component={Button}>Popover</As>
        </PopoverTrigger>
        <PopoverPortal>
          <PopoverContent>
            <PopoverArrow />
            <PopoverCloseButton class="float-end">
              <PopoverCloseIcon />
            </PopoverCloseButton>
            <PopoverTitle>Title</PopoverTitle>
            <PopoverDescription>
              Description
            </PopoverDescription>
          </PopoverContent>
        </PopoverPortal>
      </PopoverRoot>
      <DialogRoot>
        <DialogTrigger asChild>
          <As component={Button}>Dialog</As>
        </DialogTrigger>
        <DialogPortal>
          <DialogOverlay />
          <DialogContent>
            <DialogHeader>
              <DialogCloseButton />
              <DialogTitle>Header Title</DialogTitle>
              <DialogDescription>
                Header Description
              </DialogDescription>
            </DialogHeader>
            <DialogFooter>
              Footer
            </DialogFooter>
          </DialogContent>
        </DialogPortal>
      </DialogRoot>
      <TableRoot>
        <TableCaption>Caption</TableCaption>
        <TableHeader>
          <TableRow>
            <TableHead>
              <CheckboxRoot>
                <CheckboxControl />
              </CheckboxRoot>
            </TableHead>
            <TableHead>Head</TableHead>
          </TableRow>
        </TableHeader>
        <TableBody>
          <TableRow>
            <TableCell>
              <CheckboxRoot>
                <CheckboxControl />
              </CheckboxRoot>
            </TableCell>
            <TableCell>Cell</TableCell>
          </TableRow>
          <TableRow>
            <TableCell>
              <CheckboxRoot>
                <CheckboxControl />
              </CheckboxRoot>
            </TableCell>
            <TableCell>Cell</TableCell>
          </TableRow>
        </TableBody>
      </TableRoot>
      <Button onClick={showToast}>Toast</Button>
      <ConfirmButton message="Are you sure?" variant="destructive">
        Confirm Button
      </ConfirmButton>
      <PaginationRoot
        count={10}
        itemComponent={props => (
          <PaginationItem page={props.page}>
            <PaginationLink>
              {props.page}
            </PaginationLink>
          </PaginationItem>
        )}
        ellipsisComponent={() => (
          <PaginationEllipsis />
        )}
      >
        <PaginationContent>
          <PaginationPrevious />
          <PaginationItems />
          <PaginationNext />
        </PaginationContent>
      </PaginationRoot>
      <SelectRoot
        defaultValue="Apple"
        options={["Apple", "Banana", "Blueberry", "Grapes", "Pineapple"]}
        placeholder="Select a fruit…"
        itemComponent={props => (
          <SelectItem item={props.item}>
            {props.item.rawValue}
          </SelectItem>
        )}
      >
        <SelectTrigger aria-label="Fruit">
          <SelectValue<string>>
            {state => state.selectedOption()}
          </SelectValue>
        </SelectTrigger>
        <SelectContent>
          <SelectListbox />
        </SelectContent>
      </SelectRoot>
      <MenubarRoot>
        <MenubarMenu>
          <MenubarTrigger>
            Git
          </MenubarTrigger>
          <MenubarContent>
            <MenubarItem >
              Commit <MenubarShortcut>⌘+K</MenubarShortcut>
            </MenubarItem>
            <MenubarItem>
              Push <MenubarShortcut>⇧+⌘+K</MenubarShortcut>
            </MenubarItem>
            <MenubarItem disabled>
              Update Project <MenubarShortcut>⌘+T</MenubarShortcut>
            </MenubarItem>
            <MenubarSub overlap gutter={4} shift={-8}>
              <MenubarSubTrigger>
                GitHub
              </MenubarSubTrigger>
              <MenubarSubContent>
                <MenubarItem>
                  Create Pull Request…
                </MenubarItem>
                <MenubarItem>
                  View Pull Requests
                </MenubarItem>
                <MenubarItem>Sync Fork</MenubarItem>
                <MenubarSeparator />
                <MenubarItem>
                  Open on GitHub
                </MenubarItem>
              </MenubarSubContent>
            </MenubarSub>
            <MenubarSeparator />
            <MenubarCheckboxItem
            >
              Show Git Log
            </MenubarCheckboxItem>
            <MenubarCheckboxItem>
              Show History
            </MenubarCheckboxItem>
            <MenubarSeparator />
            <MenubarGroup>
              <MenubarGroupLabel>
                Branches
              </MenubarGroupLabel>
              <MenubarRadioGroup>
                <MenubarRadioItem value="main">
                  main
                </MenubarRadioItem>
                <MenubarRadioItem value="develop">
                  develop
                </MenubarRadioItem>
              </MenubarRadioGroup>
            </MenubarGroup>
          </MenubarContent>
        </MenubarMenu>
        <MenubarMenu>
          <MenubarTrigger>
            File
          </MenubarTrigger>
          <MenubarContent>
            <MenubarItem>
              New Tab <MenubarShortcut>⌘+T</MenubarShortcut>
            </MenubarItem>
            <MenubarItem>
              New Window <MenubarShortcut>⌘+N</MenubarShortcut>
            </MenubarItem>
            <MenubarItem disabled>
              New Incognito Window
            </MenubarItem>
            <MenubarSeparator />
            <MenubarSub overlap gutter={4} shift={-8}>
              <MenubarSubTrigger>
                Share
              </MenubarSubTrigger>
              <MenubarSubContent>
                <MenubarItem>
                  Email Link
                </MenubarItem>
                <MenubarItem>
                  Messages
                </MenubarItem>
                <MenubarItem>
                  Notes
                </MenubarItem>
              </MenubarSubContent>
            </MenubarSub>
            <MenubarSeparator />
            <MenubarItem>
              Print... <MenubarShortcut>⌘+P</MenubarShortcut>
            </MenubarItem>
          </MenubarContent>
        </MenubarMenu>
        <MenubarMenu>
          <MenubarTrigger>
            Edit
          </MenubarTrigger>
          <MenubarContent>
            <MenubarItem>
              Undo <MenubarShortcut>⌘+Z</MenubarShortcut>
            </MenubarItem>
            <MenubarItem>
              Redo <MenubarShortcut>⇧+⌘+Z</MenubarShortcut>
            </MenubarItem>
            <MenubarSeparator />
            <MenubarSub overlap gutter={4} shift={-8}>
              <MenubarSubTrigger>
                Find
              </MenubarSubTrigger>
              <MenubarSubContent>
                <MenubarItem>
                  Search The Web
                </MenubarItem>
                <MenubarSeparator />
                <MenubarItem>
                  Find...
                </MenubarItem>
                <MenubarItem>
                  Find Next
                </MenubarItem>
                <MenubarItem>
                  Find Previous
                </MenubarItem>
              </MenubarSubContent>
            </MenubarSub>
            <MenubarSeparator />
            <MenubarItem>
              Cut
            </MenubarItem>
            <MenubarItem>
              Copy
            </MenubarItem>
            <MenubarItem>
              Paste
            </MenubarItem>
          </MenubarContent>
        </MenubarMenu>
      </MenubarRoot>
      <TabsRoot>
        <TabsList>
          <TabsTrigger value="account">Tabs 1 Trigger</TabsTrigger>
          <TabsTrigger value="password">Tabs 2 Trigger</TabsTrigger>
        </TabsList>
        <TabsContent value="account">Tabs 1 Content</TabsContent>
        <TabsContent value="password">Tabs 2 Content</TabsContent>
      </TabsRoot>
      <AvatarRoot fallbackDelay={600}>
        <AvatarImage
          class="image__img"
          src="/vite.svg"
          alt="Vite"
        />
        <AvatarFallback>VT</AvatarFallback>
      </AvatarRoot>

      <ProgressRoot value={progress()}>
        <div class="flex justify-between">
          <ProgressLabel>Loading</ProgressLabel>
          <ProgressValueLabel>
            {progress()}%
          </ProgressValueLabel>
        </div>
        <ProgressTrack>
          <ProgressFill />
        </ProgressTrack>
      </ProgressRoot>

      <Skeleton class="h-screen" />
    </div >
  )
}
