package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
)

func loadGraphResultSet(creds *Credentials, query string) (<- chan json.RawMessage, <- chan error) {
	objs := make(chan json.RawMessage, 20)
	errs := make(chan error)

	go func() {
		defer close(objs)
		defer close(errs)

		req, err := http.NewRequest("GET", query, nil)
		if err != nil {
			errs <- errors.Wrap(err, "building request")
			return
		}

		for {
			req.Header.Add("Authorization", creds.TokenType+" "+creds.AccessToken)
			req.Header.Add("ConsistencyLevel", "eventual")

			rsp, err := http.DefaultClient.Do(req)
			if err != nil {
				errs <- errors.Wrap(err, "executing query")
				return
			}
			if rsp.StatusCode != 200 {
				/// DEBUGGING ONLY
				wtf, _ := httputil.DumpRequest(req, true)
				fmt.Println("DIAG REQUEST", base64.RawStdEncoding.EncodeToString(wtf))
				wtf, _ = httputil.DumpResponse(rsp, true)
				fmt.Println("DIAG RESPONSE", base64.RawStdEncoding.EncodeToString(wtf))

				errs <- errors.New(rsp.Status)
				return
			}
			raw, err := ioutil.ReadAll(rsp.Body)
			if err != nil {
				errs <- errors.Wrap(err, "unable to read response body")
				return
			}

			attr := struct {
				NextLink string            `json:"@odata.nextLink"`
				Values   []json.RawMessage `json:"value"`
			}{}
			err = json.Unmarshal(raw, &attr)
			if err != nil {
				errs <- errors.Wrap(err, "unable to unmarshal response")
				return
			}

			for _, value := range attr.Values {
				objs <- value
			}

			if attr.NextLink == "" {
				break
			}

			req, err = http.NewRequest("GET", attr.NextLink, nil)
			if err != nil {
				errs <- errors.Wrap(err, "failed to build subsequent request")
				return
			}
		}
	}()

	return objs, errs
}
