
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
* If the CSV file is missing or invalid, startup will fail.
* For Docker-based runs, keep input files in expected locations unless you update the runtime configuration.

---



## Write-Up

This project follows a ports-and-adapters, or hexagonal, structure. The device domain sits at the center, while HTTP delivery and CSV persistence are implemented as adapters around it. That keeps the core logic focused on the business rules instead of framework or storage details, which is one of the main goals of hexagonal architecture.

I spent approximately 5 hours end-to-end on this project.

A related design choice in the implementation is **dependency injection**.
In practice, the command handler depends on a repository interface defined in `ports.go`, 
not on a concrete CSV repository implementation. 
That makes the application easier to change and easier to test, 
because the business logic depends on abstractions rather than infrastructure details.
This style is closely aligned with dependency composition, where object graphs are assembled at the application boundary rather than hidden inside the objects themselves.)

This also improved the unit testing story. Because the handlers depend on ports rather than concrete implementations, I was able to generate and use mocks for the repository in handler tests.
In other words, the ports pattern made the unit tests more isolated and easier to reason about.

### Runtime Complexity

At startup, the application reads the device definitions from `devices.csv` and inserts them into an in-memory Go map. If there are **N** devices, startup time is **O(N)**, and repository memory usage is also **O(N)**.

For request handling:

* **Device lookup by ID** is average **O(1)** since devices are stored in a Go map keyed by device ID.
* **Adding a heartbeat** is **O(1)**, since it updates only counters and timestamps on the device.
* **Adding upload stats** is **O(1)**, because the implementation maintains only a running sum and count for upload times instead of storing the entire list of samples.
* **Calculating average upload time** is **O(1)** per read, because the average is derived from the stored running sum and count instead of scanning all previous upload samples.
* **Calculating uptime** is **O(1)** in the current design, since it uses already stored timestamps and counters rather than iterating through historical heartbeat events.
