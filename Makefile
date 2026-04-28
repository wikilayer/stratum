.PHONY: icons-sync design-system tidy

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
