# E2Engine end-to-end tests

End-to-end test suite for E2Engine.

Part of [E2Engine](https://e2engine.dev), an open-source platform for declarative end-to-end testing.

This repository verifies E2Engine through its public interfaces using the CLI, Docker, real and mocked HTTP/gRPC services, direct and socket transports, persistence, and asynchronous execution flows.

The CLI, Docker images, and helper service binaries required by the tests are built locally by the Makefile before the corresponding test suites run.

## Coverage

| Area | Coverage |
| --- | --- |
| Environment, test, and test-suite lifecycle | Create, list, get, execute, and delete flows |
| Spec validation | Environment, test, and test-suite validation, including nested HTTP/gRPC configuration |
| Spec formats | YAML and JSON execution specs |
| Direct transport | HTTP and gRPC execution |
| Socket transport | Asynchronous HTTP and gRPC execution |
| Test suites | Multi-test suite execution and persisted child executions |
| Mocked HTTP services | Fixture matching by path and body, fixture misses, and multiple-fixture no-match behavior |
| Mocked gRPC services | Fixture matching by method and message, including fixture misses |
| Real-to-mocked dependencies | Real HTTP/gRPC services calling mocked dependencies with captured calls |
| Call expectations | Required calls, expected call counts, unmatched call parameters, and missing expected dependencies |
| Execution results | Passed, failed, and error terminal states with persisted execution summaries |
| Runner lifecycle | Direct execution waits for completion; socket executions run asynchronously and sequentially |
| Docker | CLI execution in containers, mounted configuration and persistent data, mocked HTTP/gRPC execution, and communication with real services over a Docker network |

The E2E suite focuses on behavior across component boundaries. Exhaustive field-level and implementation-specific cases belong in the unit tests of the corresponding E2Engine repository.

## Repository structure

```text
cli/          CLI E2E tests, harness, test data, and generated CLI binary
docker/       Docker E2E tests, harness, configuration, and test data
services/     Helper HTTP/gRPC services, generated service binaries, and reusable service Dockerfile
util/         Shared E2E comparison and test utilities
```

Generated binaries are written to cli/bin/ and services/bin/ and are not committed. Docker tests use locally built E2Engine and helper service images.

## Development

Run the CLI E2E test suite:

```bash
make test-cli
```

Run the Docker E2E test suite:

```bash
make test-docker
```

Run the linter:

```bash
make lint
```

Run the complete verification suite:

```bash
make verify
```

The test targets build the binaries and Docker images required by their respective suites before running the tests.

## Stable branch

The `stable` branch provides a known-good version of the E2E test suite for downstream E2Engine repositories.

Downstream CI, such as the CLI build, should use `stable` rather than `main`. This prevents changes under development in the test suite from unexpectedly breaking downstream builds.

Changes are promoted in the following order:

1. Update dependencies and tests on `main`.
2. Run the complete test repository CI and ensure it passes.
3. Promote `main` to `stable`:

```bash
git switch main
git pull --ff-only

# after main CI is green
git push origin main:stable
```

The `stable` branch must always point to a commit from main that has already passed CI; promotion should only fast-forward the branch and must not create new commits.

4. Update downstream repositories that depend on the new versions.

For example, when updating Core or instrumentation dependencies:

```text
dependency release
        ↓
tests/main
        ↓
tests CI passes
        ↓
tests/stable
        ↓
downstream repository update
```

## E2Engine

This repository is part of E2Engine.

- [core](https://github.com/e2engine/core) — core domain model, execution logic, and public APIs
- [repository](https://github.com/e2engine/repository) — persistence implementations
- [runner-local](https://github.com/e2engine/runner-local) — local test execution
- [cli](https://github.com/e2engine/cli) — command-line interface
- [tests](https://github.com/e2engine/tests) — end-to-end tests for E2Engine
- [demo](https://github.com/e2engine/demo) — executable demonstration system and E2Engine usage examples
- [instrumentation-go](https://github.com/e2engine/instrumentation-go) — Go instrumentation library for E2Engine

## License

Licensed under the Apache License, Version 2.0. See [LICENSE](LICENSE).