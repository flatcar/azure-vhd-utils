rwildcard=$(foreach d,$(wildcard $1*),$(call rwildcard,$d/,$2) $(filter $(subst *,%,$2),$d))

.PHONY: all
all: test azure-vhd-utils lint

.PHONY: lint
lint: azure-vhd-utils
	go fmt ./...
	go vet ./...
	go mod tidy

azure-vhd-utils: $(call rwildcard, ., *.go) go.mod go.sum Makefile
	go build

.PHONY: test
test:
	go test ./...

.PHONY: clean
clean:
	rm -f azure-vhd-utils
