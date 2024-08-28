Gin based HTTP service.

To generate swagger docs, run 

export PATH=$PATH:$HOME/go/bin
swag init -g ./cmd/api/main.go -o docs

in the root dir.