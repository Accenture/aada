#!/bin/bash

go version

echo building linux
GOOS=linux GOARCH=amd64 go build -o binaries/linux/aada

echo building mac
GOOS=darwin GOARCH=amd64 go build -o binaries/mac/aada

echo building windows
GOOS=windows GOARCH=amd64 go build -o binaries/windows/aada.exe

echo building windows zip
zip -j binaries/windows/aada.zip binaries/windows/aada.exe

echo signing mac binary and building mac zip
gon aada.hcl

