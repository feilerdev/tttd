#
# BUILD
#
.PHONY: build
build:
	@echo "==> building system..."
	@GOOS=linux go build .
	@echo "==> done!"

#
# TESTING
#

.PHONY: test
test:
	@echo "==> testing system..."
	# TODO(alexandreliberato): use parallel tests > td-quality => $$ 
	@go test ./... -v
	@echo "==> done!"

#
# INSTALL LINTERS
#

.PHONY: install-linters
install-linters: ## Installs linters
	@echo "==> Installing staticcheck"
	@go install honnef.co/go/tools/cmd/staticcheck@latest
	@echo "==> Installing govulncheck"
	@go install golang.org/x/vuln/cmd/govulncheck@latest
	@echo "==> Installing GCI"
	@go install github.com/daixiang0/gci@latest
	@echo "Installing golangci-lint"
	@curl -sfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $$(go env GOPATH)/bin v1.55.2

.PHONY: lint
lint: ## Run lint
	@echo "==> Running go vet"
	@go vet ./...
	@echo "==> Running staticcheck"
	@staticcheck ./...
	@echo "==> Running govulncheck"
	@govulncheck ./...