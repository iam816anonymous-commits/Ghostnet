# GhostNet Developer Guide

## Architecture
GhostNet is a hybrid Rust/Go privacy network. Rust is used for the high-performance networking core (`ghostnet-core`), while Go is used for the control plane and orchestration.

## Getting Started
1. Install Go 1.24+ and Rust 1.85+.
2. Explore the core logic in `ghostnet-core/src`.
3. Run the simulation:
   ```bash
   go run simulator/cmd/simulation/main.go -users 100 -duration 10 -mode MEDIUM
   ```

## Contributing
- Follow the modular structure.
- Add unit tests for every new feature.
- Update documentation in the `docs/` directory.
