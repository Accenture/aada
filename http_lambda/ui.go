package main

import (
	"context"
	_ "embed"
	"net/http"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

//go:embed favicon.b64
var favicon string

//go:embed ui_success.html
var successPage string

func buildSuccessResponse(message string) Response {
	return Response{
		StatusCode: http.StatusOK,
		Body:       strings.Replace(successPage, "__MESSAGE__", message, 1),
		Headers: map[string]string{
			"Content-Type": "text/html",
		},
	}
}

//go:embed ui_failure.html
var failurePage string

func buildFailureResponse(message string) Response {
	return Response{
		StatusCode: http.StatusBadRequest,
		Body:       strings.Replace(failurePage, "__MESSAGE__", message, 1),
		Headers: map[string]string{
			"Content-Type": "text/html",
		},
	}
}

//go:embed ui_default.html
var defaultPage string

func buildDefaultResponse() Response {
	return Response{
		StatusCode: http.StatusOK,
		Body:       defaultPage,
		Headers: map[string]string{
			"Content-Type": "text/html",
		},
	}
}

//go:embed ui_console.html
var consolePage string

type ConsoleLink struct {
	Url          string   `cbor:"1,keyasint,omitempty"`
	DisplayNames []string `cbor:"2,keyasint,omitempty"`
	Account      string   `cbor:"3,keyasint,omitempty"`
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

	// Sort by account number
	sort.Slice(links, func(i, j int) bool {
		return links[i].Account < links[j].Account
	})

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
	s3svc := s3.NewFromConfig(awsConfig)
	presigner := s3.NewPresignClient(s3svc)
	expiration := time.Now().Add(10 * time.Minute)
	bucketName, ok := os.LookupEnv("BINARIES_BUCKET")
	if !ok {
		bucketName = "invalid-bucket-configuration"
	}

	getObject := s3.GetObjectInput{
		Bucket:          aws.String(bucketName),
		ResponseExpires: &expiration,
	}
	binaries := []string{
		"aada_mac_x64",
		"aada_mac_arm64",
		"aada_win_x64",
		"aada_win_arm",
		"aada_linux_x64",
		"aada_linux_arm32",
		"aada_linux_arm64",
	}
	body := downloadsPage

	for _, key := range binaries {
		getObject.Key = aws.String(key + ".zip")
		url, _ := presigner.PresignGetObject(context.Background(), &getObject)
		body = strings.Replace(body, "{"+strings.ToUpper(key)+"}", url.URL, 1)
	}

	return Response{
		StatusCode: http.StatusOK,
		Body:       body,
		Headers: map[string]string{
			"Content-Type": "text/html",
		},
		IsBase64Encoded: false,
	}
}
