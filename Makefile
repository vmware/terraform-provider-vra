VERSION=0.0.1
NAME=terraform-provider-cas

SRC = $(shell find . -type f -name '*.go' -not -path "./vendor/*")

.PHONY: all build check clean dev fmt simplify race release

all: check build

build:
	for os in darwin linux windows; do \
	  GOARCH=amd64 GOOS=$$os go build -o ${NAME}-$$os; \
	done

check:
	@gofmt -d ${SRC}
	@test -z "$(shell gofmt -l ${SRC} | tee /dev/stderr)" || { echo "Fix formatting issues with 'make fmt'"; exit 1; }
	@ret=0 && for d in $$(go list ./... | grep -v /vendor/); do \
		test -z "$$(golint $${d} | tee /dev/stderr)" || ret=1; \
		done; exit $$ret
	@go tool vet main.go
	@go tool vet cas
	@go tool vet sdk
	go test ./...

clean:
	rm -rf pkg
	rm terraform-provider-cas*

dev:
	GOARCH=$$(go env GOARCH) GOOS=$$(go env GOOS) go install

fmt:
	@gofmt -l -w $(SRC)

simplify:
	@gofmt -s -l -w $(SRC)

race:
	go test -race ./...

release:
	for os in darwin linux windows; do \
	  mkdir -p pkg/$$os; \
	  GOARCH=amd64 GOOS=$$os go build -o pkg/$$os/${NAME}; \
	  (cd pkg/$$os; zip ../../${NAME}_${VERSION}_$${os}_amd64.zip ${NAME}); \
	done
