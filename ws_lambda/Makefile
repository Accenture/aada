all: ws_lambda.zip

ws_lambda.zip: ws_lambda
	zip ws_lambda.zip ws_lambda

ws_lambda: *.go go.mod go.sum
	GOOS=linux GOARCH=amd64 go build -o ws_lambda .
