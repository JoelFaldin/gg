# gg

A Go-based mini DNS server. Project built to learn and understand about DNS!

## Project Structure

gg/
├── internal/        # App internal logic
|   ├── model/       # YAML model
|   ├── parser/      # Utilities to parse request
|   ├── server/      # Core logic that handles requests
|   ├── store/       # YAML config
├── data.yaml      # Data yaml
├── .gitignore
├── go.mod
├── main.go

## Prerequisites

1. Golang 1.21 or higher
2. Git installed on your machine!
3. `dig` (Domain Information Groper) installed on your system.

## Instalation

1. Clone the repo on your machine.
2. Run `go install`.
3. Start the main server with `air`.
4. In another terminal, send a request:
```
dig "@127.0.0.1" -p 8053 google.com
```

## Thanks for visiting!!
