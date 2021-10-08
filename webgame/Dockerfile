# Only build with docker compose

FROM golang:1.16 as builder

RUN mkdir /app
WORKDIR /app

RUN cp $(go env GOROOT)/misc/wasm/wasm_exec.js .

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

RUN go run ./cmd/gen ./pkg http://localhost:8080/assets/

ENV GOOS js
ENV GOARCH wasm

RUN go build -o wgmhd.wasm ./cmd/wgmhd

FROM nginx:alpine

COPY --from=builder /app/wasm_exec.js /usr/share/nginx/html/wasm_exec.js
COPY --from=builder /app/wgmhd.wasm /usr/share/nginx/html/wgmhd.wasm
COPY webgame/index.html /usr/share/nginx/html/index.html
COPY webgame/main.html /usr/share/nginx/html/main.html

