# E2Engine end-to-end tests

End-to-end test suite for E2Engine.

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

## E2Engine

This repository is part of E2Engine.

- core — core domain model, execution logic, and public APIs
- repository — persistence implementations
- runner-local — local test execution
- cli — command-line interface
- tests — end-to-end tests for E2Engine
- demo — executable demonstration system and E2Engine examples

## License

Licensed under the Apache License, Version 2.0.