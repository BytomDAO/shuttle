package swap

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
)

var (
	localURL = "http://127.0.0.1:9888/"

	buildTransactionURL     = localURL + "build-transaction"
	getTransactionURL       = localURL + "get-transaction"
	signTransactionURL      = localURL + "sign-transaction"
	decodeRawTransactionURL = localURL + "decode-raw-transaction"
	submitTransactionURL    = localURL + "submit-transaction"
	compileURL              = localURL + "compile"
	decodeProgramURL        = localURL + "decode-program"
	listAccountsURL         = localURL + "list-accounts"
	listAddressesURL        = localURL + "list-addresses"
	listBalancesURL         = localURL + "list-balances"
	listPubkeysURL          = localURL + "list-pubkeys"
	listUnspentOutputsURL   = localURL + "list-unspent-outputs"
)

type response struct {
	Status    string          `json:"status"`
	Data      json.RawMessage `json:"data"`
	ErrDetail string          `json:"error_detail"`
}

func request(url string, payload []byte, respData interface{}) error {
	resp := &response{}
	if err := post(url, payload, resp); err != nil {
		return err
	}

	if resp.Status != "success" {
		return errors.New(resp.ErrDetail)
	}

	return json.Unmarshal(resp.Data, respData)
}

func post(url string, payload []byte, result interface{}) error {
	return PostWithHeader(url, nil, payload, result)
}

func PostWithHeader(url string, header map[string]string, payload []byte, result interface{}) error {
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payload))
	if err != nil {
		return err
	}

	// set Content-Type in advance, and overwrite Content-Type if provided
	req.Header.Set("Content-Type", "application/json")
	for k, v := range header {
		req.Header.Set(k, v)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	if result == nil {
		return nil
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	return json.Unmarshal(body, result)
}
