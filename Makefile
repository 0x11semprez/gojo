.PHONY: build build-generator build-api \
        test test-generator test-api \
        run fmt lint clean

# Compiled rust binary the Go API shells out to (see
# api/internal/cryptography/keygen). Override with GENERATOR_BIN_PATH
# to point the API at a different build (e.g. in a container).
GENERATOR_BIN := generator/target/release/generator
API_BIN := api/bin/gojo

## build: build the rust generator (release) and the go api binary.
build: build-generator build-api

## build-generator: build the rust key generator in release mode.
build-generator:
	cd generator && cargo build --release

## build-api: build the go api binary.
build-api:
	cd api && go build -o bin/gojo ./cmd

## test: run both test suites (rust generator, go api).
test: test-generator test-api

## test-generator: run the rust generator's test suite.
test-generator:
	cd generator && cargo test --all-targets

## test-api: run the go api's test suite (unit + integration; needs a
## running postgres reachable via api/.env.test, and the generator
## binary built -- see build-generator).
test-api:
	cd api && go test ./...

## run: build everything, then start the api (applies pending
## migrations on startup, then listens for requests).
run: build
	cd api && ./bin/gojo

## fmt: format both codebases.
fmt:
	cd generator && cargo fmt --all
	cd api && gofmt -w .

## lint: lint both codebases (matches the checks CI runs).
lint:
	cd generator && cargo clippy --all-targets -- -D warnings
	cd api && go vet ./...
	cd api && golangci-lint run ./...

## clean: remove build artifacts.
clean:
	rm -rf generator/target
	rm -rf $(API_BIN)
