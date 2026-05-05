
# SafelyYou Fleet Monitoring Submission

A small Go service for ingesting device telemetry and exposing per-device statistics.

It reads a `devices.csv` file on startup, accepts heartbeats and upload stats for each device, and returns:

- average upload time
- uptime percentage

---

## Project Structure

```text
.
├── main.go
├── results.txt
├── go.mod
├── .gitignore
└── src
    ├── app.go
    ├── device
    │   ├── device_entity.go
    │   ├── ports.go
    │   └── handlers
    │       ├── handler.go
    │       └── handler_test.go
    └── infrastructure
        ├── controller.go
        ├── repository.go
        └── mocks

```
### File overview
The projects follow object-oriented ports and adaptors (hexagonal) architecture structure.
* `main.go`
  Application entry point. Reads CLI flags and starts the app.

* `src/app.go`
  App bootstrap and runtime dependency injection.

* `src/device/device_entity.go`
  Core device domain logic.

* `src/device/ports.go`
  Interfaces used by the domain and handlers.

* `src/device/handlers/handler.go`
  Command handler logic for device events.

* `src/device/handlers/handler_test.go`
  Unit tests for the handler layer.

* `src/infrastructure/controller.go`
  HTTP controller and route registration.

* `src/infrastructure/repository.go`
  CSV-backed device repository.

* `src/infrastructure/mocks`
  Generated test mocks.

---


## Configuration

The application supports the following flags:

* `-h`, `--host`
  Host and port to bind the server to

* `-f`, `--file`
  Path to the `devices.csv` file

### Defaults

If you do not provide flags, the application uses:

* host: `0.0.0.0:6733`
* file: `./devices.csv`

---

## Run Locally with Go

Make sure `devices.csv` is available.

Example:

```bash
go mod tidy
go run ./main.go
```

Run with custom flags:

```bash
go run ./main.go --host 0.0.0.0:6733 --file ./devices.csv
```

---

## Run with Dockerfile

This project includes a two-stage `Dockerfile`.

### Important

If you want to run with the provided `Dockerfile`, place `devices.csv` alongside the code at the project root, next to:

* `main.go`
* `go.mod`
* `Dockerfile`

Example expected layout:

```text
.
├── Dockerfile
├── main.go
├── go.mod
├── devices.csv
└── src
```

### Build the image

```bash
docker build -t fleet-app .
```

### Run the container

```bash
docker run --rm -p 6733:6733 fleet-app
```

---

## Run End-to-End Test with the Device Simulator

To run the application together with the device simulator and generate `results.txt`, use Docker Compose.

### Before running

Make sure the simulator binary filename in `docker-compose.yml` exactly matches the simulator file you downloaded.

For example, if your downloaded simulator file is:

```text
device-simulator-linux-arm64
```

then your compose command should reference that exact file name.

### Example

```bash
docker compose up --build
```

This will:

* build and start the app
* start the simulator
* run the simulator against the app
* generate `results.txt` if your compose setup writes simulator output to that file

### Notes

* verify the simulator binary name in `docker-compose.yml`
* verify the simulator command-line arguments match the simulator you downloaded
* verify `devices.csv` is present

---

## Notes

* If the CSV file is missing or invalid, startup will fail.
* For Docker-based runs, keep input files in expected locations unless you update the runtime configuration.

---