# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Overview

Fence is a [Cedar](https://www.cedarpolicy.com/) (`cedar-go`) based authorization system written in Go (module `github.com/binarymatt/fence`). It exposes a Connect-RPC server (`cmd/agent`) that stores Cedar policies and entities and answers `IsAllowed` requests, plus a Go client library (`pkg/client`) that can evaluate authz either locally or by calling the server.

Planned run modes (from README): **local** (policies/entities from local files via client), **cached local** (local files refreshed on a goroutine), **proxied** (client calls server for every decision), **cached** (client/agent caches server state and refreshes periodically).

## Tooling & Commands

Tools are pinned via devbox (`devbox.json`) and loaded with direnv (`.envrc`); task runner is `go-task` (`Taskfile.yml`).

- Build: `go build ./...`
- Run server: `go run ./cmd/agent -c sample_config.yaml` (config path also via `FENCE_CONFIG`)
- Test all: `task go:test` (= `go test ./... -cover`)
- Single test: `go test ./internal/service -run TestGetPolicySet -v`
- Regenerate protobuf/Connect code: `task proto` (`buf lint && buf generate`) — output goes to `gen/`, never edit `gen/` by hand
- Regenerate mocks: `task go:mocks` (mockery; writes `mocks.go` into `pkg/providers` and `internal/service` for all interfaces — don't hand-edit those)
- Container/deploy: `task container:build`, `task infra:plan`, `task deploy` (podman + terraform against `registry.homelab.lan`; Taskfile references a `terraform/` dir, but the terraform files currently live in `examples/terraform/`)

Badger builds require CGO (see `Dockerfile`).

## Architecture

**API definition** — `proto/fence/v1/`:
- `fence.proto`: `FenceService.IsAllowed` plus shared types (`UID`, `Entity`, `Policy`, `Diagnostics`). Uses protovalidate annotations.
- `admin.proto`: `FenceAdminService` — CRUD for entities and policies.

**Server path** (`cmd/agent` → `pkg/server` → `internal/service`):
- `pkg/server/config.go`: koanf YAML config (`listen_address`, `db_path`, `db_type`, `log_level`, `proxy.{address,token}`) with defaults.
- `pkg/server/server.go`: picks a `DataStore` by `db_type` (`badger` or `sql` = SQLite via bun; SQL tables are created on startup), mounts both Connect handlers on one mux with logging + CORS middleware, `/health`, HTTP/1 + h2c.
- `internal/service.Service` implements both `FenceServiceHandler` and `FenceAdminServiceHandler`. On every `IsAllowed` it loads the full policy set and entity map from the `DataStore` and calls `cedar.Authorize`.
- `internal/service.DataStore` is an interface with **unexported** methods, implemented by `sqlDataStore` (`db.go`, bun models in `models.go`) and `badgerDataStore` (`badger.go`, proto-marshalled values under `policy_<id>` / `entity_<type>_<id>` keys).

**Client path** (`pkg/client` → `pkg/providers`):
- `providers.FenceProvider` interface: `IsAllowed(ctx, principal, action, resource)` + `Refresh(ctx)`.
  - `FileProvider`: loads a Cedar policy file and JSON entity file through `afero.Fs` (use `afero.NewMemMapFs` in tests).
  - `RemoteServerProvider`: Connect client to `FenceService` with bearer-token + retryablehttp transport.
- `pkg/client.Client` wraps a provider; `ContextWithClient` / `ClientFromContext` stash it in a context.

**Translation** — `internal/translation` converts Cedar `Decision`/`Diagnostic` into `fencev1.IsAllowedResponse` (empty reasons + deny ⇒ synthetic "default deny" reason). Shared by the service and `FileProvider`.

## Testing notes

- Assertions use `github.com/shoenig/test/must` (testify only for generated mocks).
- SQL datastore tests use in-memory SQLite via `setupDB(t, loadFixtures)` in `internal/service/service_test.go`, loading `testdata/fixture.yaml` with bun `dbfixture`.

## Current state

Work in progress: `FenceProvider.IsAllowed` takes no Cedar context, so `pkg/client.Client.Authorize` drops `req.Context`, and since `Authorize` takes no `context.Context` it calls the provider with `context.Background()`.
