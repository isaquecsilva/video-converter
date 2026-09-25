.PHONY: build

run:
	@go run .

build:
	@go build -trimpath -ldflags='-s -w' -o ./bin/video_converter.exe

tests:
	@go test ./... -v
