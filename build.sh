#!/bin/bash

apt-get update -y && apt-get install -y --no-install-recommends \
	zip \
	xorg-dev \
	libx11-dev \
	libgl1-mesa-dev \
	libasound2-dev \
	libgles2-mesa-dev \
	libalut-dev \
	libxcursor-dev \
	libxi-dev \
	libxinerama-dev \
	libxrandr-dev \
	libxxf86vm-dev \
	libglfw3-dev \
	xvfb \
	xauth

cd /app

go mod download

echo "running tests"
xvfb-run go test -race -coverprofile=coverage.out -covermode=atomic ./...

if [ $? -eq 0 ]
then
  echo "tests passed"
else
  echo "tests failied" >&2
  exit 1
fi

export GOARCH=amd64
export GOOS=linux

echo "building linux version"
go build -ldflags "-w" -o bin/walk-good-maybe-hd-linux-amd64 .

export GOARCH=amd64
export GOOS=windows

echo "building windows version"
go build -ldflags "-w" -o bin/walk-good-maybe-hd-windows-amd64.exe .

export GOOS=js
export GOARCH=wasm

echo "building wasm version"
go build -o bin/walk-good-maybe-hd.wasm .

cd bin

zip walk-good-maybe-hd-linux-amd64 walk-good-maybe-hd-linux-amd64
zip walk-good-maybe-hd-windows-amd64.zip walk-good-maybe-hd-windows-amd64.exe
zip walk-good-maybe-hd-wasm.zip walk-good-maybe-hd.wasm

ls

echo "done zipping"
