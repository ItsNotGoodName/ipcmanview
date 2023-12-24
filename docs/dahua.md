# Event Codes

| Code                  | Action(s) | Description                                     |
| --------------------- | --------- | ----------------------------------------------- |
| NTPAdjustTime         | Pulse     | Fires when time changes via NTP.                |
| TimeChange            | Pulse     | Fires when time changes.                        |
| InterVideoAccess      | Pulse     | Fires when there is a login.                    |
| StorageChange         | Pulse     | Fires when a file is read.                      |
| NewFile               |           | Fires when a file is created.                   |
| IntelliFrame          | Pulse     | Fires when there is detection.                  |
| RtspSessionDisconnect | Start     | Fires when a RTSP client disconnects.           |
| SystemState           | Pulse     | Fires probably when device ready state changes. |

# Recommended Event Rules

- NTPAdjustTime
- TimeChange
- InterVideoAccess
- StorageChange
- NewFile
