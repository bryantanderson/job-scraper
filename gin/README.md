To generate swagger docs, run in the root dir:

export PATH=$PATH:$HOME/go/bin
swag init -g ./cmd/api/main.go -o docs

To run tests: 

go test ./... -cover

To run tests with a coverprofile:

go test ./... -coverprofile=reports/cover.out && go tool cover -html=reports/cover.out