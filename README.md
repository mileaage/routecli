# RouteCLI

Simple CLI for building HTTP servers with configurable routes and middleware.

## Install

```bash
go install github.com/mileaage/routecli@latest
```

## Usage

```bash
./routecli -start        # Start server on :8080
./routecli -help         # Show help
```

## Config

Create `config.yaml`:

```yaml
routes:
  - "/"
  - "/api/health" 
  - "/api/users"

middlewares:
  - "logger"
```

## Build

```bash
git clone https://github.com/mileaage/routecli.git
cd routecli
go build
```

Requires Go 1.19+