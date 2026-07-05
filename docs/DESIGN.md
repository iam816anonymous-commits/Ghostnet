# GhostNet Privacy Network vNext - Design Documentation

## Architecture Overview

GhostNet vNext is a modular privacy network designed to resist traffic analysis through a multi-stage pipeline. The architecture is a hybrid of **Rust** for performance-critical packet processing and **Go** for orchestration and control plane.

### 1. GhostNet Core (Rust)
The core networking engine is implemented in Rust to ensure memory safety, high performance, and predictable resource usage.
- **Packet Normalizer**: Standardizes packet sizes (1024 bytes) via padding and fragmentation.
- **Mix Queue Manager**: Handles batching, shuffling, and random release intervals.
- **Cover Traffic Engine**: Generates adaptive dummy traffic.
- **Traffic Classifier**: Categorizes traffic and manages dynamic queues.
- **Routing Engine**: Selects multi-hop paths (4 hops).

### 2. Control Plane & Orchestration (Go)
Go is used for high-level management and simulation.
- **Identity Manager**: Ensures per-site isolation.
- **Simulation Suite**: Orchestrates thousands of concurrent users.
- **Privacy Dashboard**: Reports metrics and privacy scores.

## Component Details & Privacy Benefits

### 1. Adaptive Traffic Classifier
- **Function**: Automatically classifies traffic (e.g., Small Web, API, Video).
- **Privacy Benefit**: Allows the system to apply different mixing strategies based on traffic sensitivity and volume.
- **Trade-off**: Over-classification can reduce the anonymity set if too many small queues are created.

### 2. Mix Queue Manager
- **Function**: Implements batching (Min Batch Size) and random release intervals.
- **Privacy Benefit**: Prevents timing correlation attacks by breaking the 1:1 temporal relationship between incoming and outgoing packets.
- **Trade-off**: Increases latency, especially when traffic volume is low.

### 3. Adaptive Cover Traffic Engine
- **Function**: Generates dummy packets during idle periods or low-volume bursts.
- **Privacy Benefit**: Maintains a constant or semi-constant traffic profile.
- **Trade-off**: Consumes additional bandwidth and battery.

### 4. Packet Normalizer
- **Function**: Forces all packets to a fixed size (1024 bytes).
- **Privacy Benefit**: Eliminates packet-size fingerprinting.
- **Trade-off**: Adds overhead due to padding and fragmentation headers.

## Threat Model & Limitations

- **Traffic Correlation**: Resistant to global passive observers through mixing and cover traffic.
- **Timing Analysis**: Mitigated by Mix Queues.
- **Packet-Size Analysis**: Neutralized by Packet Normalization.
- **Limitations**: We do not claim complete anonymity. Endpoint compromise or controlled entry/exit nodes remain risks.

## Performance vs. Privacy

| Mode | Privacy Score | Latency | Bandwidth Overhead |
| :--- | :--- | :--- | :--- |
| LOW | 3-5 | Low | Minimal |
| MEDIUM | 5-7 | Moderate | Low |
| HIGH | 7-9 | High | Moderate |
| MAXIMUM | 9-10 | Very High | High |
