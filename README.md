# IPCManView

**🚧 WORK IN PROGRESS, EACH COMMIT MOST LIKELY BREAKS THE LAST 🚧**

Application for managing and viewing Dahua devices.

TODO: explain why this program exists

https://github.com/ItsNotGoodName/ipcmanview/assets/35015993/f2c73450-ba61-4d39-b3c5-92eaa719fced

# Features

- Single binary<sub>1</sub>
- View device information (e.g. software version, license, storage, …)
- Subscribe to device events
- View live stream of cameras
- View snapshot of cameras
- Publish to MQTT with Home Assistant MQTT discovery
- View files on devices (only local, SFTP, and FTP storage locations are supported)

1. Streaming requires [MediaMTX](https://github.com/bluenviron/mediamtx), and [MQTT](https://mqtt.org/) requires a [MQTT broker](https://mosquitto.org/).

# Usage

```
ipcmanview serve
```

## Configuration

| Environment Variable       | Default           | Description                                                                                                                                   |
| -------------------------- | ----------------- | --------------------------------------------------------------------------------------------------------------------------------------------- |
| `DIR`                      | "ipcmanview_data" | Directory path for storing data.                                                                                                              |
| `HTTP_HOST`                |                   | HTTP host to listen on (e.g. "127.0.0.1").                                                                                                    |
| `HTTP_PORT`                | 8080              | HTTP port to listen on.                                                                                                                       |
| `MQTT_ADDRESS`             |                   | MQTT server address (e.g. "mqtt://192.168.1.20:1883").                                                                                        |
| `MQTT_TOPIC`               | "ipcmanview"      | MQTT server topic to publish messages.                                                                                                        |
| `MQTT_USERNAME`            |                   | MQTT server username for authentication.                                                                                                      |
| `MQTT_PASSWORD`            |                   | MQTT server password for authentication.                                                                                                      |
| `MQTT_HA`                  | false             | Enable Home Assistant MQTT discovery.                                                                                                         |
| `MQTT_HA_TOPIC`            | "homeassistant"   | Home Assistant MQTT discover topic.                                                                                                           |
| `MEDIAMTX_HOST`            |                   | MediaMTX host address (e.g. "192.168.1.20").                                                                                                  |
| `MEDIAMTX_WEBRTC_PORT`     | 8889              | MediaMTX WebRTC port.                                                                                                                         |
| `MEDIAMTX_HLS_PORT`        | 8888              | MediaMTX HLS port.                                                                                                                            |
| `MEDIAMTX_PATH_TEMPLATE`   |                   | [Template](https://pkg.go.dev/text/template) for generating MediaMTX paths (e.g. `ipcmanview_dahua_{{.DeviceID}}_{{.Channel}}_{{.Subtype}}`). |
| `MEDIAMTX_STREAM_PROTOCOL` | "webrtc"          | MediaMTX stream protocol<sub>1</sub> ("webrtc" or "hls").                                                                                     |

1. See MediaMTX [docs](https://github.com/bluenviron/mediamtx#web-browsers-1) on whether to use WebRTC or HLS.

# Roadmap

Roadmap is in order of importance.

- Support syncing the config `VideoInMode` (Camera > Conditions > Profile Management) with sunrise and sunset
- Support editing the config `General` (System > General)
- Support editing the config `Email` (Network > SMTP(Email))
- Support editing the config `VideoAnalyseRule`
- View DAV files in local storage via RTSP (see 4.1.3 in the Dahua HTTP API PDF)
- Create and cache thumbnails for files
- Act as a HomeKit bridge for viewing cameras
- Support two-way talk on cameras that support it (see [./pkg/dahuacgi/audio.go](./pkg/dahuacgi/audio.go))
- Support OpenAPI (but I don't want to write more YAML then there is code in the handlers)
