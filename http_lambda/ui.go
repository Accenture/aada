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
        @keyframes countdown {
            from { width: 100px; }
            to   { width: 0; }
        }
        #countdown {
            border-radius: 2px;
            margin-left: auto;
            margin-right: auto;
            width: 100px;
            height: 5px;
            background-color: black;
            animation: countdown 5s linear forwards;
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
    <div id="countdown"></div>
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

const consolePage = `<!DOCTYPE html>
<html>
	<head>
		<title>AADA AWS Console Access</title>
		<style>
			body {
				font-family: sans-serif;
				text-align: center;
			}
			#accountsList {
				margin-left: 2em;
				margin-right: 2em;
				display: grid;
				grid-template-columns: repeat(auto-fill, minmax(200px, 1fr));
				grid-gap: 2em;
			}
			#accountsList a {
				padding-top: 0.75em;
				padding-bottom: 0.75em;
				background: #fb992720;
				border: 1px solid darkorange;
			}
			.account {
				margin-top: 0.25em;
				margin-bottom: 0.50em;
				font-size: xx-small;
			}
		</style>
	</head>
	<body>
		<p><strong>AADA</strong></p>
		<svg height='100px' width='100px'  fill="#0296CC" xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink" version="1.1" x="0px" y="0px" viewBox="0 0 54 56" enable-background="new 0 0 54 56" xml:space="preserve">
			<g>
				<path fill="#0296CC" d="M46.91,47.5l6.48-3.24l-17.64-6.5l6.5,17.64l3.24-6.48l6.8,6.8C52.49,55.9,52.74,56,53,56   s0.51-0.1,0.71-0.29c0.39-0.39,0.39-1.02,0-1.41L46.91,47.5z M39.13,41.13l9.23,3.4l-3.89,1.94l-1.94,3.89L39.13,41.13z"></path>
				<path fill="#0296CC" d="M27,14h18c0.55,0,1-0.45,1-1s-0.45-1-1-1H27c-0.55,0-1,0.45-1,1S26.45,14,27,14z"></path>
				<path fill="#0296CC" d="M27,28h18c0.55,0,1-0.45,1-1s-0.45-1-1-1H27c-0.55,0-1,0.45-1,1S26.45,28,27,28z"></path>
				<path fill="#0296CC" d="M20.58,7.24c-0.42-0.36-1.05-0.31-1.41,0.11l-6.25,7.29L9.1,11.77c-0.44-0.33-1.07-0.24-1.4,0.2   c-0.33,0.44-0.24,1.07,0.2,1.4l4.57,3.43c0.18,0.13,0.39,0.2,0.6,0.2c0.28,0,0.56-0.12,0.76-0.35l6.86-8   C21.05,8.23,21,7.6,20.58,7.24z"></path>
				<path fill="#0296CC" d="M20.58,21.24c-0.42-0.36-1.05-0.31-1.41,0.11l-6.25,7.29L9.1,25.77c-0.44-0.33-1.07-0.24-1.4,0.2   c-0.33,0.44-0.24,1.07,0.2,1.4l4.57,3.43c0.18,0.13,0.39,0.2,0.6,0.2c0.28,0,0.56-0.12,0.76-0.35l6.86-8   C21.05,22.23,21,21.6,20.58,21.24z"></path>
				<path fill="#0296CC" d="M19.17,36.35l-6.25,7.29L9.1,40.77c-0.44-0.33-1.07-0.24-1.4,0.2c-0.33,0.44-0.24,1.07,0.2,1.4l4.57,3.43   c0.18,0.13,0.39,0.2,0.6,0.2c0.28,0,0.56-0.12,0.76-0.35l6.86-8c0.36-0.42,0.31-1.05-0.11-1.41   C20.16,35.88,19.53,35.93,19.17,36.35z"></path>
				<path fill="#0296CC" d="M41.82,0H12.18C5.46,0,0,5.46,0,12.18v29.65C0,48.54,5.46,54,12.18,54H32c0.55,0,1-0.45,1-1s-0.45-1-1-1   H12.18C6.56,52,2,47.44,2,41.82V12.18C2,6.56,6.56,2,12.18,2h29.65C47.44,2,52,6.56,52,12.18V37c0,0.55,0.45,1,1,1s1-0.45,1-1   V12.18C54,5.46,48.54,0,41.82,0z"></path>
			</g>
		</svg>
		<p>Available Accounts</p>
		<div id="accountsList">__ACCOUNTS_LIST__</div>
    </body>
</html>
`

type ConsoleLink struct {
	Url          string
	DisplayNames []string
	Account      string
}

func buildConsolePage(links []ConsoleLink) Response {
	html := ""
	for _, link := range links {
		html = "<a href=\"" + link.Url + "\">"
		html = html + "<div class=\"account\">" + link.Account + "</a>"
		for _, displayName := range link.DisplayNames {
			html = html + "<div class=\"name\">" + displayName + "</div>"
		}
		html = html + "</a>\n"
	}

	return Response{
		StatusCode: http.StatusOK,
		Body:       strings.Replace(consolePage, "__ACCOUNTS_LIST__", html, 1),
		Headers: map[string]string{
			"Content-Type": "text/html",
		},
	}
}
