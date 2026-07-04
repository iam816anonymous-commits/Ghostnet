# GhostNet Privacy Network vNext

GhostNet is a next-generation privacy network focusing on traffic analysis resistance.

## Getting Started

### Prerequisites
- Go 1.24+

### Running the Simulation
To run a network simulation:
```bash
go run cmd/simulation/main.go -users 100 -duration 10 -mode MEDIUM
```

### Running Tests
```bash
go test ./...
```

## Project Structure

### GhostNet Core (Rust) - `ghostnet-core/`
- `src/classifier`: Traffic classification and dynamic queue management.
- `src/cover`: Adaptive dummy traffic generation.
- `src/normalizer`: Packet padding and fragmentation.
- `src/queue`: Mix networking, batching, and shuffling.
- `src/routing`: Multi-hop path selection.

### GhostNet Control Plane (Go) - `pkg/`
- `pkg/dashboard`: Privacy metrics and visualization.
- `pkg/identity`: Per-site identity isolation.
- `pkg/pipeline`: Orchestration pipeline.
