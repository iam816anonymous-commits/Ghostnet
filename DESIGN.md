# GhostNet Privacy Network vNext - Design Documentation

## Architecture Overview

GhostNet vNext is a modular privacy network designed to resist traffic analysis through a multi-stage pipeline:

1.  **Identity Manager**: Ensures per-site isolation of cookies and local storage.
2.  **Traffic Classifier**: Categorizes traffic to allow for specialized queue management.
3.  **Mix Queue Manager**: Batches packets, reorders them, and releases them at random intervals.
4.  **Cover Traffic Engine**: Generates adaptive dummy traffic to obfuscate real communication patterns.
5.  **Packet Normalizer**: Standardizes packet sizes via padding and fragmentation.
6.  **Multi-Hop Routing Engine**: Selects random paths through a network of relays.

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
- **Privacy Benefit**: Maintains a constant or semi-constant traffic profile, making it difficult for an observer to detect when real communication is occurring.
- **Trade-off**: Consumes additional bandwidth and battery.

### 4. Packet Normalizer
- **Function**: Forces all packets to a fixed size (1024 bytes).
- **Privacy Benefit**: Eliminates packet-size fingerprinting.
- **Trade-off**: Adds overhead due to padding and fragmentation headers.

### 5. Multi-Hop Routing Engine
- **Function**: 4-hop routing (Entry -> Middle -> Middle -> Exit).
- **Privacy Benefit**: Ensures no single relay knows both the source and the destination.
- **Trade-off**: Increased latency and higher risk of path failure.

## Threat Model & Limitations

- **Traffic Correlation**: Resistant to global passive observers through mixing and cover traffic.
- **Timing Analysis**: Mitigated by Mix Queues.
- **Packet-Size Analysis**: Neutralized by Packet Normalization.
- **Limitations**: We do not claim complete anonymity. A powerful adversary controlling both Entry and Exit nodes may still perform correlation, although mixing makes this significantly harder.

## Performance vs. Privacy

| Mode | Privacy Score | Latency | Bandwidth Overhead |
| :--- | :--- | :--- | :--- |
| LOW | 3-5 | Low | Minimal |
| MEDIUM | 5-7 | Moderate | Low |
| HIGH | 7-9 | High | Moderate |
| MAXIMUM | 9-10 | Very High | High |
