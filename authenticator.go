package main

import (
	"context"
	"encoding/base64"
	"encoding/xml"
	"github.com/PuerkitoBio/goquery"
	"github.com/anaskhan96/soup"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/http/httputil"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

const federationUrl = "https://federation-sts.accenture.com/adfs/ls/IdpInitiatedSignOn.aspx?loginToRp=urn:amazon:webservices"

var Trace *os.File

type Authenticator struct {
	Username string
	Password string
	VIPCode  string

	client            *http.Client
	submissionFormUrl string
	params            url.Values
	samlToken         string

	AWSConfig aws.Config

	accessKey    string
	secretKey    string
	sessionToken string
	expiration   time.Time
}

// Requests the primary login page and extracts the unique client request for this login session
func (a *Authenticator) getSubmissionFormUrl() {
	// There must be a cookie jar in place to carry the session context through from request to request
	jar, err := cookiejar.New(nil)
	if err != nil {
		log.Fatal(err)
	}
	a.client = &http.Client{
		Jar: jar,
	}

	req, err := http.NewRequest("get", federationUrl, nil)
	if err != nil {
		log.Fatal(err)
	}

	dumpRequest(req, false)

	rsp, err := a.client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	dumpResponse(rsp)

	responseBody, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		log.Fatal(err)
	}
	_ = rsp.Body.Close()

	tree := soup.HTMLParse(string(responseBody))
	node := tree.Find("form", "id", "options")
	a.submissionFormUrl = node.Attrs()["action"]
}

// Submit the username and password to the login url in exchange for a login context
func (a *Authenticator) submitCredentials() {
	a.params = make(url.Values)
	a.params.Set("AuthMethod", "")
	a.params.Set("Kmsi", "true")
	a.params.Set("UserName", a.Username)
	a.params.Set("Password", a.Password)
	body := a.params.Encode()

	rsp := a.postRequest(a.submissionFormUrl, body, false)
	if rsp.StatusCode != http.StatusOK {
		log.Fatal("failed to authenticate, check username or password")
	}

	a.params = make(url.Values) // Clears the existing parameters

	doc, err := goquery.NewDocumentFromReader(rsp.Body)
	if err != nil {
		log.Fatal(err)
	}
	rsp.Body.Close()

	selector := doc.Find("span[id=errorText]")
	if len(selector.Nodes) > 0 {
		errorMessage := strings.TrimSpace(selector.Text())
		if errorMessage != "" {
			log.Fatal(selector.Text())
		}
	}

	selector = doc.Find("form[id=loginForm] div input")
	selector.Each(func(i int, selection *goquery.Selection) {
		a.params.Set(
			selection.AttrOr("name", "Unknown"),
			selection.AttrOr("value", "Unknown"),
		)
	})
}

// Submit the VIP security code to the login context
func (a *Authenticator) submitVipCode() {
	a.params.Set("SecurityCode", a.VIPCode)

	body := a.params.Encode()

	rsp := a.postRequest(a.submissionFormUrl, body, true)
	if rsp.StatusCode != http.StatusOK {
		log.Fatal("failed to submit Symantec VIP security code")
	}

	doc, err := goquery.NewDocumentFromReader(rsp.Body)
	if err != nil {
		log.Fatal(err)
	}
	rsp.Body.Close()

	selection := doc.Find("label[id=errorText]")
	errorText := selection.Text()
	if len(errorText) > 0 {
		log.Fatal(errorText)
	}

	selection = doc.Find("input[name=SAMLResponse]")
	if encodedToken, ok := selection.Attr("value"); ok {
		a.samlToken = encodedToken
	} else {
		log.Fatal("failed to retrieve SAML response")
	}
}

func (a *Authenticator) getSAMLRoles() []string {
	decodedToken, err := base64.StdEncoding.DecodeString(a.samlToken)
	if err != nil {
		log.Fatal(err)
	}

	samlDocument := &SAMLXml{}
	err = xml.Unmarshal(decodedToken, samlDocument)
	if err != nil {
		if err != io.EOF {
			log.Fatal(err)
		}
	}

	for _, attr := range samlDocument.Attributes {
		if attr.Name == "https://aws.amazon.com/SAML/Attributes/Role" {
			return attr.Value
		}
	}

	return []string{}
}

func (a *Authenticator) assumeRole(principal string, role string, duration time.Duration) {
	svc := sts.New(a.AWSConfig)
	rsp, err := svc.AssumeRoleWithSAMLRequest(&sts.AssumeRoleWithSAMLInput{
		DurationSeconds: aws.Int64(duration.Milliseconds() / 1000),
		PrincipalArn:    aws.String(principal),
		RoleArn:         aws.String(role),
		SAMLAssertion:   aws.String(a.samlToken),
	}).Send(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	a.accessKey = *rsp.Credentials.AccessKeyId
	a.secretKey = *rsp.Credentials.SecretAccessKey
	a.sessionToken = *rsp.Credentials.SessionToken
	a.expiration = *rsp.Credentials.Expiration
}

func (a *Authenticator) postRequest(url string, body string, dumpBody bool) *http.Response {
	req, err := http.NewRequest("POST", url, strings.NewReader(body))
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("Content-Length", strconv.Itoa(len(body)))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	dumpRequest(req, dumpBody)

	rsp, err := a.client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	dumpResponse(rsp)

	return rsp
}

func dumpResponse(response *http.Response) {
	if Trace == nil {
		return
	}
	Trace.WriteString("=== RESPONSE =============================================================================\n")
	raw, err := httputil.DumpResponse(response, true)
	if err != nil {
		log.Fatal(err)
	}
	Trace.Write(raw)
}

func dumpRequest(request *http.Request, dumpBody bool) {
	if Trace == nil {
		return
	}
	Trace.WriteString("--- REQUEST ------------------------------------------------------------------------------\n")
	raw, err := httputil.DumpRequest(request, dumpBody)
	if err != nil {
		log.Fatal(err)
	}
	Trace.Write(raw)
	if !dumpBody {
		Trace.WriteString(" *** body redacted ***\n")
	}
}
