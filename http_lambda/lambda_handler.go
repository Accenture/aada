package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
)

// https://docs.aws.amazon.com/apigateway/latest/developerguide/set-up-lambda-proxy-integrations.html#api-gateway-simple-proxy-for-lambda-input-format
type Request struct {
	Method          string              `json:"httpMethod"`
	Path            string              `json:"path"`
	Body            string              `json:"body"`
	IsBase64Encoded bool                `json:"isBase64Encoded"`
	Query           map[string]string   `json:"queryStringParameters"`
	QuerySS         map[string][]string `json:"multiValueQueryStringParameters"`
}

// https://docs.aws.amazon.com/apigateway/latest/developerguide/set-up-lambda-proxy-integrations.html#api-gateway-simple-proxy-for-lambda-output-format
type Response struct {
	StatusCode int               `json:"statusCode"`
	Body       string            `json:"body"`
	Headers    map[string]string `json:"headers"`
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
	fmt.Println(string(raw))
	in := Request{}
	err := json.Unmarshal(raw, &in)
	if err != nil {
		return Response{
			StatusCode: 500,
		}, err
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
			}, nil
		case "/authenticator":
			code, ok := in.Query["code"]
			if !ok {
				return buildFailureResponse("missing request information"), nil
			}
			state, ok := in.Query["state"]
			if !ok {
				return buildFailureResponse("missing login state"), nil
			}
			return processOIDCRequest(ctx, state, code, "", wsurl)
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
	return buildFailureResponse("malformed request"), nil
}