.PHONY: format lint lint-tools test-build test build icons-sync design-system tidy

STATICCHECK_VERSION ?= v0.8.0

format:
	gofmt -w .

lint:
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

build:
	go build ./...

# Regenerate static/icons.svg from static/icons.txt by pulling each
# name from unpkg.com/lucide-static (or a custom URL for brand icons
# Lucide doesn't ship). Re-runnable safely.
icons-sync:
	go run ./cmd/icons

# Open the standalone design-system reference in the default browser.
# It uses file:// — no server needed.
design-system:
	open design-system/index.html

tidy:
	go mod tidy
