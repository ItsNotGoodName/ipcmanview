import { createAsync, revalidate } from "@solidjs/router"
import { ErrorBoundary, Suspense, resetErrorBoundaries } from "solid-js"

import { formatDate, parseDate, useLoading } from "~/hooks"
import { CardContent, CardHeader, CardRoot, CardTitle } from "~/ui/Card"
import { getProfile } from "./Profile.data"
import { AlertDescription, AlertRoot, AlertTitle } from "~/ui/Alert"
import { Button } from "~/ui/Button"
import { Skeleton } from "~/ui/Skeleton"

export function Profile() {
  const data = createAsync(getProfile)
  const [loading, refreshData] = useLoading(() => revalidate(getProfile.key).then(resetErrorBoundaries))

  return (
    <div class="p-4">
      <ErrorBoundary fallback={(error) => (
        <AlertRoot>
          <AlertTitle>{error.message}</AlertTitle>
          <AlertDescription>
            <Button onClick={refreshData} disabled={loading()}>Retry</Button>
          </AlertDescription>
        </AlertRoot>
      )}>
        <Suspense fallback={
          <Skeleton class="w-full h-32" />
        }>
          <CardRoot>
            <CardHeader>
              <CardTitle>User</CardTitle>
            </CardHeader>
            <CardContent>
              <div>{data()?.username}</div>
              <div>{formatDate(parseDate(data()?.createdAt))}</div>
              <div>{formatDate(parseDate(data()?.updatedAt))}</div>
            </CardContent>
          </CardRoot>
        </Suspense>
      </ErrorBoundary>
    </div>
  )
}


