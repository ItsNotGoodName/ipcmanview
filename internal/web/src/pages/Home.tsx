import { createAsync, useSearchParams } from "@solidjs/router"
import { CardRoot, } from "~/ui/Card"
import { getHomePage } from "./Home.data"
import { ErrorBoundary, For, Show, Suspense } from "solid-js"
import { BiRegularCctv } from "solid-icons/bi"
import { PageError, PageLoading } from "~/ui/Page"
import { LayoutNormal } from "~/ui/Layout"
import { TabsContent, TabsList, TabsRoot, TabsTrigger } from "~/ui/Tabs"
import { TableBody, TableCell, TableHead, TableHeader, TableRoot, TableRow } from "~/ui/Table"
import { GetHomePageResp_Device } from "~/twirp/rpc"
import { getDeviceDetail } from "./data"
import { Skeleton } from "~/ui/Skeleton"
import { ToggleButton } from "@kobalte/core"

export function Home() {
  const data = createAsync(getHomePage)
  const [searchParams, setSearchParams] = useSearchParams()

  return (
    <LayoutNormal>
      <ErrorBoundary fallback={(e) => <PageError error={e} />}>
        <Suspense fallback={<PageLoading />}>
          <div class="flex gap-2">
            <div>
              <CardRoot class="flex gap-2 p-4">
                <div class="flex items-center">
                  <BiRegularCctv class="h-8 w-8" />
                </div>
                <div>
                  <div>Devices</div>
                  <div class="text-xl font-bold">{data()?.devices.length}</div>
                </div>
              </CardRoot>
            </div>
          </div>
          <TabsRoot value={searchParams.tab || "device"} onChange={(value) => setSearchParams({ tab: value })}>
            <TabsList class="w-full overflow-x-auto overflow-y-hidden">
              <TabsTrigger value="device" >Device</TabsTrigger>
              <TabsTrigger value="status" >Status</TabsTrigger>
              <TabsTrigger value="detail" >Detail</TabsTrigger>
              <TabsTrigger value="software-version" >Software Version</TabsTrigger>
              <TabsTrigger value="license" >License</TabsTrigger>
              <TabsTrigger value="storage" >Storage</TabsTrigger>
            </TabsList>
            <TabsContent value="device">
              <DeviceTable devices={data()?.devices} />
            </TabsContent>
            <TabsContent value="status">
              <StatusTable />
            </TabsContent>
            <TabsContent value="detail">
              <DetailTable devices={data()?.devices} />
            </TabsContent>
            <TabsContent value="software-version">
              <SoftwareVersionTable />
            </TabsContent>
            <TabsContent value="license">
              <LicenseTable />
            </TabsContent>
            <TabsContent value="storage">
              <StorageTable />
            </TabsContent>
          </TabsRoot>
        </Suspense>
      </ErrorBoundary>
    </LayoutNormal>
  )
}

function DeviceTable(props: { devices?: GetHomePageResp_Device[] }) {
  return (
    <TableRoot>
      <TableHeader>
        <TableRow>
          <TableHead>Name</TableHead>
        </TableRow>
      </TableHeader>
      <TableBody>
        <For each={props.devices}>
          {v => (
            <TableRow>
              <TableCell>
                {v.name}
              </TableCell>
            </TableRow>
          )}
        </For>
      </TableBody>
    </TableRoot>
  )
}

function StatusTable() {
  return (
    <TableRoot>
      <TableHeader>
        <TableRow>
          <TableHead>Device</TableHead>
          <TableHead>URL</TableHead>
          <TableHead>Username</TableHead>
          <TableHead>Location</TableHead>
          <TableHead>Seed</TableHead>
          <TableHead>RPC Error</TableHead>
          <TableHead>RPC State</TableHead>
          <TableHead>RPC Last Login</TableHead>
        </TableRow>
      </TableHeader>
      <TableBody>
      </TableBody>
    </TableRoot>
  )
}

