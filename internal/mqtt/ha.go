package mqtt

type HaAvailability struct {
	// The payload that represents the available state.
	PayloadAvailable string `json:"payload_available,omitempty"`
	// The payload that represents the unavailable state.
	PayloadNotAvailable string `json:"payload_not_available,omitempty"`
	// An MQTT topic subscribed to receive availability (online/offline) updates.
	Topic string `json:"topic,omitempty"`
	// Defines a template to extract device’s availability from the topic. To determine the devices’s availability result of this template will be compared to payload_available and payload_not_available.
	ValueTemplate string `json:"value_template,omitempty"`
}

type HaDeviceMap struct {
	// A link to the webpage that can manage the configuration of this device. Can be either an HTTP or HTTPS link.
	ConfigurationUrl string `json:"configuration_url,omitempty"`
	// A list of connections of the device to the outside world as a list of tuples [connection_type, connection_identifier]. For example the MAC address of a network interface: "connections": [["mac", "02:5b:26:a8:dc:12"]].
	Connections [][2]string `json:"connections,omitempty"`
	// The hardware version of the device.
	HwVersion string `json:"hw_version,omitempty"` //
	// A list of IDs that uniquely identify the device. For example a serial number.
	Identifiers []string `json:"identifiers,omitempty"`
	// The manufacturer of the device.
	Manufacturer string `json:"manufacturer,omitempty"`
	// The model of the device.
	Model string `json:"model,omitempty"`
	// The name of the device.
	Name string `json:"name,omitempty"`
	// Suggest an area if the device isn’t in one yet.
	SuggestedArea string `json:"suggested_area,omitempty"`
	// The firmware version of the device.
	SwVersion string `json:"sw_version,omitempty"`
	// Identifier of a device that routes messages between this device and Home Assistant. Examples of such devices are hubs, or parent devices of a sub-device. This is used to show device topology in Home Assistant.
	ViaDevice string `json:"via_device,omitempty"`
}

type HaEntity struct {
	// A list of MQTT topics subscribed to receive availability (online/offline) updates. Must not be used together with availability_topic.
	Availability []HaAvailability `json:"availability,omitempty"`
	// When availability is configured, this controls the conditions needed to set the entity to available. Valid entries are all, any, and latest. If set to all, payload_available must be received on all configured availability topics before the entity is marked as online. If set to any, payload_available must be received on at least one configured availability topic before the entity is marked as online. If set to latest, the last payload_available or payload_not_available received on any configured availability topic controls the availability.
	AvailabilityMode string `json:"availability_mode,omitempty"`
	// Defines a template to extract device’s availability from the availability_topic. To determine the devices’s availability result of this template will be compared to payload_available and payload_not_available.
	AvailabilityTemplate string `json:"availability_template,omitempty"`
	// The MQTT topic subscribed to receive availability (online/offline) updates. Must not be used together with availability.
	AvailabilityTopic string `json:"availability_topic,omitempty"`
	// Information about the device this entity is a part of to tie it into the device registry. Only works when unique_id is set. At least one of identifiers or connections must be present to identify the device.
	Device HaDeviceMap `json:"device,omitempty"`
	// The type/class of the entity to set the icon in the frontend. The device_class can be null.
	DeviceClass string `json:"device_class,omitempty"`
	// Flag which defines if the entity should be enabled when first added.
	EnabledByDefault bool `json:"enabled_by_default,omitempty"`
	// The encoding of the published messages.
	Encoding string `json:"encoding,omitempty"`
	// The category of the entity.
	EntityCategory string `json:"entity_category,omitempty"` // (Optional, default: None)
	// Icon for the entity.
	Icon string `json:"icon,omitempty"`
	// Defines a template to extract the JSON dictionary from messages received on the json_attributes_topic. Usage example can be found in MQTT sensor documentation.
	JsonAttributesTemplate string `json:"json_attributes_template,omitempty"`
	// The MQTT topic subscribed to receive a JSON dictionary payload and then set as sensor attributes. Usage example can be found in MQTT sensor documentation.
	JsonAttributesTopic string `json:"json_attributes_topic,omitempty"`
	// The name to use when displaying this entity.
	Name string `json:"name,omitempty"`
	// Used instead of name for automatic generation of entity_id
	ObjectId string `json:"object_id,omitempty"`
	// The payload that represents the available state.
	PayloadAvailable string `json:"payload_available,omitempty"`
	// The payload that represents the unavailable state.
	PayloadNotAvailable string `json:"payload_not_available,omitempty"`
	// The maximum QoS level to be used when receiving and publishing messages.
	Qos int `json:"qos,omitempty"`
	// An ID that uniquely identifies this entity. If two entities have the same unique ID, Home Assistant will raise an exception.
	UniqueId string `json:"unique_id,omitempty"`
}

func NewHaEntity(conn Conn) HaEntity {
	return HaEntity{
		Availability: []HaAvailability{
			{Topic: TopicServerState(conn.Topic)},
		},
		AvailabilityMode: "all",
	}
}

// https://www.home-assistant.io/integrations/event.mqtt/
type HaEvent struct {
	HaEntity
	// A list of valid event_type strings.
	EventTypes []string `json:"event_types,omitempty"`
	// The MQTT topic subscribed to receive JSON event payloads. The JSON payload should contain the event_type element. The event type should be one of the configured event_types.
	StateTopic string `json:"state_topic,omitempty"`
}

// https://www.home-assistant.io/integrations/binary_sensor.mqtt/
type HaBinarySensor struct {
	HaEntity
	// The MQTT topic subscribed to receive sensor’s state.
	StateTopic string `json:"state_topic,omitempty"`
}
