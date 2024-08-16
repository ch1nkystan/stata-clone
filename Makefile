run-linters:
	golangci-lint run

run-server:
	go run cmd/server/*.go

run-migrator:
	go run cmd/migrator/*.go

run-worker:
	go run cmd/worker/*.go

run-ticker:
	go run cmd/ticker/*.go
