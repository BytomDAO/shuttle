package swap

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
)

var (
	localURL = "http://127.0.0.1:9888/"

	buildTransactionURL   = localURL + "build-transaction"
	getTransactionURL     = localURL + "get-transaction"
	signTransactionURL    = localURL + "sign-transaction"
	submitTransactionURL  = localURL + "submit-transaction"
	compileURL            = localURL + "compile"
	decodeProgramURL      = localURL + "decode-program"
	listAccountsURL       = localURL + "list-accounts"
	listAddressesURL      = localURL + "list-addresses"
	listBalancesURL       = localURL + "list-balances"
	listPubkeysURL        = localURL + "list-pubkeys"
	listUnspentOutputsURL = localURL + "list-unspent-outputs"
)

func request(URL string, data []byte) []byte {
	req, err := http.NewRequest("POST", URL, bytes.NewBuffer(data))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}
	// fmt.Println("response Body:", string(body))
	return body
}
