import { A, createAsync } from "@solidjs/router"
import { CardRoot, } from "~/ui/Card"
import { getHomePage } from "./Home.data"
import { ErrorBoundary, ParentProps, Suspense } from "solid-js"
import { BiRegularCctv } from "solid-icons/bi"
import { PageError, PageLoading } from "~/ui/Page"
import { LayoutNormal } from "~/ui/Layout"
import { RiBusinessMailLine, RiDeviceHardDrive2Line, RiDocumentFile2Line, RiWeatherFlashlightLine } from "solid-icons/ri"

export function Home() {
  const data = createAsync(getHomePage)

  return (
    <LayoutNormal>
      <ErrorBoundary fallback={(e) => <PageError error={e} />}>
        <Suspense fallback={<PageLoading />}>
          <RowRoot>
            <RowItem>
              <StatRoot>
                <A class="flex items-center" href="/devices">
                  <BiRegularCctv class="h-8 w-8" />
                </A>
                <div class="flex-1">
                  <StatTitle>Devices</StatTitle>
                  <StatValue>{data()?.devices.length}</StatValue>
                </div>
              </StatRoot>
            </RowItem>
            <RowItem>
              <StatRoot>
                <A class="flex items-center" href="/files">
                  <RiDocumentFile2Line class="h-8 w-8" />
                </A>
                <div class="flex-1">
                  <StatTitle>Files</StatTitle>
                  <StatValue>{data()?.fileCount.toString()}</StatValue>
                </div>
              </StatRoot>
            </RowItem>
            <RowItem>
              <StatRoot>
                <A class="flex items-center" href="/events">
                  <RiWeatherFlashlightLine class="h-8 w-8" />
                </A>
                <div class="flex-1">
                  <StatTitle>Events</StatTitle>
                  <StatValue>{data()?.eventCount.toString()}</StatValue>
                </div>
              </StatRoot>
            </RowItem>
            <RowItem>
              <StatRoot>
                <div class="flex items-center">
                  <RiBusinessMailLine class="h-8 w-8" />
                </div>
                <div class="flex-1">
                  <StatTitle>Emails</StatTitle>
                  <StatValue>N/A</StatValue>
                </div>
              </StatRoot>
            </RowItem>
            <RowItem>
              <StatRoot>
                <div class="flex items-center">
                  <RiDeviceHardDrive2Line class="h-8 w-8" />
                </div>
                <div class="flex-1">
                  <StatTitle>Disk usage</StatTitle>
                  <StatValue>N/A</StatValue>
                </div>
              </StatRoot>
            </RowItem>
          </RowRoot>
        </Suspense>
      </ErrorBoundary>
    </LayoutNormal>
  )
}

function RowRoot(props: ParentProps) {
  return <div class="flex flex-col flex-wrap gap-2 sm:flex-row">{props.children}</div>
}

function RowItem(props: ParentProps) {
  return <div class="sm:max-w-48 flex-1">{props.children}</div>
}

function StatRoot(props: ParentProps) {
  return <CardRoot class="flex gap-2 p-4">{props.children}</CardRoot>
}

function StatTitle(props: ParentProps) {
  return <div class="text-nowrap">{props.children}</div>
}

function StatValue(props: ParentProps) {
  return <div class="text-nowrap text-lg font-bold">{props.children}</div>
}
