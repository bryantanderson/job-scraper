FROM golang:1.22

WORKDIR /go/src/app

COPY . .

RUN go install -mod=mod github.com/githubnemo/CompileDaemon

EXPOSE 8090

ENTRYPOINT CompileDaemon -log-prefix=false --directory="./" --build="go build cmd/api/main.go" --command="./main"
