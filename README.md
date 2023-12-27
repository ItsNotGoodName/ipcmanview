# IPCManView

**🚧 WORK IN PROGRESS, EACH COMMIT MOST LIKELY BREAKS THE LAST 🚧**

Application for managing and viewing Dahua devices.

TODO: explain why this program exists

https://github.com/ItsNotGoodName/ipcmanview/assets/35015993/f2c73450-ba61-4d39-b3c5-92eaa719fced

# Features

- Single binary
- View device information (e.g. software version, license, storage, …)
- Subscribe to device events
- Show snapshot of cameras
- Publish to MQTT with Home Assistant MQTT discovery
- View files on devices (only local, SFTP, and FTP storage locations are supported)

# Usage

```
ipcmanview serve
```

## Configuration

| Environment Variable | Default           | Description                                           |
| -------------------- | ----------------- | ----------------------------------------------------- |
| `DIR`                | "ipcmanview_data" | Directory path for storing data.                      |
| `HTTP_HOST`          | ""                | HTTP host to listen on (e.g. "127.0.0.1").            |
| `HTTP_PORT`          | 8080              | HTTP port to listen on.                               |
| `MQTT_ADDRESS`       | ""                | MQTT server address (e.g. "mqtt://example.com:1883"). |
| `MQTT_TOPIC`         | "ipcmanview"      | MQTT server topic to publish messages.                |
| `MQTT_USERNAME`      | ""                | MQTT server username for authentication.              |
| `MQTT_PASSWORD`      | ""                | MQTT server password for authentication.              |
| `MQTT_HA`            | false             | Enable Home Assistant MQTT discovery.                 |
| `MQTT_HA_TOPIC`      | "homeassistant"   | Home Assistant MQTT discover topic.                   |

# Roadmap

Roadmap is in order of importance.

- Support syncing the config `VideoInMode` (Camera > Conditions > Profile Management) with sunrise and sunset
- Support editing the config `General` (System > General)
- Support editing the config `Email` (Network > SMTP(Email))
- Support editing the config `VideoAnalyseRule`
- View live stream on cameras
- View DAV files in local storage via RTSP (see 4.1.3 in the Dahua HTTP API PDF)
- Create and cache thumbnails for files
- Act as a HomeKit bridge for viewing cameras
- Support two-way talk on cameras that support it (see [./pkg/dahuacgi/audio.go](./pkg/dahuacgi/audio.go))
- Support OpenAPI (but I don't want to write more YAML then there is code in the handlers)
