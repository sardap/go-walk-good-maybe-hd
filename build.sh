#!/bin/bash

apt-get update -y 

apt-get install -y --no-install-recommends xorg-dev libgl1-mesa-dev

cd /app

go mod download

echo "running tests"
go test ./... 

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
go build .
