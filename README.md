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

## Thanks for visiting!!
