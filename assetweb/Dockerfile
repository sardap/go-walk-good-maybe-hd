# Only build with docker compose

FROM golang:1.16-alpine as builder

RUN mkdir /app
WORKDIR /app

COPY ./assetweb/go.mod .
COPY ./assetweb/go.sum .
RUN go mod download

COPY ./assetweb .

RUN go build .

FROM alpine:3.13.6

RUN mkdir /app

COPY --from=builder /app/assetweb /app/assetweb

COPY ./out /app/out

ENV GIN_MODE=release

ENTRYPOINT [ "/app/assetweb" ]
CMD [ "/app/out" ]
