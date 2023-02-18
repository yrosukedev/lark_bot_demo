GOOS=linux GOARCH=amd64 go build -o ./out/v0.0.1/main
zip -j ./out/v0.0.1/main.zip ./out/v0.0.1/main
