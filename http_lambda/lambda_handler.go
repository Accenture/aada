package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"net/http"
	"net/url"
	"os"
	"time"
)

// https://docs.aws.amazon.com/apigateway/latest/developerguide/set-up-lambda-proxy-integrations.html#api-gateway-simple-proxy-for-lambda-input-format
type Request struct {
	Method          string              `json:"httpMethod"`
	Path            string              `json:"path"`
	Body            string              `json:"body"`
	Headers         map[string]string   `json:"headers"`
	IsBase64Encoded bool                `json:"isBase64Encoded"`
	Query           map[string]string   `json:"queryStringParameters"`
	QuerySS         map[string][]string `json:"multiValueQueryStringParameters"`
}

// https://docs.aws.amazon.com/apigateway/latest/developerguide/set-up-lambda-proxy-integrations.html#api-gateway-simple-proxy-for-lambda-output-format
type Response struct {
	StatusCode int               `json:"statusCode"`
	Body       string            `json:"body"`
	Headers    map[string]string `json:"headers"`
	IsBase64Encoded bool `json:"isBase64Encoded"`
}

/// Interesting query
//
//  https://graph.microsoft.com/v1.0/me/transitiveMemberOf?$search="displayName:AABG.CIE.US.NA"
//  ConsistencyLevel: eventual
//

type Credentials struct {
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
	AccessToken string `json:"access_token"`
}

const NotAllowed = "{\"status\":\"denied\",\"message\":\"not allowed\"}"
const InvalidGroupName = "{\"status\":\"denied\",\"message\":\"invalid group name\"}"
const RoleAssumptionFailure = "{\"status\":\"denied\",\"message\":\"role assumption failure\"}"

func lambdaHandler(ctx context.Context, raw json.RawMessage) (Response, error) {
	rsp, err := internalLambdaHandler(ctx, raw)

	rsp.Headers["Cache-Control"] = "private, max-age=3600"
	rsp.Headers["Content-Security-Policy"] = "default-src 'self' aabg.io *.aabg.io; script-src 'self' 'unsafe-eval' 'unsafe-inline'; img-src 'self' data:; style-src 'self' 'unsafe-inline'; block-all-mixed-content;"
	rsp.Headers["Feature-Policy"] = "usb 'none'; geolocation 'none'; microphone 'none'; camera 'none'"
	rsp.Headers["Referrer-Policy"] = "strict-origin"
	rsp.Headers["Server"] = "ACN-httpd"
	rsp.Headers["Strict-Transport-Security"] = "max-age=31536000; includeSubdomains; preload" // 1 year
	rsp.Headers["X-Content-Type-Options"] = "NoSniff"
	rsp.Headers["X-Frame-Options"] = "SAMEORIGIN"
	rsp.Headers["X-XSS-Protection"] = "1; mode=block"

	fmt.Println("INFO HTTP response", rsp.StatusCode)

	return rsp, err
}

func internalLambdaHandler(ctx context.Context, raw json.RawMessage) (Response, error) {
	fmt.Println("INFO", string(raw))

	startTime := time.Now()
	defer func() {
		endTime := time.Now()
		fmt.Println("TIME", endTime.Sub(startTime).Milliseconds())
	}()

	in := Request{}
	err := json.Unmarshal(raw, &in)
	if err != nil {
		return Response{
			StatusCode: 500,
		}, err
	}

	host := in.Headers["X-Forwarded-Host"]
	if len(host) == 0 {
		host = in.Headers["X-Forwarded-For"]
	}
	if shouldThrottle(host) {
		fmt.Println("THROTTLING", host)
		return Response{
			StatusCode: 429,
		}, nil
	}

	wsurl, ok := os.LookupEnv("WS_CONN_URL")
	if !ok {
		return Response{StatusCode: http.StatusInternalServerError}, nil
	}

	switch in.Method {
	case "GET":
		switch in.Path {
		case "/favicon.ico":
			return Response{
				StatusCode: http.StatusOK,
				Body:       favicon,
				Headers: map[string]string{
					"Content-Type": "image/x-icon",
				},
				IsBase64Encoded: true,
			}, nil
		case "/downloads":
			return buildDownloadsPage(), nil
		case "/authenticator":
			code, ok := in.Query["code"]
			if !ok {
				return buildFailureResponse("missing request information"), nil
			}
			state, ok := in.Query["state"]
			if !ok {
				return buildFailureResponse("missing login state"), nil
			}

			if state == "awsconsole" || state == "awsconsole2" {
				return buildAWSConsoleDisplay(code)
			}
			return processOIDCRequest(ctx, state, code, "", wsurl)
		case "/awsconsole", "/awsconsole2":
			// Initiate a redirect for authentication
			nonce := uuid.NewString()
			rqv := url.Values{}
			rqv.Set("nonce", nonce)
			rqv.Set("state", in.Path[1:])
			rqv.Set("client_id", "dbf2de86-2e04-4086-bc86-bbc8b47076d5")
			rqv.Set("response_type", "code")
			rqv.Set("response_mode", "query")
			rqv.Set("scope", "openid profile email")
			rqv.Set("redirect_uri", "https://aabg.io/authenticator")
			return Response{
				StatusCode: http.StatusFound,
				Headers: map[string]string{
					"Location": authUrl + "?" + rqv.Encode(),
				},
			}, nil
		}
	case "POST":
		switch in.Path {
		case "/authenticator":
			// The post should definitely be encoded.  If not, there's something else wrong.
			if in.IsBase64Encoded {
				rq, err := base64.StdEncoding.DecodeString(in.Body)
				if err != nil {
					return buildFailureResponse("unable to decode response"), nil
				}
				rqv, err := url.ParseQuery(string(rq))
				if err != nil {
					return buildFailureResponse("unable to parse response"), nil
				}

				return processOIDCRequest(ctx, rqv.Get("state"), rqv.Get("code"), rqv.Get("id_token"), wsurl)
			}
		}
	}

	// If we got here, the client is doing something strange.  Slurp a few more tokens from the token bucket to slow
	// down bad actors.
	throttleCount(host, 5)

	return buildFailureResponse("malformed request"), nil
}
