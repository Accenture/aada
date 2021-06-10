package main

import (
	"encoding/json"
	"github.com/briandowns/spinner"
	"github.com/fxamacker/cbor"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

func main() {
	query := "https://graph.microsoft.com/v1.0/groups?$filter=startswith(displayName,'AWS-')&$select=id,displayName"
	dumpGraphInfo(query, "groups.cbor")
	query = "https://graph.microsoft.com/v1.0/applications?$filter=startswith(displayName,'113614')&$select=id,displayName"
	dumpGraphInfo(query, "acp_apps.cbor")
}

func dumpGraphInfo(query string, fileName string) {
	doer := spinner.New(spinner.CharSets[14], 100*time.Millisecond)
	doer.Start()

	me := "Bearer " + os.Args[1]

	masterList := make(map[string]string)

	req, err := http.NewRequest("GET", query, nil)
	if err != nil {
		log.Fatal(err)
	}
	for i := 1; i < 10000; i++ {
		doer.Suffix = " " + strconv.Itoa(i)

		req.Header.Set("Authorization", me)
		rsp, err := http.DefaultClient.Do(req)
		if err != nil {
			log.Fatal(err)
		}

		body, err := ioutil.ReadAll(rsp.Body)
		if err != nil {
			log.Fatal(err)
		}
		err = rsp.Body.Close()
		if err != nil {
			log.Fatal(err)
		}

		pack := struct {
			NextLink string `json:"@odata.nextLink"`
			Value []struct{
				Id string `json:"id"`
				DisplayName string `json:"displayName"`
			}`json:"value"`
		}{}

		err = json.Unmarshal(body, &pack)
		if err != nil {
			log.Fatal(err)
		}

		for _, pv := range pack.Value {
			masterList[pv.Id] = pv.DisplayName
		}

		if pack.NextLink == "" {
			break
		}

		req, err = http.NewRequest("GET", pack.NextLink, nil)
	}

	doer.Suffix = " writing file"
	raw, err := cbor.Marshal(&masterList, cbor.CanonicalEncOptions())
	if err != nil {
		log.Fatal(err)
	}

	err = ioutil.WriteFile(fileName, raw, 0664)
	if err != nil {
		log.Fatal(err)
	}

	doer.Stop()
}
