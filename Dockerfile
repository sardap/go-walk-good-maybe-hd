FROM golang:1.16 as builder

RUN apt-get update -y && apt-get install -y --no-install-recommends xorg-dev libgl1-mesa-dev

RUN mkdir /app
WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

FROM builder as linux

RUN go test ./... && go build .

FROM builder as windows

ENV GOARCH amd64
ENV GOOS windows

RUN go build .

FROM alpine:latest

RUN mkdir /input

COPY --from=linux /app/walk-good-maybe-hd /input/walk-good-maybe-hd-linux
COPY --from=windows /app/walk-good-maybe-hd.exe /input/walk-good-maybe-hd-windows.exe

VOLUME [ "/output" ]
WORKDIR /output

CMD [ "cp", "-a", "/input/.", "/output/" ]

