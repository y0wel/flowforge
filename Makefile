build:
	go build -o flowforge ./cmd/flowforge

run:
	go run ./cmd/flowforge run $(workflow)

test:
	go test ./...

vet:
	go vet ./...
