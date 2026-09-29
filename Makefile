.DEFAULT_GOAL := build

.PHONY: install-tools format lint comments test-build test build icons-sync design-system design-system-build tidy

STATICCHECK_VERSION ?= v0.8.1

install-tools: lint-tools
	python3 -m pip install --quiet --upgrade git+https://github.com/botforge-pro/commentcensor.git

format:
	gofmt -w .

comments:
	commentcensor *.go

lint: comments
	go vet ./...
	gofmt -l . | (! grep .)
	staticcheck ./...

lint-tools:
	go install honnef.co/go/tools/cmd/staticcheck@$(STATICCHECK_VERSION)

# Compile package and test code without running anything.
test-build:
	go build ./...
	go test -run '^$$' ./...

test:
	go test ./...

build: lint test-build test
	go build ./...

# Regenerate static/icons.svg from static/icons.txt by pulling each
# name from unpkg.com/lucide-static (or a custom URL for brand icons
# Lucide doesn't ship). Re-runnable safely.
icons-sync:
	go run ./cmd/icons

# Serve the design-system reference and open it in the default browser.
design-system:
	go run ./cmd/design-system -open

# Render the same templates and assets for GitHub Pages.
design-system-build:
	go run ./cmd/design-system -build

tidy:
	go mod tidy
