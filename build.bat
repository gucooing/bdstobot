set CGO_ENABLED=0
set GOARCH=arm64
set GOOS=linux
go build -ldflags "-w -s" -o linux-arm64 main.go
set GOOS=windows
go build -ldflags="-s -w" -o windows-arm64.exe  main.go

set GOARCH=amd64
set GOOS=linux
go build -ldflags "-w -s" -o linux-amd64  main.go
set GOOS=windows
go build -ldflags="-s -w" -o windows-amd64.exe  main.go
