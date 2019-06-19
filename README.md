# go-clicker
Golang UI clicker remote controller via web interface

## Windows vendored binaries

Windows users, download the [zip file][1] if wish to vendor the ngrok.exe and libwinpthread-1.dll on your system (or not found in PATH).

[1]: https://github.com/asfaltboy/go-clicker/releases/download/0.1/go-clicker-windows.zip

## Ngrok Auth

First-time users should authenticate to ngrok to use the TCP tunnel feature. Sign up to, and browse to https://dashboard.ngrok.com/auth to get the token, then simply run:

```
ngrok authtoken <token>
```

## Running

By default running the `main.exe` will start listening on localhost:8088 for both HTTP and WebSocket connections. We will also target this host in the served HTML as the WebSocket connection target.

```
main.exe
```

To start with a different host:port target, use:

```
main.exe -host
```

It may be useful to use ngrok to bypass firewall and network routing. To start with ngrok, simply run

```
main.exe -ngrok
```
