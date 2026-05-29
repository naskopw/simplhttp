.PHONY: build-ui build-go build clean build-production

build-ui:
	cd ui && npm install && npm run build

build-go:
	go build -o simphttp main.go

# Optimized build for production:
# - CGO_ENABLED=0: Static binary for portability
# - -ldflags="-s -w": Strips debug info/symbols for smaller size
# - -trimpath: Removes file system paths for security/reproducibility
build-production: build-ui
	CGO_ENABLED=0 go build -ldflags="-s -w" -trimpath -o simphttp main.go

build: build-ui build-go

run: build
	./simphttp

clean:
	rm -rf ui/dist
	rm -f simphttp
