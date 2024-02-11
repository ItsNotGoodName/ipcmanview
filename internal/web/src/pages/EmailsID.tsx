import Humanize from "humanize-plus"
import { A, createAsync, useSearchParams } from "@solidjs/router"
import { Crud } from "~/components/Crud"
import { Shared } from "~/components/Shared"
import { formatDate, parseDate } from "~/lib/utils"
import { buttonVariants } from "~/ui/Button"
import { CardRoot } from "~/ui/Card"
import { RiArrowsArrowLeftLine, RiDeviceHardDrive2Line, RiMediaImageLine, RiSystemDownloadLine } from "solid-icons/ri"
import { LayoutNormal } from "~/ui/Layout"
import { TabsContent, TabsList, TabsRoot, TabsTrigger } from "~/ui/Tabs"
import { BreadcrumbsItem, BreadcrumbsLink, BreadcrumbsList, BreadcrumbsRoot, BreadcrumbsSeparator } from "~/ui/Breadcrumbs"
import { As } from "@kobalte/core"
import { getEmailsIDPage } from "./EmailsID.data"
import { ErrorBoundary, For, Show, Suspense } from "solid-js"
import { Skeleton } from "~/ui/Skeleton"
import { PageError } from "~/ui/Page"
import { Badge } from "~/ui/Badge"
import { Seperator } from "~/ui/Seperator"
import { Image } from "@kobalte/core"
import { TooltipContent, TooltipRoot, TooltipTrigger } from "~/ui/Tooltip"

export function EmailsID({ params }: any) {
  const [searchParams, setSearchParams] = useSearchParams()
  const data = createAsync(() => getEmailsIDPage(BigInt(params.id)))
  const query = () => searchParams.tab ? "?tab=" + searchParams.tab : ""
  const backPage = () => Math.ceil(Number(data()?.emailSeen) / 10) || 1

  return (
    <LayoutNormal class="max-w-4xl">
      <Shared.Title>
        <BreadcrumbsRoot>
          <BreadcrumbsList>
            <BreadcrumbsItem>
              <BreadcrumbsLink asChild>
                <As component={A} href="/emails">
                  Emails
                </As>
              </BreadcrumbsLink>
              <BreadcrumbsSeparator />
            </BreadcrumbsItem>
            <BreadcrumbsItem>
              <BreadcrumbsLink>
                {params.id}
              </BreadcrumbsLink>
            </BreadcrumbsItem>
          </BreadcrumbsList>
        </BreadcrumbsRoot>
      </Shared.Title>
      <ErrorBoundary fallback={(e) => <PageError error={e} />}>
        <Suspense fallback={<Skeleton class="h-32" />}>
          <div class="flex items-center justify-between gap-2">
            <div>
              <A href={`/emails?page=${backPage()}`} title="Back" class={buttonVariants({ size: "icon", variant: "ghost" })}>
                <RiArrowsArrowLeftLine class="h-5 w-5" />
              </A>
            </div>
            <div class="flex items-center gap-2">
              <div>{data()?.emailSeen.toString()} of {data()?.emailCount.toString()}</div>
              <Crud.PageButtonsLinks
                previousPage={`/emails/${data()?.previousEmailId}${query()}`}
                previousPageDisabled={data()?.previousEmailId == data()?.id}
                nextPageDisabled={data()?.nextEmailId == data()?.id}
                nextPage={`/emails/${data()?.nextEmailId}${query()}`}
              />
            </div>
          </div>
          <CardRoot>
            <div class="overflow-x-auto p-2">
              <table>
                <tbody>
                  <tr>
                    <th class="px-2">From</th>
                    <td class="px-2">{data()?.from}</td>
                  </tr>
                  <tr>
                    <th class="px-2">Subject</th>
                    <td class="px-2">{data()?.subject}</td>
                  </tr>
                  <tr>
                    <th class="px-2">To</th>
                    <td class="flex gap-2 px-2">
                      <For each={data()?.to}>
                        {v => <Badge>{v}</Badge>}
                      </For>
                    </td>
                  </tr>
                  <tr>
                    <th class="px-2">Date</th>
                    <td class="px-2">{formatDate(parseDate(data()?.date))}</td>
                  </tr>
                  <tr>
                    <th class="px-2">Created At</th>
                    <td class="px-2">{formatDate(parseDate(data()?.createdAtTime))}</td>
                  </tr>
                </tbody>
              </table>
            </div>
          </CardRoot>
          <TabsRoot value={searchParams.tab || "text"} onChange={(value) => setSearchParams({ tab: value })}>
            <div class="overflow-x-auto">
              <TabsList>
                <TabsTrigger value="text" >Text</TabsTrigger>
                <TabsTrigger value="attachments" class="flex items-center gap-2">
                  Attachments
                  <Show when={data()?.attachments.length || 0 > 0}>
                    <Badge>{data()?.attachments.length}</Badge>
                  </Show>
                </TabsTrigger>
              </TabsList>
            </div>
            <TabsContent value="text">
              <pre class="overflow-x-auto">{data()?.text}</pre>
            </TabsContent>
            <TabsContent value="attachments" class="flex flex-wrap gap-4">
              <For each={data()?.attachments}>
                {v => (
                  <div class="sm:max-w-48 flex w-full flex-col rounded-b border">
                    <Image.Root class="mx-auto max-h-48 w-full">
                      <Image.Img src={v.thumbnailUrl} class="h-full w-full object-contain" />
                      <Image.Fallback>
                        <RiMediaImageLine class="h-full w-full object-contain" />
                      </Image.Fallback>
                    </Image.Root>
                    <Seperator />
                    <div class="p-2">
                      <div>
                        <TooltipRoot>
                          <TooltipTrigger class="w-full truncate">{v.name}</TooltipTrigger>
                          <TooltipContent>{v.name}</TooltipContent>
                        </TooltipRoot>
                      </div>
                      <div class="flex items-center justify-between gap-2">
                        <div title="Size" class="flex items-center gap-1">
                          <RiDeviceHardDrive2Line class="h-5 w-5" />
                          {Humanize.fileSize(Number(v.size))}
                        </div>
                        <a href={v.url} target="_blank" title="Download">
                          <RiSystemDownloadLine class="h-5 w-5" />
                        </a>
                      </div>
                    </div>
                  </div>
                )}
              </For>
            </TabsContent>
          </TabsRoot>
        </Suspense>
      </ErrorBoundary>
    </LayoutNormal>
  )
}
