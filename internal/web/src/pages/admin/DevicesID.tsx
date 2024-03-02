import { A, createAsync } from "@solidjs/router";
import { ErrorBoundary } from "solid-js";
import { PageError } from "~/ui/Page";
import { LayoutNormal } from "~/ui/Layout";
import { getAdminDevicesIDPage } from "./DevicesID.data";
import { Shared } from "~/components/Shared";
import { BreadcrumbsItem, BreadcrumbsLink, BreadcrumbsRoot, BreadcrumbsSeparator } from "~/ui/Breadcrumbs";

export function AdminDevicesID(props: any) {
  const id = () => props.params.id || 0
  const data = createAsync(() => getAdminDevicesIDPage(id()))

  return (
    <LayoutNormal class="max-w-4xl">
      <Shared.Title>
        <BreadcrumbsRoot>
          <BreadcrumbsItem>
            <BreadcrumbsLink as={A} href="../">
              Devices
            </BreadcrumbsLink>
            <BreadcrumbsSeparator />
          </BreadcrumbsItem>
          <BreadcrumbsItem>
            {id()}
          </BreadcrumbsItem>
        </BreadcrumbsRoot>
      </Shared.Title>
      <ErrorBoundary fallback={(e) => <PageError error={e} />}>
        <pre>
          {JSON.stringify(data(), null, 2)}
        </pre>
      </ErrorBoundary>
    </LayoutNormal>
  )
}

export default AdminDevicesID
