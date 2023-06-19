all: http_lambda.zip

http_lambda.zip: bootstrap
	zip http_lambda.zip bootstrap

bootstrap: *.go go.mod go.sum ui*.html
	GOOS=linux GOARCH=arm64 go build -o bootstrap -tags lambda.norpc .
