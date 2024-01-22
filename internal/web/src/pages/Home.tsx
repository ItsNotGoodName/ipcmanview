import { createAsync } from "@solidjs/router"
import { CardRoot } from "~/ui/Card"
import { getHome } from "./Home.data"
import { ErrorBoundary, Suspense } from "solid-js"
import { Loading } from "~/ui/Loading"
import { AlertRoot, AlertTitle } from "~/ui/Alert"
import { BiRegularCctv } from "solid-icons/bi"

export function Home() {
  const data = createAsync(getHome)

  return (
    <ErrorBoundary fallback={(e: Error) =>
      <AlertRoot variant="destructive">
        <AlertTitle>
          {e.message}
        </AlertTitle>
      </AlertRoot>
    }>
      <Suspense fallback={<Loading />}>
        <div class="flex p-4">
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
  )
}

