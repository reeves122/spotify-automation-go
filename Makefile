run:
	go run cmd/main.go

test:
	go test ./...

coverage:
	go test -json -coverprofile=cover.out ./... > result.json
	go tool cover -func cover.out
	go tool cover -html=cover.out

mocks:
	go install github.com/golang/mock/mockgen@latest
	mockgen -source adapter/adapter.go -destination mocks/mock_adapter/mock.go -package mock_adapter

fmt:
	go fmt ./...

tidy:
	go mod tidy