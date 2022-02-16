package main

import (
	_ "embed"
	"net/http"
	"sort"
	"strings"
)

//go:embed favicon.b64
var favicon string

//go:embed ui_success.html
var successPage string

func buildSuccessResponse() Response {
	return Response{
		StatusCode: http.StatusOK,
		Body:       successPage,
		Headers: map[string]string{
			"Content-Type": "text/html",
		},
	}
}

//go:embed ui_failure.html
var failurePage string

func buildFailureResponse(message string) Response {
	return Response{
		StatusCode: http.StatusOK,
		Body:       strings.Replace(failurePage, "__MESSAGE__", message, 1),
		Headers: map[string]string{
			"Content-Type": "text/html",
		},
	}
}

//go:embed ui_console.html
var consolePage string

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

//go:embed ui_downloads.html
var downloadsPage string

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
