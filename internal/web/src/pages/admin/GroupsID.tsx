import { createAsync } from "@solidjs/router";
import { AdminGroupsIDPageSearchParams, getAdminGroupsIDPage } from "./GroupsID.data";
import { ErrorBoundary } from "solid-js";
import { PageError } from "~/ui/Page";
import { PageProps } from "~/lib/utils";

export function AdminGroupsID(props: PageProps<AdminGroupsIDPageSearchParams>) {
  const data = createAsync(() => getAdminGroupsIDPage({ id: BigInt(props.params.id || 0) }))

  return (
    <div class="flex justify-center p-4">
      <div class="flex w-full max-w-4xl flex-col gap-2">
        <ErrorBoundary fallback={(e: Error) => <PageError error={e} />}>
          <pre>
            {JSON.stringify(data(), null, 2)}
          </pre>
        </ErrorBoundary>
      </div>
    </div>
  )
}
