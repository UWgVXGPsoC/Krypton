go clean
$Env:GOOS = 'windows'; $Env:GOARCH = '386'; go build -ldflags="-s -w"