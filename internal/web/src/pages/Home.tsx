import Humanize from "humanize-plus"
import { A, createAsync, revalidate } from "@solidjs/router"
import { CardRoot, } from "~/ui/Card"
import { getHomePage } from "./Home.data"
import { ErrorBoundary, For, ParentProps, Show, Suspense } from "solid-js"
import { BiRegularCctv } from "solid-icons/bi"
import { PageError, PageLoading } from "~/ui/Page"
import { LayoutNormal } from "~/ui/Layout"
import { RiBusinessMailLine, RiMediaVideoLine, RiDeviceHardDrive2Line, RiDocumentFile2Line, RiEditorAttachment2, RiWeatherFlashlightLine, RiMediaImageLine, RiSystemDownloadLine } from "solid-icons/ri"
import { Shared } from "~/components/Shared"
import { formatDate, parseDate } from "~/lib/utils"
import { Seperator } from "~/ui/Seperator"
import { TooltipArrow, TooltipContent, TooltipRoot, TooltipTrigger } from "~/ui/Tooltip"
import { createDate, createTimeAgo } from "@solid-primitives/date"
import { Image } from "@kobalte/core"
import { linkVariants } from "~/ui/Link"
import { useBus } from "~/providers/bus"

export function Home() {
  const data = createAsync(() => getHomePage())
  const bus = useBus()

  bus.event.listen((e) => {
    if (e.action.startsWith("dahua-email:") || e.action.startsWith("dahua-device:"))
      revalidate(getHomePage.key)
  })

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
            <div class="flex-1 lg:max-w-sm">
              <CardRoot class="p-4">
                <Shared.Title>Latest emails</Shared.Title>
                <div>
                  <For each={data()?.emails}>
                    {v => {
                      const [createdAt] = createDate(() => parseDate(v.createdAtTime));
                      const [createdAtAgo] = createTimeAgo(createdAt);

                      return (
                        <div class="hover:bg-muted/50 flex flex-col border-b transition-colors sm:flex-row">
                          <A href={`/emails/${v.id}`} class="flex flex-1 flex-col gap-2 p-2 max-sm:pb-1 sm:flex-row sm:pr-1">
                            <div class="sm:min-w-32 flex">
                              <TooltipRoot>
                                <TooltipTrigger class="truncate text-start text-sm font-bold">{createdAtAgo()}</TooltipTrigger>
                                <TooltipContent>
                                  <TooltipArrow />
                                  {formatDate(createdAt())}
                                </TooltipContent>
                              </TooltipRoot>
                            </div>
                            <div class="flex-1 truncate">
                              {v.subject}
                            </div>
                          </A>
                          <Show when={v.attachmentCount > 0}>
                            <A href={`/emails/${v.id}?tab=attachments`} class="p-2 max-sm:pt-1 sm:pl-1">
                              <TooltipRoot>
                                <TooltipTrigger class="flex h-full items-center">
                                  <RiEditorAttachment2 class="h-5 w-5" />
                                </TooltipTrigger>
                                <TooltipContent>
                                  <TooltipArrow />
                                  {v.attachmentCount} {Humanize.pluralize(v.attachmentCount, "attachment")}
                                </TooltipContent>
                              </TooltipRoot>
                            </A>
                          </Show>
                        </div>
                      )
                    }}
                  </For>
                </div>
              </CardRoot>
            </div>
            <div class="flex flex-1 flex-col gap-4">
              <Shared.Title>Latest files</Shared.Title>
              <div class="grid grid-cols-2 gap-4 sm:grid-cols-4 xl:grid-cols-6 2xl:grid-cols-8">
                <For each={data()?.files}>
                  {(v) => {
                    const [startTime] = createDate(() => parseDate(v.startTime));
                    const [startTimeAgo] = createTimeAgo(startTime);

                    return (
                      <div>
                        <div class="hover:bg-accent/50 sm:max-w-48 flex w-full flex-col rounded-b border transition-all">
                          <A href={`/files/${v.id}`} >
                            <Image.Root class="mx-auto max-h-48 w-full">
                              <Image.Img src={v.thumbnailUrl} class="h-full w-full object-contain" />
                              <Image.Fallback>
                                <Show when={v.type == "jpg"} fallback={
                                  <RiMediaVideoLine class="h-full w-full object-contain" />
                                }>
                                  <RiMediaImageLine class="h-full w-full object-contain" />
                                </Show>
                              </Image.Fallback>
                            </Image.Root>
                          </A>
                          <Seperator />
                          <div class="flex items-center justify-between gap-2 p-2">
                            <TooltipRoot>
                              <TooltipTrigger class="truncate text-sm">{startTimeAgo()}</TooltipTrigger>
                              <TooltipContent>
                                <TooltipArrow />
                                {formatDate(startTime())}
                              </TooltipContent>
                            </TooltipRoot>
                            <a href={v.url} target="_blank" title="Download">
                              <RiSystemDownloadLine class="h-5 w-5" />
                            </a>
                          </div>
                        </div>
                      </div>
                    )
                  }}
                </For>
              </div>
            </div>
          </div>
          <div class="flex flex-col sm:flex-row">
            <CardRoot class="p-4">
              <Shared.Title>Build</Shared.Title>
              <div class="relative overflow-x-auto">
                <table class="w-full">
                  <tbody>
                    <tr class="border-b">
                      <td class="p-2">Commit</td>
                      <td class="p-2"><a class={linkVariants()} href={data()?.build?.commitUrl}>{data()?.build?.commit}</a></td>
                    </tr>
                    <tr class="border-b">
                      <td class="p-2">Date</td>
                      <td class="p-2">{formatDate(parseDate(data()?.build?.date))}</td>
                    </tr>
                    <tr class="border-b">
                      <td class="p-2">Version</td>
                      <td class="p-2"><a class={linkVariants()} href={data()?.build?.releaseUrl}>{data()?.build?.version}</a></td>
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
  return <h2 class="text-nowrap">{props.children}</h2>
}

function StatValue(props: ParentProps) {
  return <p class="text-nowrap text-lg font-bold">{props.children}</p>
}
