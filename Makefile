.PHONY: build clean deploy

build:
	env GOOS=linux go build -ldflags="-s -w" -o bin/websocket/connect websocket/connect/main.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/websocket/message websocket/message/main.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/websocket/disconnect websocket/disconnect/main.go

clean:
	rm -rf ./bin

deploy: clean build
	sls deploy --verbose
