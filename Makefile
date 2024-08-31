.PHONY: build
build:
	@echo "==> building system..."
	@GOOS=linux go build .
	@echo "==> done!"

.PHONY: test
test:
	@echo "==> testing system..."
	@go test ./...
	@echo "==> done!"
