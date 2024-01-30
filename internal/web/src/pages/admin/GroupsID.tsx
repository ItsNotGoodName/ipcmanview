import { createAsync } from "@solidjs/router";
import { getAdminGroupsIDPage } from "./GroupsID.data";
import { ErrorBoundary } from "solid-js";
import { PageError } from "~/ui/Page";
import { LayoutNormal } from "~/ui/Layout";

export function AdminGroupsID(props: any) {
  const data = createAsync(() => getAdminGroupsIDPage(BigInt(props.params.id || 0)))

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
