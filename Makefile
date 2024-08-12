# golang settings
GO ?= go
GOLINT ?= golangci-lint
GOOS := $(shell $(GO) env GOOS)
GOARCH := $(shell $(GO) env GOARCH)
BIN := $(abspath ./bin/$(GOOS)_$(GOARCH))
GO_ENV ?= GOPRIVATE=github.com/tusmasoma GOBIN=$(BIN)

# tools
$(shell mkdir -p $(BIN))

GOLANGCI_LINT_VERSION := v1.55.2
$(BIN)/golangci-lint-$(GOLANGCI_LINT_VERSION):
	unlink $(BIN)/golangci-lint || true
	$(GO_ENV) ${GO} install github.com/golangci/golangci-lint/cmd/golangci-lint@$(GOLANGCI_LINT_VERSION)
	mv $(BIN)/golangci-lint $(BIN)/golangci-lint-$(GOLANGCI_LINT_VERSION)
	ln -s $(BIN)/golangci-lint-$(GOLANGCI_LINT_VERSION) $(BIN)/golangci-lint

MOCKGEN_VERSION := 1.6.0
$(BIN)/mockgen-$(MOCKGEN_VERSION):
	unlink $(BIN)/mockgen || true
	$(GO_ENV) ${GO} install github.com/golang/mock/mockgen@v$(MOCKGEN_VERSION)
	mv $(BIN)/mockgen $(BIN)/mockgen-$(MOCKGEN_VERSION)
	ln -s $(BIN)/mockgen-$(MOCKGEN_VERSION) $(BIN)/mockgen

GOIMPORTS_VERSION := v0.19.0
$(BIN)/goimports-$(GOIMPORTS_VERSION):
	unlink $(BIN)/goimports || true
	$(GO_ENV) ${GO} install golang.org/x/tools/cmd/goimports@$(GOIMPORTS_VERSION)
	mv $(BIN)/goimports $(BIN)/goimports-$(GOIMPORTS_VERSION)
	ln -s $(BIN)/goimports-$(GOIMPORTS_VERSION) $(BIN)/goimports

GOFUMPT_VERSION := v0.6.0
$(BIN)/gofumpt-$(GOFUMPT_VERSION):
	unlink $(BIN)/gofumpt || true
	$(GO_ENV) ${GO} install mvdan.cc/gofumpt@$(GOFUMPT_VERSION)
	mv $(BIN)/gofumpt $(BIN)/gofumpt-$(GOFUMPT_VERSION)
	ln -s $(BIN)/gofumpt-$(GOFUMPT_VERSION) $(BIN)/gofumpt

# テストターゲット: テストを実行
.PHONY: test
test:
	$(GO) test -v -count=1 ./...

# golangci-lint: lint for all under the PKG
.PHONY: lint
lint: PKG ?= ./...
lint: $(BIN)/golangci-lint-$(GOLANGCI_LINT_VERSION)
	$(BIN)/golangci-lint run -c ./.golangci.yml $(PKG)

# golangci-lint: lint for all go files have diff
.PHONY: lint-diff
lint-diff: PKG ?= ./...
lint-diff: $(BIN)/golangci-lint-$(GOLANGCI_LINT_VERSION)
	$(BIN)/golangci-lint run -c ./.golangci.yml $(PKG) | reviewdog -f=golangci-lint -diff="git diff origin/main"

.PHONY: fmt
fmt: $(BIN)/goimports-$(GOIMPORTS_VERSION) $(BIN)/gofumpt-$(GOFUMPT_VERSION)
	FILES=$$(find . -type f -name "*.go") && \
	${GO_ENV} $(BIN)/goimports -local "github.com/tusmasoma/go-clean-arc" -w $${FILES} && \
	${GO_ENV} $(BIN)/gofumpt -l -w $${FILES}

.PHONY: generate
generate: generate-deps
	@for dir in $$(find . -type d | sed '1,1d' | sed 's@./@@') ; do \
		if [ -n "$$(git diff --name-only origin/main "$${dir}")" ]; then \
			echo "go generate ./$${dir}/..." && \
			(cd "$${dir}" && PATH="$(BIN):$(PATH)" ${GO_ENV} ${GO} generate ./...) || exit 1; \
		fi; \
	done
	$(MAKE) fmt

.PHONY: generate-all
generate-all: generate-deps
	@PATH="$(BIN):$(PATH)" ${GO_ENV} ${GO} generate ./...
	$(MAKE) fmt

.PHONY: generate-deps
generate-deps: $(BIN)/mockgen-$(MOCKGEN_VERSION)

.PHONY: tidy
tidy:
	$(GO) mod tidy

.PHONY: build
build:
	$(GO) build -v ./...

.PHONY: bin-clean
bin-clean:
	$(RM) -r ./bin