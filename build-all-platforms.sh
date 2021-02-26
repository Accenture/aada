#!/bin/bash

GOOS=linux GOARCH=amd64 go build -o binaries/linux/aada
GOOS=darwin GOARCH=amd64 go build -o binaries/mac/aada
GOOS=windows GOARCH=amd64 go build -o binaries/windows/aada.exe

gon aada.hcl

