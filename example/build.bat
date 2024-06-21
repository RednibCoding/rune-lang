@echo off

:: Check if the win directory exists, if not, create it
if not exist win (
    mkdir win
)

:: Check if the linux directory exists, if not, create it
if not exist linux (
    mkdir linux
)

:: Check if the macos directory exists, if not, create it
if not exist macos (
    mkdir macos
)

:: Cross-compile for Windows
set GOOS=windows
set GOARCH=amd64
go build -ldflags="-s -w" -o win\rune.exe

:: Cross-compile for Linux
set GOOS=linux
set GOARCH=amd64
go build -ldflags="-s -w" -o linux\rune

:: Cross-compile for macOS
set GOOS=darwin
set GOARCH=amd64
go build -ldflags="-s -w" -o macos\rune

echo Cross-compilation completed.
pause
