FROM golang:1.18.3-alpine3.16 as builder

WORKDIR /app

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY main.go ./

RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags '-w' -o main

FROM scratch as final

COPY --from=builder /app/main /main

ENTRYPOINT ["/main"]
