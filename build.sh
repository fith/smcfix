#!/bin/bash

# Uncomment the ones you want to build.

# macOS 64-bit
# GOOS=darwin GOARCH=amd64 go build -o bin/mac/smcfix smcfix.go
# macOS 32-bit
# GOOS=darwin GOARCH=386 go build -o bin/mac/smcfix32 smcfix.go

# Windows 64-bit
# GOOS=windows GOARCH=amd64 go build -o bin/win/smcfix.exe smcfix.go
# Windows 32-bit
# GOOS=windows GOARCH=386 go build -o bin/win/smcfix32.exe smcfix.go

# Linux 64-bit
# GOOS=linux GOARCH=amd64 go build -o bin/linux/smcfix smcfix.go
# Linux 32-bit
# GOOS=linux GOARCH=386 go build -o bin/linux/smcfix32 smcfix.go


# Packaging with icon
# https://developer.fyne.io/started/packaging
# go install fyne.io/fyne/v2/cmd/fyne@latest
# macOS
# fyne package -os darwin -icon assets/icon.png
# Windows
# fyne package -os windows -icon assets/icon.png
# Linux
# fyne package -os linux -icon assets/icon.png



# macOS package in dmg
# https://github.com/create-dmg/create-dmg
# create-dmg \
#  --app-drop-link 196 48 \
#  --window-size 396 196 \
#  --icon-size 100 \
#  --icon "smcfix.app" 48 48 \
#  "smcfix.dmg" \
#  "smcfix.app"
