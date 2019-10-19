.DEFAULT_GOAL := build

OS_LIST=windows linux darwin
ARCH_LIST=amd64 386

.PHONY: test
test: ## Run all test
	@gotest ./... -v -cover -race

.PHONY: build
build: ## Build binary file
	@go build -o aptoilet *.go

.PHONY: cross-compile
cross-compile: ## Build binaries for Windows, Linux, and macOS of x64 and x86
	@mkdir bin
	@for GOOS in ${OS_LIST}; do \
		for GOARCH in ${ARCH_LIST}; do \
			if [ $$GOOS = windows ]; then \
				GOOS=$$GOOS GOARCH=$$GOARCH go build -o bin/aptoilet-$$GOOS-$$GOARCH.exe main.go; \
			else \
				GOOS=$$GOOS GOARCH=$$GOARCH go build -o bin/aptoilet-$$GOOS-$$GOARCH main.go; \
			fi \
		done \
	done

.PHONY: clean
clean:
	@-rm -rf bin/

.PHONY: help
help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
