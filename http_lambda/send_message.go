package main

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	v4 "github.com/aws/aws-sdk-go-v2/aws/signer/v4"
	"github.com/aws/aws-sdk-go-v2/config"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"time"
)

func sendMessageToClient(ctx context.Context, gwurl string, client string, message string) error {
	rawBody := bytes.NewBufferString(message)

	hasher := sha256.New()
	hasher.Reset()
	hasher.Write(rawBody.Bytes())
	sum := hex.EncodeToString(hasher.Sum(nil))

	region, ok := os.LookupEnv("AWS_REGION")
	if !ok {
		region = "us-east-1" // Default only
	}

	// url is expected to be https://etc/@connections, so we're adding /client-id onto the end
	rawurl, err := url.Parse(gwurl + "/" + client)
	if err != nil {
		return err
	}

	fmt.Printf("posting to %s\n", rawurl.String())

	req := &http.Request{
		Method: "POST",
		URL:    rawurl,
		Header: make(http.Header),
		Body:   ioutil.NopCloser(rawBody),
	}
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return err
	}
	creds, err := cfg.Credentials.Retrieve(ctx)
	if err != nil {
		return err
	}
	signer := v4.NewSigner()
	err = signer.SignHTTP(ctx, creds, req, sum, "execute-api", region, time.Now())
	if err != nil {
		return err
	}

	rsp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Printf("response error %d - %s\n", rsp.StatusCode, rsp.Status)
		return err
	}
	content, _ := ioutil.ReadAll(rsp.Body)
	_ = rsp.Body.Close()
	fmt.Printf("response status %d - %s - %s\n", rsp.StatusCode, rsp.Status, string(content))
	return nil
}
