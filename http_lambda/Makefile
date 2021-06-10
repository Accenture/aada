all: http_lambda.zip

http_lambda.zip: http_lambda
	zip http_lambda.zip http_lambda

http_lambda: *.go go.mod go.sum
	GOOS=linux GOARCH=amd64 go build -o http_lambda .
