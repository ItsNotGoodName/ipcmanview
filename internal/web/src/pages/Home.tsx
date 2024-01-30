import { createAsync } from "@solidjs/router"
import { CardRoot } from "~/ui/Card"
import { getHomePage } from "./Home.data"
import { ErrorBoundary, Suspense } from "solid-js"
import { BiRegularCctv } from "solid-icons/bi"
import { PageError, PageLoading } from "~/ui/Page"
import { LayoutNormal } from "~/ui/Layout"

export function Home() {
  const data = createAsync(getHomePage)

  return (
    <LayoutNormal>
      <ErrorBoundary fallback={(e) => <PageError error={e} />}>
        <Suspense fallback={<PageLoading class="flex-1" />}>
          <div class="flex">
            <CardRoot class="flex gap-2 p-4">
              <div class="flex items-center">
                <BiRegularCctv class="h-8 w-8" />
              </div>
              <div>
                <div>Total Devices</div>
                <div class="text-xl font-bold">{data()?.deviceCount.toString()}</div>
              </div>
            </CardRoot>
          </div>
        </Suspense>
      </ErrorBoundary>
    </LayoutNormal>
  )
}

