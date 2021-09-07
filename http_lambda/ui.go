package main

import (
	_ "embed"
	"net/http"
	"sort"
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
<html><head><meta http-equiv="Content-Type" content="text/html; charset=windows-1252">
    <title>AADA AWS Console Access</title>
    <style type="text/css">
        body {
            font-family: sans-serif;
            text-align: center;
        }
        #accountsList {
            margin-left: 2em;
            margin-right: 2em;
            display: grid;
            grid-template-columns: repeat(auto-fill, minmax(225px, 1fr));
            grid-gap: 2em;
        }
        .link {
            padding-top: 0.75em;
            padding-bottom: 0.75em;
            background: #36688d;
            border: 1px solid #5486ab;
            border-radius: 4px;
        }
        .account {
            margin-top: 0.25em;
            margin-bottom: 0.50em;
            font-size: small;
            text-decoration: underline;
        }
        a {
            color: #dddddd;
            transition-duration: 1s;
            text-decoration: none;
        }
        a:visited {
            color: #dddddd;
            transition-duration: 1s;
        }
        a:hover {
            color: white;
            transition-duration: 1s;
        }
        #titles {
            position: absolute;
            right: 2%;
            top: 2%;
        }
        .title {
            display: none;
            color: #dddddd;
            transition-duration: 1s;
            text-decoration: none;
            text-transform: uppercase;
            font-weight: bold;
        }
    </style>
</head>
<body>
<p><strong>AADA</strong></p>
<svg height="100px" width="100px" fill="#0296CC" xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink" version="1.1" x="0px" y="0px" viewBox="0 0 54 56" enable-background="new 0 0 54 56" xml:space="preserve">
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
<p>Available Accounts and Roles</p>
<div id="accountsList">__ACCOUNTS_LIST__</div>
<div id="titles">
    <label for="titleToggle">Toggle Titles</label>
    <input id="titleToggle" type="checkbox" />
</div>
<script type="application/javascript">
    document.getElementById("titleToggle").onchange = function (ev) {
        let goal = "none";
        let checkbox = ev.target;
        localStorage.setItem("showTitles", checkbox.checked);
        if (checkbox.checked) {
            goal = "block";
        }
        let titles = document.getElementsByClassName("title");
        for (i = 0; i < titles.length; i++) {
            titles[i].style.display = goal;
        }
    }
    let titles = document.getElementsByClassName("title");
    for (let i = 0; i < titles.length; i++) {
        let target = titles[i];
        target.addEventListener("input", function (ev) {
            let title = ev.target;
            localStorage.setItem("titles/" + title.dataset.account, title.innerHTML);
        });
        let title = localStorage.getItem("titles/" + target.dataset.account);
        if (title != null) {
            target.innerHTML = title;
        }
    }
    if (localStorage.getItem("showTitles") === "true") {
        document.getElementById("titleToggle").click();
    }
</script>
</body></html>
`
type ConsoleLink struct {
	Url          string
	DisplayNames []string
	Account      string
}

func buildConsolePage(links []ConsoleLink) Response {
	sort.Slice(links, func(i, j int) bool {
		left := strings.ToLower(links[i].DisplayNames[0])
		right := strings.ToLower(links[j].DisplayNames[0])
		return left < right
	})

	// <div class="link">
	//   <div class="title" data-account="155076074322" contenteditable>Title</div>
	//   <a href="https://myapps.microsoft.com/signin/113614_155076074322_ACP_AWS_APP/973938ac-9d51-4b06-a33b-db31e1abb0b8?tenantId=e0793d39-0939-496d-b129-198edd916feb">
	//     <div class="account">155076074322</div>
	//     <div class="name">Administrator</div>
	//   </a>
	// </div>

	html := ""
	for _, link := range links {
		html += "\n<div class=\"link\"><div class=\"title\" data-account=\"" + link.Account + "\" contenteditable>Title</div>"
		html += "<a href=\"" + link.Url + "\">"
		html = html + "<div class=\"account\">" + link.Account + "</div>"
		for _, displayName := range link.DisplayNames {
			html = html + "<div class=\"name\">" + displayName + "</div>"
		}
		html = html + "</a></div>"
	}

	return Response{
		StatusCode: http.StatusOK,
		Body:       strings.Replace(consolePage, "__ACCOUNTS_LIST__", html, 1),
		Headers: map[string]string{
			"Content-Type": "text/html",
		},
	}
}

const downloadsPage = `<!DOCTYPE html>
<html><head><meta http-equiv="Content-Type" content="text/html; charset=windows-1252">
    <title>AADA Downloads</title>
    <link rel="preconnect" href="https://fonts.googleapis.com">
    <link rel="preconnect" href="https://fonts.gstatic.com" crossorigin>
    <link href="https://fonts.googleapis.com/css2?family=Roboto+Mono&display=swap" rel="stylesheet">
    <style>
        body {
            font-family: sans-serif;
            text-align: center;
        }
        #binariesList {
            margin-left: 2em;
            margin-right: 2em;
            display: flex;
            flex-direction: column;
            align-items: center;
        }
        a {
            margin: 0.25em;
            font-family: 'Roboto Mono', monospace;
            color: #dddddd;
            transition-duration: 1s;
            text-decoration: none;
            padding-top: 0.75em;
            padding-bottom: 0.75em;
            background: #36688d;
            border: 1px solid #5486ab;
            border-radius: 4px;
            width: 20em;
        }
        a:visited {
            color: #dddddd;
            transition-duration: 1s;
        }
        a:hover {
            color: white;
            transition-duration: 1s;
        }
    </style>
