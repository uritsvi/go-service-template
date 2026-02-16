# Go REST API Template

Production-ready, feature-driven REST API template in Go. Uses Gin, env-based configuration, a logger interface, and Swagger/OpenAPI. Optimized for AI-assisted development and extension.

**Renaming the project:** To use this template under a different project name, run the following from the repo root (replace `YOUR-PROJECT-NAME` with your desired module name, e.g. `my-api`). On macOS use `sed -i ''`; on Linux use `sed -i`:

```bash
grep -rl 'go-service-template' --include='*.go' --include='go.mod' --include='*.yaml' --include='*.yml' . 2>/dev/null | xargs sed -i 's/go-service-template/YOUR-PROJECT-NAME/g'
```

Then run `go mod tidy`. Optionally rename the repository directory to match.

**Clone and run without committing:** You can use this template without creating a new repo or committing. Clone the repository, copy `.example.env` to `.env`, set `SERVER_PORT`, run `go mod tidy` and `go run ./cmd/server`. The project runs as-is with no git push or new remote. Use it to try the template locally or as a starting point; when you’re ready, rename the project (see above) and optionally init a new git repo or add your own remote.

## Overview

This template provides a minimal but complete REST server with:

- **Feature-driven layout**: Each feature lives under `internal/<feature>/` with its own config, handler, routes, and response types.
- **Health route**: Example feature at `GET /health` (or configurable base path) returning `{"status":"ok"}`.
- **Swagger UI**: OpenAPI docs at `/swagger/index.html`.
- **Configuration**: All settings from environment variables; godotenv for `.env` loading.
- **Logger interface**: Central logging interface; no `fmt.Println` or `log.Println`.

## Architecture

- **`cmd/server/`**: Application entrypoint. Loads env, builds config, wires features, starts Gin.
- **`internal/config/`**: Root configuration and logger interface. All feature configs are part of the root config.
- **`internal/<feature>/`**: One directory per feature. Each has:
  - **`config/`**: Feature-specific config struct (filled from root config).
  - **handler.go**: HTTP handlers.
  - **routes.go**: Route registration for the feature.
  - **response.go**: Response structs/helpers (and optionally **model.go** if the feature has domain models).

Features do not depend on each other; they are registered in `main` and receive config and logger via dependency injection.

## Feature-driven development

Every feature follows the same pattern:

1. **Config**: `internal/<feature>/config/config.go` with a struct. Loaded at startup and attached to the root config.
2. **Handler**: Business logic and HTTP handling. Receives logger (and optionally config) via constructor.
3. **Routes**: A `Register(router, handler, config)` (or similar) that mounts the feature’s routes on a Gin group.
4. **Response**: Shared response types and helpers so all responses are structured JSON.

The template includes the **health route feature** as the example. Add new features by adding a new directory under `internal/` and registering it in `cmd/server/main.go`.

## Configuration

- Configuration is **env-only**. No hardcoded ports or hosts.
- **godotenv** loads a `.env` file at startup.
- Copy `.example.env` to `.env` and set values. Variables are grouped by module (server, health, reserved for future features).

Required for running the server:

- `SERVER_PORT`: Port the server listens on (e.g. `8080`).


## Run locally

1. Copy `.example.env` to `.env` and set `SERVER_PORT`.
2. Install dependencies: `go mod tidy`.
3. Generate Swagger docs (see below), then run:

   ```bash
   go run ./cmd/server
   ```

4. Open `http://localhost:8080/health` and `http://localhost:8080/swagger/index.html`.

## Run with Docker

The Dockerfile is a placeholder. **Dockerfile implementation will be provided in a later prompt.** When available, use:

```bash
docker build -t go-service-template .
docker run -p 8080:8080 --env-file .env go-service-template
```

## Generate Swagger docs

From the project root:

```bash
swag init -g cmd/server/main.go --parseDependency --parseInternal -o docs
```

Install the swag CLI if needed:

```bash
go install github.com/swaggo/swag/cmd/swag@latest
```

Then rebuild and run the server. Swagger UI is at **`/swagger/index.html`**.

## Extending with new features

1. Create `internal/<feature>/` with:
   - `config/config.go` – config struct (loaded from env via root config).
   - `handler.go` – handler struct and methods.
   - `routes.go` – `Register(router, handler, config)`.
   - `response.go` – response types/helpers (and `model.go` if needed).
2. In `internal/config/config.go`: add a field for the new feature’s config (e.g. `MyFeature *myfeatureconfig.Config`) and load it in `Load()`.
3. In `cmd/server/main.go`: load the feature config, create the handler (inject logger), call `feature.Register(router, handler, cfg.MyFeature)`.
4. Add any new env vars to `.example.env` and to the feature’s config loader.

Keep features isolated and use the logger interface only (no stdlib logging).

## AI compatibility

The template is structured for AI-assisted development:

- **Clear separation of concerns**: Config, handlers, routes, and responses are in small, focused files.
- **Explicit types**: Structs and interfaces are named and used consistently.
- **No global mutable state**: Config and logger are passed in; no hidden globals.
- **Dependency injection**: Handlers receive logger and config via constructors or arguments.
- **Single source of rules**: All architectural and coding rules are in **`.cursor/agents/go-agent.mda`**. When changing structure or dependencies, update that file.

## Debugger agent

The project is set up to work with a **debugger agent** from the Awesome Cursor Agent GitHub project (or similar Cursor/agent ecosystems).

- **Installation**: Follow the agent’s repo (e.g. clone, install CLI or Cursor extension). See [Awesome Cursor resources](https://github.com/hao-ji-xing/awesome-cursor) for links and options.
- **Integration**: The agent can use Cursor rules (e.g. `.cursor/rules` or project instructions), MCP servers, or project-specific config to attach to this codebase. Ensure `.cursor/agents/go-agent.mda` and this README are available as context.
- **Use in AI-driven development**: The agent helps with runtime visibility (e.g. logs, traces), debugging, and hypothesis testing while you extend features or fix bugs, without adding `fmt`/`log` in the code.

Replace the placeholder repo/URL with your chosen debugger agent when you have one.
