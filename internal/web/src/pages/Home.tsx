import { A, createAsync } from "@solidjs/router"
import { CardRoot, } from "~/ui/Card"
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
        <Suspense fallback={<PageLoading />}>
          <div class="flex gap-2">
            <div>
              <CardRoot class="flex gap-2 p-4">
                <A class="flex items-center" href="/devices">
                  <BiRegularCctv class="h-8 w-8" />
                </A>
                <div>
                  <div>Devices</div>
                  <div class="text-xl font-bold">{data()?.devices.length}</div>
                </div>
              </CardRoot>
            </div>
          </div>
        </Suspense>
      </ErrorBoundary>
    </LayoutNormal>
  )
}
