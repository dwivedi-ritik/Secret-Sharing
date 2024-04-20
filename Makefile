run: build
	./.bin/app 
build: clear
	@go build -o ./.bin/app ./cmd/api/main.go
clear:
	@rm -rf ./bin/*