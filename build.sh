#!/bin/bash

# macOS 64-bit
GOOS=darwin GOARCH=amd64 go build -o bin/mac/smcfix smcfix.go
# macOS 32-bit
# GOOS=darwin GOARCH=386 go build -o bin/mac/smcfix32 smcfix.go

# Windows 64-bit
GOOS=windows GOARCH=amd64 go build -o bin/win/smcfix.exe smcfix.go
# Windows 32-bit
# GOOS=windows GOARCH=386 go build -o bin/win/smcfix32.exe smcfix.go

# Linux 64-bit
GOOS=linux GOARCH=amd64 go build -o bin/linux/smcfix smcfix.go
# Linux 32-bit
# GOOS=linux GOARCH=386 go build -o bin/linux/smcfix32 smcfix.go