import { createAsync } from "@solidjs/router"
import { CardRoot } from "~/ui/Card"
import { getHome } from "./Home.data"
import { ErrorBoundary, Suspense } from "solid-js"
import { BiRegularCctv } from "solid-icons/bi"
import { PageError, PageLoading } from "~/ui/Page"

export function Home() {
  const data = createAsync(getHome)

  return (
    <div class="flex p-4">
      <ErrorBoundary fallback={(e: Error) => <PageError error={e} />}>
        <Suspense fallback={<PageLoading class="flex-1" />}>
          <CardRoot class="flex gap-2 p-4">
            <div class="flex items-center">
              <BiRegularCctv class="h-8 w-8" />
            </div>
            <div>
              <div>Total Devices</div>
              <div class="text-xl font-bold">{data()?.deviceCount.toString()}</div>
            </div>
          </CardRoot>
        </Suspense>
      </ErrorBoundary>
    </div>
  )
}

