FROM golang:1.18.3-alpine3.16 as server

WORKDIR /src

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY main.go ./

RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags '-w' -o main

FROM golang:1.18.3-alpine3.16 as wasm-builder

WORKDIR /src

COPY wasm/ ./

RUN go mod tidy

RUN GOOS=js GOARCH=wasm go build -o particle.wasm main.go

FROM node:18-alpine3.16 as content

WORKDIR /src

COPY content/ ./

RUN npm install
RUN npm run build-css

FROM scratch as final

COPY --from=server /src/main /main
COPY --from=content /src/templates /usr/src/content/templates
COPY --from=content /src/static /usr/src/content/static
COPY --from=content /src/routes.json /usr/src/content/routes.json
COPY --from=wasm-builder /src/particle.wasm /usr/src/content/static/particle.wasm

ENTRYPOINT ["/main"]
