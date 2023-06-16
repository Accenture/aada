all: ws_lambda.zip

ws_lambda.zip: bootstrap
	zip ws_lambda.zip bootstrap

bootstrap: *.go go.mod go.sum
	GOOS=linux GOARCH=arm64 go build -o bootstrap -tags lambda.norpc .
