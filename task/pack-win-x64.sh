#!/bin/bash -u
# This script compiles project and prepares package for Windows amd64.

fyne package -os windows --name balance_win_x64.exe
mv balance_win_x64.exe "$GOPATH/bin/"
