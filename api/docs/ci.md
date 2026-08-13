```mermaid
flowchart LR

Push["git push"]

Rust["Rust workflow\n(generator/)"]
Go["Go workflow\n(api/)"]

RCheck["cargo check --all-targets"]
RFormat["cargo fmt --check"]
RLint["cargo clippy -D warnings"]
RTest["cargo test --all-targets"]

GCheck["go build + go vet"]
GFormat["gofmt -l"]
GLint["golangci-lint"]
GTest["go test ./..."]

Push --> Rust
Push --> Go

Rust --> RCheck
Rust --> RFormat
Rust --> RLint
Rust --> RTest

Go --> GCheck
Go --> GFormat
Go --> GLint
Go --> GTest
```

Two independent workflows, one per codebase, both triggered on every push. Within each, the four jobs run in parallel and any single failure fails the workflow. A new push to the same branch cancels whatever run is still in progress for it (`concurrency: cancel-in-progress`), so only the latest commit's result matters.

See [.github/workflows/rust.yml](../../.github/workflows/rust.yml) and [.github/workflows/go.yml](../../.github/workflows/go.yml) for the actual job definitions.