</head>
<body>
<p><strong>AADA Downloads</strong></p>
<svg height='100px' width='100px' fill="#19772e" xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink" xmlns:sketch="http://www.bohemiancoding.com/sketch/ns" viewBox="0 0 12 12" version="1.1" x="0px" y="0px">
	<title>Noun_10x10</title>
	<desc>Created with Sketch.</desc>
	<g stroke="none" stroke-width="1" fill="none" fill-rule="evenodd" sketch:type="MSPage">
		<g sketch:type="MSArtboardGroup" fill="#19772e">
			<g sketch:type="MSLayerGroup" transform="translate(1.000000, 2.000000)">
				<path d="M5,9 C4.875,9 4.79166667,8.95909091 4.70833333,8.87727273 L3.45833333,7.65 C3.29166667,7.48636364 3.29166667,7.24090909 3.45833333,7.07727273 C3.625,6.91363636 3.875,6.91363636 4.04166667,7.07727273 L5.29166667,8.30454545 C5.45833333,8.46818182 5.45833333,8.71363636 5.29166667,8.87727273 C5.20833333,8.95909091 5.125,9 5,9 L5,9 Z" sketch:type="MSShapeGroup"></path>
				<path d="M2.38541667,3.16022727 C2.45833333,3.08863636 2.5,2.98636364 2.5,2.86363636 C2.5,1.71818182 3.41666667,0.818181818 4.58333333,0.818181818 C5.5,0.818181818 6.29166667,1.39090909 6.54166667,2.25 C6.625,2.45454545 6.875,2.57727273 7.08333333,2.49545455 C7.10829568,2.4913698 7.13159638,2.48524575 7.15327691,2.47732675 C7.26611665,2.45454545 7.3971598,2.45454545 7.5,2.45454545 C8.41666667,2.45454545 9.16666667,3.19090909 9.16666667,4.09090909 C9.16666667,4.99090909 8.41666667,5.72727273 7.5,5.72727273 L6.66666667,5.72727273 C6.41666667,5.72727273 6.25,5.89090909 6.25,6.13636364 C6.25,6.38181818 6.41666667,6.54545455 6.66666667,6.54545455 L7.5,6.54545455 C8.875,6.54545455 10,5.44090909 10,4.09090909 C10,2.74090909 8.875,1.63636364 7.5,1.63636364 C7.40914201,1.63636364 7.31828402,1.64414451 7.22742604,1.65631288 C6.72946884,0.661035962 5.70689689,0 4.58333333,0 C3.08764506,0 1.87434996,1.07437766 1.69072185,2.48944595 C0.717443081,2.66537889 0,3.48663508 0,4.5 C0,5.64545455 0.916666667,6.54545455 2.08333333,6.54545455 L3.33333333,6.54545455 C3.58333333,6.54545455 3.75,6.38181818 3.75,6.13636364 C3.75,5.89090909 3.58333333,5.72727273 3.33333333,5.72727273 L2.08333333,5.72727273 C1.375,5.72727273 0.833333333,5.19545455 0.833333333,4.5 C0.833333333,3.80454545 1.375,3.27272727 2.08333333,3.27272727 C2.20833333,3.27272727 2.3125,3.23181818 2.38541667,3.16022727 Z M5,9 C4.75,9 4.58333333,8.83636364 4.58333333,8.59090909 L4.58333333,4.90909091 C4.58333333,4.66363636 4.75,4.5 5,4.5 C5.25,4.5 5.41666667,4.66363636 5.41666667,4.90909091 L5.41666667,8.59090909 C5.41666667,8.83636364 5.25,9 5,9 Z" sketch:type="MSShapeGroup"></path>
				<path d="M5,9 C4.875,9 4.79166667,8.95909091 4.70833333,8.87727273 C4.54166667,8.71363636 4.54166667,8.46818182 4.70833333,8.30454545 L5.95833333,7.07727273 C6.125,6.91363636 6.375,6.91363636 6.54166667,7.07727273 C6.70833333,7.24090909 6.70833333,7.48636364 6.54166667,7.65 L5.29166667,8.87727273 C5.20833333,8.95909091 5.125,9 5,9 L5,9 Z" sketch:type="MSShapeGroup"></path>
			</g>
		</g>
	</g>
</svg>
<p>Latest Release v1.0.6</p>
<div id="binariesList">
    <a href="https://aada-pet-werewolf-binaries.s3.amazonaws.com/aada_mac_x64.zip" download>aada_mac_x64.zip</a>
    <a href="https://aada-pet-werewolf-binaries.s3.amazonaws.com/aada_mac_m1.zip" download>aada_mac_m1.zip</a>
    <a href="https://aada-pet-werewolf-binaries.s3.amazonaws.com/aada_win_x64.zip" download>aada_win_x64.zip</a>
    <a href="https://aada-pet-werewolf-binaries.s3.amazonaws.com/aada_win_arm.zip" download>aada_win_arm.zip</a>
    <a href="https://aada-pet-werewolf-binaries.s3.amazonaws.com/aada_linux_x64.zip" download>aada_linux_x64.zip</a>
    <a href="https://aada-pet-werewolf-binaries.s3.amazonaws.com/aada_linux_arm32.zip" download>aada_linux_arm32.zip</a>
    <a href="https://aada-pet-werewolf-binaries.s3.amazonaws.com/aada_linux_arm64.zip" download>aada_linux_arm64.zip</a>
</div>
</body></html>
`

func buildDownloadsPage() Response {
	return Response{
		StatusCode: http.StatusOK,
		Body:       downloadsPage,
		Headers: map[string]string{
			"Content-Type": "text/html",
		},
		IsBase64Encoded: false,
	}
}
