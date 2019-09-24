# things with buzzers: websocket

A WebSocket server to publish messages when someone pushed a hardware game show buzzer.

## Compile for Raspberry Pi

```sh
GOARM=7 GOARCH=arm GOOS=linux go build -o twb-websocket
```