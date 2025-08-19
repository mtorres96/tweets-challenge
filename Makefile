.PHONY: run test swagger deps tools swagger-run

GOBIN := $(shell go env GOPATH)/bin
SWAG  := $(GOBIN)/swag

run:
	go run ./cmd/api

test:
	go test ./...

deps:
	go mod tidy

# Instala herramientas locales si faltan
tools:
	@test -x "$(SWAG)" || (echo "Installing swag CLI..." && go install github.com/swaggo/swag/cmd/swag@v1.16.3)

# Genera la doc de Swagger en ./docs (instala swag si falta)
swagger: tools
	"$(SWAG)" init -g cmd/api/main.go -o docs

# Alternativa sin instalar en PATH (one-off)
swagger-run:
	go run github.com/swaggo/swag/cmd/swag@v1.16.3 init -g cmd/api/main.go -o docs