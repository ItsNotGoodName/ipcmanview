# Events

Events are produced by the device and are accessed through the `eventManager` CGI API.

| Code                  | Action(s) | Description                                     |
| --------------------- | --------- | ----------------------------------------------- |
| NTPAdjustTime         | Pulse     | Fires when time changes via NTP.                |
| TimeChange            | Pulse     | Fires when time changes.                        |
| InterVideoAccess      | Pulse     | Fires when there is a login.                    |
| StorageChange         | Pulse     | Fires when a file is read.                      |
| NewFile               | Pulse     | Fires when a file is created.                   |
| IntelliFrame          | Pulse     | Fires when there is (video) detection.          |
| RtspSessionDisconnect | Start     | Fires when a RTSP client disconnects.           |
| SystemState           | Pulse     | Fires probably when device ready state changes. |

## Event Rules

Event rules can be used to control the flow of events to consumers.

https://github.com/ItsNotGoodName/ipcmanview/assets/35015993/01704920-ada3-4a5c-87a5-ebd51c6f18cc

The following events can be considerably noisy.

- NTPAdjustTime
- TimeChange
- StorageChange
- InterVideoAccess
- NewFile
- RtspSessionDisconnect

# Bugs

## First pictures not saved randomly

| Model          | Build Date |
| -------------- | ---------- |
| IPC-T5442TM-AS | 2020-10-19 |
| IPC-T5442TM-AS | 2020-12-03 |

This bug will cause the camera to randomly not save the first picture when a `VideoAnalyseRule` event is triggered.
This was found when I compared the pictures stored on the camera with pictures from email notifications.
I have yet to find a way to reproduce it successfully and as far as I know, this does not affect videos.
