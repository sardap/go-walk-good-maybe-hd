#!/bin/bash

apt-get update -y 

apt-get install -y --no-install-recommends xorg-dev libgl1-mesa-dev zip

cd /app

go mod download

echo "running tests"
go test ./...

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
go build -o bin/walk-good-maybe-hd-amd64-linux .

export GOARCH=amd64
export GOOS=windows

echo "building windows version"
go build -o bin/walk-good-maybe-hd-amd64-windows.exe .

export GOOS=js
export GOARCH=wasm

echo "building wasm version"
go build -o bin/walk-good-maybe-hd-wasm .

cd bin

zip walk-good-maybe-hd-amd64-linux walk-good-maybe-hd-amd64-linux
zip walk-good-maybe-hd-amd64-windows.zip walk-good-maybe-hd-amd64-windows.exe
zip walk-good-maybe-hd-wasm.zip walk-good-maybe-hd-wasm

echo "done zipping"
