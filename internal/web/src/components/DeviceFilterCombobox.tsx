import { createAsync, } from "@solidjs/router"
import { RiSystemFilterLine } from "solid-icons/ri"
import { getListDevices } from "~/pages/data"
import { ListDevicesResp_Device } from "~/twirp/rpc"
import { ComboboxContent, ComboboxControl, ComboboxIcon, ComboboxInput, ComboboxItem, ComboboxItemLabel, ComboboxListbox, ComboboxReset, ComboboxRoot, ComboboxTrigger, ComboboxState } from "~/ui/Combobox"

export function DeviceFilterCombobox(props: { setDeviceIDs: (ids: string[]) => void, deviceIDs: string[] }) {
  const listDevices = createAsync(() => getListDevices())

  return (
    <ComboboxRoot<ListDevicesResp_Device>
      multiple
      optionValue="id"
      optionTextValue="name"
      optionLabel="name"
      options={listDevices() || []}
      placeholder="Device"
      value={props.deviceIDs.map(v => ({ id: v, name: "" }))}
      onChange={(value) => props.setDeviceIDs(value.map(v => v.id))}
      itemComponent={props => (
        <ComboboxItem item={props.item}>
          <ComboboxItemLabel>{props.item.rawValue.name}</ComboboxItemLabel>
        </ComboboxItem>
      )}
    >
      <ComboboxControl<ListDevicesResp_Device> aria-label="Device">
        {state => (
          <ComboboxTrigger>
            <ComboboxIcon as={RiSystemFilterLine} class="size-4" />
            Device
            <ComboboxState state={state} getOptionString={(option) => option.name} />
            <ComboboxReset state={state} class="size-4" />
          </ComboboxTrigger>
        )}
      </ComboboxControl>
      <ComboboxContent>
        <ComboboxInput />
        <ComboboxListbox />
      </ComboboxContent>
    </ComboboxRoot >
  )
}
