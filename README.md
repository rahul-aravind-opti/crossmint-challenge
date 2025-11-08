# Crossmint Megaverse Challenge

## Overview
This repository contains the solution for the Crossmint Megaverse coding challenge. The CLI builds two astral layouts:
- Phase 1 renders an X-shaped constellation of 🪐 POLYanets between rows 2 and 8.
- Phase 2 reproduces the Crossmint logo using 🌙 SOLoons, ☄ comETHs, and 🪐 POLYanets sourced from the goal map endpoint.

The implementation focuses on reliability under tight API rate limits. It includes resiliency tooling, clear domain modeling, and repeatable workflows for initialization and execution.

## Prerequisites
- Go 1.21+
- Internet access to https://challenge.crossmint.io/api
- Candidate ID from Crossmint (use `megaverse init --candidate <id>`).

## Quick Start
1. Clone the repository and change into the project directory.
2. Copy `config/config.yaml` if you need custom settings.
3. Provide your candidate ID via `megaverse init --candidate <id>` or set `CROSSMINT_CANDIDATE_ID`.
4. Run a phase command:
   - `megaverse phase1` generates the Phase 1 cross pattern.
   - `megaverse phase2` builds the Crossmint logo using the goal map.
5. Inspect the live map with `megaverse status` or the Crossmint dashboard as needed.

All commands respect the configured timeout, rate limit, and retry budget to stay within the API allowances.

## CLI Commands
- `megaverse init` creates or updates the configuration file with your candidate ID.
- `megaverse phase1` runs the cross-pattern strategy in parallel workers.
- `megaverse phase2` downloads the goal map, plans the layout, and materialises it in parallel.
- `megaverse status` prints a summary of the current megaverse grid.

## Architecture Highlights
- `cmd/megaverse`: program entry point wiring configuration, services, and CLI.
- `internal/interfaces/cli`: Cobra commands orchestrating user actions and timeouts.
- `internal/application`: application services plus strategy pattern implementations for each phase.
- `internal/domain/entities`: core entities (e.g., `Polyanet`, `Soloon`, `Megaverse`) with validation.
- `internal/infrastructure/api`: HTTP client with rate limiting, exponential backoff, and retry-go integration.
- `pkg/ratelimit`: thin wrapper around `golang.org/x/time/rate` for shared limiter usage.
- `pkg/retry`: adapter around `github.com/avast/retry-go/v4` exposing a challenge-friendly configuration.

## Resiliency Tooling
- Rate limiting is enforced before every HTTP call to avoid 429 responses.
- Retry policies support exponential backoff with bounded delays and context cancellation.
- Creation strategies aggregate errors so partial failures are surfaced without aborting the whole run.

## Configuration
`config/config.yaml` exposes sane defaults. Key sections include:
- `api.base_url`, `api.timeout`, and `api.candidate_id`
- `api.retry` (attempts, delays, multiplier)
- `api.rate_limit.requests_per_second`
- `execution.max_workers`, `execution.batch_size`, `execution.timeout`

Environment variables compatible with Viper (e.g., `CROSSMINT_API_TIMEOUT`) override file values at runtime.

## Testing
Run all unit tests with:
```
go test ./...
```

The suite covers retry behaviour, strategy plan generation, and infrastructure helpers. Extend or focus tests by targeting individual packages, for example `go test ./internal/application/strategies`.

## Troubleshooting
- **429 Too Many Requests**: lower `api.rate_limit.requests_per_second` or increase retry attempts.
- **Timeouts**: raise `execution.timeout` or check network connectivity.
- **Invalid map state**: run `megaverse status` to inspect current grid contents before re-running a phase.

Happy minting! 🚀