function DetailTable(props: { devices?: GetHomePageResp_Device[] }) {
  const colspan = 9

  return (
    <TableRoot>
      <TableHeader>
        <TableRow>
          <TableHead>Device</TableHead>
          <TableHead>SN</TableHead>
          <TableHead>Device Class</TableHead>
          <TableHead>Device Type</TableHead>
          <TableHead>Hardware Version</TableHead>
          <TableHead>Market Area</TableHead>
          <TableHead>Process Info</TableHead>
          <TableHead>Vendor</TableHead>
          <TableHead>Onvif Version</TableHead>
          <TableHead>Algorithm Version</TableHead>
        </TableRow>
      </TableHeader>
      <TableBody>
        <For each={props.devices}>
          {item => {
            const data = createAsync(() => getDeviceDetail(item.id))

            return (
              <TableRow>
                <TableCell>
                  {item.name}
                </TableCell>
                <ErrorBoundary fallback={(e: Error) => (
                  <TableCell colspan={colspan} class="py-0">
                    <div class="bg-destructive text-destructive-foreground rounded p-2">
                      {e.message}
                    </div>
                  </TableCell>
                )}>
                  <Suspense fallback={
                    <TableCell colspan={colspan} class="py-0">
                      <Skeleton class="h-8" />
                    </TableCell>
                  }>
                    <TableCell>
                      <ToggleButton.Root>
                        {state => (
                          <Show when={state.pressed()} fallback={<>***************</>}>
                            {data()?.sn}
                          </Show>
                        )}
                      </ToggleButton.Root>
                    </TableCell>
                    <TableCell>
                      {data()?.deviceClass}
                    </TableCell>
                    <TableCell>
                      {data()?.deviceType}
                    </TableCell>
                    <TableCell>
                      {data()?.hardwareVersion}
                    </TableCell>
                    <TableCell>
                      {data()?.marketArea}
                    </TableCell>
                    <TableCell>
                      {data()?.processInfo}
                    </TableCell>
                    <TableCell>
                      {data()?.vendor}
                    </TableCell>
                    <TableCell>
                      {data()?.onvifVersion}
                    </TableCell>
                    <TableCell>
                      {data()?.algorithmVersion}
                    </TableCell>
                  </Suspense>
                </ErrorBoundary>
              </TableRow>
            )
          }}
        </For>
      </TableBody>
    </TableRoot>
  )
}

function SoftwareVersionTable() {
  return (
    <TableRoot>
      <TableHeader>
        <TableRow>
          <TableHead>Device</TableHead>
          <TableHead>Build</TableHead>
          <TableHead>Build Date</TableHead>
          <TableHead>Security Base Line Version</TableHead>
          <TableHead>Version</TableHead>
          <TableHead>Web Version</TableHead>
        </TableRow>
      </TableHeader>
      <TableBody>
      </TableBody>
    </TableRoot>
  )
}

function LicenseTable() {
  return (
    <TableRoot>
      <TableHeader>
        <TableRow>
          <TableHead>Device</TableHead>
          <TableHead>Abroad Info</TableHead>
          <TableHead>All Type</TableHead>
          <TableHead>Digit Channel</TableHead>
          <TableHead>Effective Days</TableHead>
          <TableHead>Effective Time</TableHead>
          <TableHead>License ID</TableHead>
          <TableHead>Product Type</TableHead>
          <TableHead>Status</TableHead>
          <TableHead>Username</TableHead>
        </TableRow>
      </TableHeader>
      <TableBody>
      </TableBody>
    </TableRoot>
  )
}

function StorageTable() {
  return (
    <TableRoot>
      <TableHeader>
        <TableRow>
          <TableHead>Device</TableHead>
          <TableHead>Name</TableHead>
          <TableHead>State</TableHead>
          <TableHead>Type</TableHead>
          <TableHead>Total</TableHead>
          <TableHead>Used</TableHead>
          <TableHead>Is Error</TableHead>
        </TableRow>
      </TableHeader>
      <TableBody>
      </TableBody>
    </TableRoot>
  )
}
