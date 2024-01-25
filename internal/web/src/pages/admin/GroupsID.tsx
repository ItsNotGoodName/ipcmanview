import { createAsync } from "@solidjs/router";
import { AdminGroupsIDPageSearchParams, getAdminGroupsIDPage } from "./GroupsID.data";
import { ErrorBoundary } from "solid-js";
import { PageError } from "~/ui/Page";
import { PageProps } from "~/lib/utils";
import { LayoutNormal } from "~/ui/Layout";

export function AdminGroupsID(props: PageProps<AdminGroupsIDPageSearchParams>) {
  const data = createAsync(() => getAdminGroupsIDPage({ id: BigInt(props.params.id || 0) }))

  return (
    <LayoutNormal>
      <ErrorBoundary fallback={(e: Error) => <PageError error={e} />}>
        <pre>
          {JSON.stringify(data(), null, 2)}
        </pre>
      </ErrorBoundary>
    </LayoutNormal>
  )
}
