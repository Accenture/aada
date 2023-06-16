package main

const (
	EventTypeConnect    = "CONNECT"
	EventTypeMessage    = "MESSAGE"
	EventTypeDisconnect = "DISCONNECT"
)

type Identity struct {
	SourceIp string `json:"sourceIp"`
}

type RequestContext struct {
	AccountId         string   `json:"acountId"`
	ApiId             string   `json:"apiId"`
	DomainName        string   `json:"domainName"`
	DomainPrefix      string   `json:"domainPrefix"`
	Path              string   `json:"path"`
	Protocol          string   `json:"protocol"`
	RouteKey          string   `json:"routeKey"`
	EventType         string   `json:"eventType"`
	ExtendedRequestId string   `json:"extendedRequestId"`
	ConnectionId      string   `json:"connectionId"`
	ResourceId        string   `json:"resourceId"`
	ResourcePath      string   `json:"resourcePath"`
	Stage             string   `json:"stage"`
	Identity          Identity `json:"identity"`
}

type Event struct {
	Version         string         `json:"version"`
	Resource        string         `json:"resource"`
	Path            string         `json:"path"`
	Context         RequestContext `json:"requestContext"`
	HttpMethod      string         `json:"httpMethod"`
	Body            string         `json:"body"`
	IsBase64Encoded bool           `json:"isBase64Encoded"`
}

type HTTPResponse struct {
	StatusCode      int               `json:"statusCode"`
	Body            string            `json:"body"`
	IsBase64Encoded bool              `json:"isBase64Encoded"`
	Headers         map[string]string `json:"headers"`
}
