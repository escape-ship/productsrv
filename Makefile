all: build

init:
	@echo "Initializing..."
	@$(MAKE) tool_download
	@$(MAKE) build

build:
	@echo "Building..."
	@go mod tidy
	@go mod download
	@$(MAKE) sqlc_gen
	@${MAKE} build_alone

pushall:
	@docker build -t ghcr.io/escape-ship/productsrv:latest .
	@docker push ghcr.io/escape-ship/productsrv:latest

build_alone:
	@go build -tags migrate -o bin/$(shell basename $(PWD)) ./cmd

sqlc_gen:
	@echo "Generating sqlc..."
	@cd internal/infra/sqlc && \
	sqlc generate

tool_download:
	$(MAKE) sqlc_download

sqlc_download:
	@echo "Downloading sqlc..."
	@go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest

run:
	@echo "Running..."
	@./bin/$(shell basename $(PWD))

linter-golangci: ### check by golangci linter
	golangci-lint run
.PHONY: linter-golangci

clean:
	rm -f bin/$(shell basename $(PWD))