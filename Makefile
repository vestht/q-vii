default: build push

build:
	@CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o bootstrap -ldflags "-s -w" handler.go
	@zip handler.zip bootstrap
	@rm -rf bootstrap

push:
	@aws lambda update-function-code --function-name sg-rules-wd --zip-file fileb://handler.zip
	@rm -rf handler.zip