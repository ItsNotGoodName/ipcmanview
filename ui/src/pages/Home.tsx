import { Component, ErrorBoundary, Match, Switch, createResource, } from "solid-js";
import { useService } from "~/providers/service";
import { LayoutDefault } from "~/ui/Layout";

export const Home: Component = () => {
  const { dahuaService } = useService()

  const [cameraCount, data] = createResource(() => dahuaService.cameraCount().then((res) => res.count))

  return (
    <LayoutDefault>
      <ErrorBoundary
        fallback={(err) => <div>Error: {err.toString()}</div>}
      >
        <Switch>
          <Match when={cameraCount.state == "pending"}>
            I am loading...
          </Match>
          <Match when={cameraCount.state == "errored"}>
            <div onClick={data.refetch}>
              This is an error = {cameraCount.error.toString()}
            </div>
          </Match>
          <Match when={cameraCount()}>
            There are {cameraCount()} cameras.
          </Match>
        </Switch>
      </ErrorBoundary>
    </LayoutDefault>
  )
};

