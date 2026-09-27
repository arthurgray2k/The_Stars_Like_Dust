.PHONY: all build test test-coverage vet fmt clean run-biron run-artemisia run-aratap

BINARY=stars

all: fmt vet test build

build:
	go build -o $(BINARY) ./cmd/stars

test:
	go test -v ./...

test-coverage:
	go test -v -cover ./...

vet:
	go vet ./...

fmt:
	go fmt ./...

clean:
	rm -f $(BINARY)
	rm -rf data/*.db*

run-biron: build
	./$(BINARY) --pov biron_farrill

run-artemisia: build
	./$(BINARY) --pov artemisia_hinriad

run-aratap: build
	./$(BINARY) --pov simok_aratap
