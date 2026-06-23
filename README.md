# gg

A Go-based mini DNS server. Project built to learn and understand about DNS!

## Project Structure

```
gg/
├── internal/        # App internal logic
|   ├── config/      # Functions to expose env variables
|   ├── model/       # YAML model
|   ├── parser/      # Utilities to parse request
|   ├── server/      # Core logic that handles requests
|   ├── store/       # YAML config
├── data.yaml      # Data yaml
├── .gitignore
├── go.mod
├── main.go
```

## Prerequisites

1. Golang 1.21 or higher
2. Git installed on your machine!
3. `dig` (Domain Information Groper) installed on your system.

## Environment variables

The project uses environment variables to facilitate `port`, `file` and `address` configuration. Example:

```
PORT=8053
FILE=data.yaml
ADDRESS=8.8.8.8:53
```

## Instalation & Usage

1. Clone the repo on your machine.
2. Run `go install`.
3. Start the main server with `air`.
4. In another terminal, send a request:
```
dig "@127.0.0.1" -p 8053 google.com
```
5. You can also use the `+short` extension to only print the address:
```
dig "@127.0.0.1" -p 8053 google.com +short
```

## Using it as the DNS server of your system

To test the server against actual traffic, I configured my machine to use this server as the DNS server of the whole system. Don't worry, it's safe!

<details>
  <summary>Windows</summary>
  
  Go to `Configuration`,  `Network and Internet`, `Wifi` or `Ethernet` (your active connection), `Hardware Properties`, and edit the DNS assignment to `Manual`.
  
  There, you can toggle `IPv4` to `On`, and enter `127.0.0.1`. You can also do the same with `IPv6`, set it to `::1`.

  With that, your windows machine will be using this DNS server!
</details>

<details>
  <summary>Linux</summary>

  (Here, we will focus con `Arch Linux`, as it's the only Linux machine I have.)

  Modify `/etc/systemd/resolved.conf`. Remove the `#` character in the `DNS=` line, and add the IP address of the server.

  ```
  DNS=127.0.0.1:1053 fe80::1
  FallbackDNS=1.1.1.1
  ```

  Then, restart the service:

  ```
  sudo systemctl restart systemd-resolved
  ```

  With that, your (Arch) Linux system will be using this DNS server!

</details>

## Thanks for visiting!!
