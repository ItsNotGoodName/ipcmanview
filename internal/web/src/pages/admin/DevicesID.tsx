import { createAsync } from "@solidjs/router";
import { ErrorBoundary } from "solid-js";
import { PageError } from "~/ui/Page";
import { LayoutNormal } from "~/ui/Layout";
import { getAdminDevicesIDPage } from "./DevicesID.data";

export function AdminDevicesID(props: any) {
  const data = createAsync(() => getAdminDevicesIDPage(BigInt(props.params.id || 0)))

  return (
    <LayoutNormal>
      <ErrorBoundary fallback={(e) => <PageError error={e} />}>
        <pre>
          {JSON.stringify(data(), null, 2)}
        </pre>
      </ErrorBoundary>
    </LayoutNormal>
  )
}
