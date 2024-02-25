import { Timestamp } from "~/twirp/google/protobuf/timestamp";
import { createAsync, revalidate, useSearchParams } from "@solidjs/router";
import { ErrorBoundary, Suspense, resetErrorBoundaries } from "solid-js";
import { getHello } from "./Twirp.data"
import { createLoading } from "~/lib/utils";

function queryInt(str: string | undefined) {
  if (!str) {
    return 0
  }

  const id = parseInt(str)
  if (isNaN(id)) {
    return 0
  }

  return id
}

export function Twirp() {
  const data = createAsync(() => getHello())
  const [loading, refreshData] = createLoading(() => revalidate(getHello.key).then(resetErrorBoundaries))
  const text = () => data() ? data()!.text + " " + Timestamp.toDate(Timestamp.create(data()!.currentTime)) : ""

  const [searchParams, setSearchParams] = useSearchParams()
  const count = () => queryInt(searchParams.count)
  const incrementCount = () => setSearchParams({ count: count() + 1 })

  return (
    <>
      <button onClick={incrementCount}>{count()}</button>
      <br />
      <ErrorBoundary fallback={(error) => (
        <>
          <div>{error.message}</div>
          <button onClick={refreshData} disabled={loading()}>Retry</button>
        </>
      )}>
        <Suspense fallback={<>Loading...</>}>
          <div>
            {text()}
          </div>
          <button onClick={refreshData} disabled={loading()}>Refresh</button>
        </Suspense>
      </ErrorBoundary>
    </>
  )
}
