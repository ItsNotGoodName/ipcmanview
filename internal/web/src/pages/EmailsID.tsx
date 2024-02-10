import { useSearchParams } from "@solidjs/router"
import { Crud } from "~/components/Crud"
import { Shared } from "~/components/Shared"
import { formatDate } from "~/lib/utils"
import { Button } from "~/ui/Button"
import { CardRoot } from "~/ui/Card"
import { RiArrowsArrowLeftLine } from "solid-icons/ri"
import { LayoutNormal } from "~/ui/Layout"
import { TabsContent, TabsList, TabsRoot, TabsTrigger } from "~/ui/Tabs"

export function EmailsID({ params }: any) {
  const [searchParams, setSearchParams] = useSearchParams()

  return (
    <LayoutNormal class="max-w-4xl">
      <Shared.Title>
        Emails / 2
      </Shared.Title>
      <div class="flex items-center gap-2 justify-between">
        <div>
          <Button size="icon" variant="ghost" title="Back">
            <RiArrowsArrowLeftLine class="h-5 w-5" />
          </Button>
        </div>
        <div class="flex items-center gap-2">
          <div>1 of 10</div>
          <Crud.PageButtons />
        </div>
      </div>
      <CardRoot>
        <div class="overflow-x-auto p-2">
          <table>
            <tbody>
              <tr>
                <th class="px-2">From</th>
                <td class="px-2">from@example.com</td>
              </tr>
              <tr>
                <th class="px-2">Subject</th>
                <td class="px-2">Example Subject</td>
              </tr>
              <tr>
                <th class="px-2">To</th>
                <td class="px-2">to@example.com</td>
              </tr>
              <tr>
                <th class="px-2">Date</th>
                <td class="px-2">{formatDate(new Date())}</td>
              </tr>
              <tr>
                <th class="px-2">Created At</th>
                <td class="px-2">{formatDate(new Date())}</td>
              </tr>
            </tbody>
          </table>
        </div>
      </CardRoot>
      <TabsRoot value={searchParams.tab || "text"} onChange={(value) => setSearchParams({ tab: value })}>
        <div class="flex overflow-x-auto">
          <TabsList>
            <TabsTrigger value="text" >Text</TabsTrigger>
            <TabsTrigger value="attachments" >Attachments</TabsTrigger>
          </TabsList>
        </div>
        <TabsContent value="text">
          Text
        </TabsContent>
        <TabsContent value="attachments">
          Attachments
        </TabsContent>
      </TabsRoot>
    </LayoutNormal>
  )
}
