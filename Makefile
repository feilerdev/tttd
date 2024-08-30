.PHONY: build
build:
	@echo "==> building system..."
	@GOOS=linux go build ./...
	@echo "==> done!"