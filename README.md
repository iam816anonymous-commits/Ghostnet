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
- `pkg/classifier`: Traffic classification logic.
- `pkg/common`: Shared types and constants.
- `pkg/cover`: Adaptive dummy traffic generation.
- `pkg/dashboard`: Privacy metrics and visualization.
- `pkg/identity`: Per-site identity isolation.
- `pkg/normalizer`: Packet padding and fragmentation.
- `pkg/pipeline`: The main processing pipeline.
- `pkg/queue`: Batching and mixing engine.
- `pkg/routing`: Multi-hop path selection and relay management.
