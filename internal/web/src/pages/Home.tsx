import { A, createAsync } from "@solidjs/router"
import { CardRoot, } from "~/ui/Card"
import { getHomePage } from "./Home.data"
import { ErrorBoundary, For, ParentProps, Show, Suspense } from "solid-js"
import { BiRegularCctv } from "solid-icons/bi"
import { PageError, PageLoading } from "~/ui/Page"
import { LayoutNormal } from "~/ui/Layout"
import { RiBusinessMailLine, RiMediaVideoLine, RiDeviceHardDrive2Line, RiDocumentFile2Line, RiEditorAttachment2, RiWeatherFlashlightLine, RiMediaImageLine } from "solid-icons/ri"
import { Shared } from "~/components/Shared"
import { formatDate } from "~/lib/utils"
import { Seperator } from "~/ui/Seperator"
import { TooltipContent, TooltipRoot, TooltipTrigger } from "~/ui/Tooltip"

export function Home() {
  const data = createAsync(getHomePage)

  return (
    <LayoutNormal>
      <ErrorBoundary fallback={(e) => <PageError error={e} />}>
        <Suspense fallback={<PageLoading />}>
          <div class="flex flex-col flex-wrap gap-4 sm:flex-row">
            <StatParent>
              <StatRoot>
                <A class="flex items-center" href="/devices">
                  <BiRegularCctv class="h-8 w-8" />
                </A>
                <div class="flex-1">
                  <StatTitle>Devices</StatTitle>
                  <StatValue>{data()?.devices.length}</StatValue>
                </div>
              </StatRoot>
            </StatParent>
            <StatParent>
              <StatRoot>
                <A class="flex items-center" href="/emails">
                  <RiBusinessMailLine class="h-8 w-8" />
                </A>
                <div class="flex-1">
                  <StatTitle>Emails</StatTitle>
                  <StatValue>{data()?.emailCount.toString()}</StatValue>
                </div>
              </StatRoot>
            </StatParent>
            <StatParent>
              <StatRoot>
                <A class="flex items-center" href="/events">
                  <RiWeatherFlashlightLine class="h-8 w-8" />
                </A>
                <div class="flex-1">
                  <StatTitle>Events</StatTitle>
                  <StatValue>{data()?.eventCount.toString()}</StatValue>
                </div>
              </StatRoot>
            </StatParent>
            <StatParent>
              <StatRoot>
                <A class="flex items-center" href="/files">
                  <RiDocumentFile2Line class="h-8 w-8" />
                </A>
                <div class="flex-1">
                  <StatTitle>Files</StatTitle>
                  <StatValue>{data()?.fileCount.toString()}</StatValue>
                </div>
              </StatRoot>
            </StatParent>
            <StatParent>
              <StatRoot>
                <div class="flex items-center">
                  <RiDeviceHardDrive2Line class="h-8 w-8" />
                </div>
                <div class="flex-1">
                  <StatTitle>Disk usage</StatTitle>
                  <StatValue>N/A</StatValue>
                </div>
              </StatRoot>
            </StatParent>
          </div>
          <div class="flex flex-col gap-4 lg:flex-row">
            <CardRoot class="flex-shrink-0 p-4 lg:max-w-sm">
              <Shared.Title>Latest emails</Shared.Title>
              <table class="w-full table-fixed">
                <tbody>
                  <For each={Array(5)}>
                    {_ =>
                      <tr class="hover:bg-muted/50 flex flex-col overflow-hidden border-b py-2 transition-colors sm:flex-row">
                        <td class="text-nowrap px-2 font-bold">
                          <TooltipRoot>
                            <TooltipTrigger>19 minutes ago</TooltipTrigger>
                            <TooltipContent>
                              {formatDate(new Date())}
                            </TooltipContent>
                          </TooltipRoot>
                        </td>
                        <td class="w-full truncate px-2">
                          <A href={`/emails/3`}>
                            Lorem ipsum dolor sit amet consectetur, adipisicing elit. Dolore, illum nostrum sit iusto architecto nobis tenetur repellat enim cumque esse quia consequatur hic, quidem velit magni modi alias beatae eos.
                          </A>
                        </td>
                        <td class="px-1">
                          <A href={`/emails/3?tab=attachments`}>
                            <TooltipRoot>
                              <TooltipTrigger class="flex h-full items-center">
                                <RiEditorAttachment2 />
                              </TooltipTrigger>
                              <TooltipContent>
                                2 attachment(s)
                              </TooltipContent>
                            </TooltipRoot>
                          </A>
                        </td>
                      </tr>
                    }
                  </For>
                </tbody>
              </table>
            </CardRoot>
            <CardRoot class="flex-2 p-4">
              <Shared.Title>Latest files</Shared.Title>
              <div class="grid gap-4 pt-4 sm:grid-cols-4 xl:grid-cols-6 2xl:grid-cols-8">
                <For each={Array(8)}>
                  {(_, i) =>
                    <div class="hover:bg-accent/50 flex flex-col gap-1 rounded border p-2 transition-all">
                      <A href={`/files/3`}>
                        <Show when={i() % 2 == 0} fallback={
                          <RiMediaImageLine class="aspect-square h-full w-full" />
                        }>
                          <RiMediaVideoLine class="aspect-square h-full w-full" />
                        </Show>
                      </A>
                      <Seperator />
                      <TooltipRoot>
                        <TooltipTrigger class="text-sm">19 minutes ago</TooltipTrigger>
                        <TooltipContent>
                          {formatDate(new Date())}
                        </TooltipContent>
                      </TooltipRoot>
                    </div>
                  }
                </For>
              </div>
            </CardRoot>
          </div>
          <div class="max-w-sm">
            <CardRoot class="p-4 ">
              <Shared.Title>Build</Shared.Title>
              <div class="relative overflow-x-auto">
                <table class="w-full">
                  <tbody>
                    <tr class="border-b">
                      <td class="p-2">Commit</td>
                      <td class="p-2"><a href={data()?.build?.commitUrl}>{data()?.build?.commit}</a></td>
                    </tr>
                    <tr class="border-b">
                      <td class="p-2">Date</td>
                      <td class="p-2">{data()?.build?.date}</td>
                    </tr>
                    <tr class="border-b">
                      <td class="p-2">Version</td>
                      <td class="p-2">{data()?.build?.version}</td>
                    </tr>
                  </tbody>
                </table>
              </div>
            </CardRoot>
          </div>
        </Suspense>
      </ErrorBoundary>
    </LayoutNormal >
  )
}

function StatParent(props: ParentProps) {
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
