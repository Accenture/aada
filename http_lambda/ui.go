package main

import (
	_ "embed"
	"net/http"
	"strings"
)

//go:embed favicon.b64
var favicon string

const successPage = `<!DOCTYPE html>
<html>
	<head>
		<title>AADA Auth Status</title>
		<style>
			body {
				font-family: sans-serif;
				text-align: center;
				background-color: #00ff0020;
			}
		</style>
    </head>
    <body>
		<p><strong>AADA</strong></p>
		<svg height='100px' width='100px' fill="#000000" xmlns="http://www.w3.org/2000/svg" data-name="Layer 1" viewBox="0 0 100 100" x="0px" y="0px">
			<title>OK</title>
			<path d="M51.25,81.75A31.75,31.75,0,1,0,19.5,50,31.79,31.79,0,0,0,51.25,81.75Zm0-56.5A24.75,24.75,0,1,1,26.5,50,24.78,24.78,0,0,1,51.25,25.25Z"></path>
			<polygon points="66.75 43.35 61.8 38.4 48.13 52.07 42.9 46.84 37.95 51.79 48.13 61.97 66.75 43.35"></polygon>
		</svg>
		<p>login successful</p>
		<script>
			setTimeout(function() { window.close() }, 5000);
		</script>
    </body>
</html>
`

func buildSuccessResponse() Response {
	return Response{
		StatusCode: http.StatusOK,
		Body:       successPage,
		Headers: map[string]string{
			"Content-Type": "text/html",
		},
	}
}

const failurePage = `<!DOCTYPE html>
<html>
	<head>
		<title>AADA Auth Status</title>
		<style>
			body {
				font-family: sans-serif;
				text-align: center;
				background-color: #ff000020;
			}
		</style>
    </head>
    <body>
		<p><strong>AADA</strong></p>
		<svg height='100px' width='100px' fill="#000000" xmlns="http://www.w3.org/2000/svg" data-name="Layer 1" viewBox="0 0 100 100" x="0px" y="0px">
			<title>OK</title>
			<path d="M71.58,27.55h0a31.75,31.75,0,1,0,0,44.9A31.54,31.54,0,0,0,71.58,27.55Zm-4.95,40a24.75,24.75,0,1,1,0-35h0A24.78,24.78,0,0,1,66.63,67.5Z"></path>
			<polygon points="56.14 38.04 49.13 45.05 42.11 38.04 37.16 42.99 44.18 50 37.16 57.01 42.11 61.96 49.13 54.95 56.14 61.96 61.09 57.01 54.08 50 61.09 42.99 56.14 38.04"></polygon>
		</svg>
		<p>__MESSAGE__</p>
    </body>
</html>
`

func buildFailureResponse(message string) Response {
	return Response{
		StatusCode: http.StatusOK,
		Body:       strings.Replace(failurePage, "__MESSAGE__", message, 1),
		Headers: map[string]string{
			"Content-Type": "text/html",
		},
	}
}
