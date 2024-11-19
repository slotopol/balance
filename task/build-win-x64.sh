#!/bin/bash -u
# This script compiles project for Windows amd64.

wd=$(realpath -s "$(dirname "$0")/..")
cp -ruv "$wd/appdata/"* "$HOME/AppData/Roaming/fyne/slotopol.balance"

go env -w GOOS=windows GOARCH=amd64 CGO_ENABLED=1
go build -o "$GOPATH/bin/balance_win_x64.exe" -v -tags="" $wd
