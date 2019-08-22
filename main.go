package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

var (
	localURL = "http://127.0.0.1:9888/"

	buildTransactionURL  = localURL + "build-transaction"
	getTransactionURL    = localURL + "get-transaction"
	signTransactionURL   = localURL + "sign-transaction"
	submitTransactionURL = localURL + "submit-transaction"

	compileURL            = localURL + "compile"
	decodeProgramURL      = localURL + "decode-program"
	listAccountsURL       = localURL + "list-accounts"
	listAddressesURL      = localURL + "list-addresses"
	listBalancesURL       = localURL + "list-balances"
	listPubkeysURL        = localURL + "list-pubkeys"
	listUnspentOutputsURL = localURL + "list-unspent-outputs"
)

func main() {
	balances := listBalances("a1")
	fmt.Println(balances)
	listAccounts()
	addresses := listAddresses("a1")
	fmt.Println(addresses)
	pubkeyInfo := listPubkeys("a1")
	fmt.Println(pubkeyInfo)
}

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
	fmt.Println("response Body:", string(body))
	return body
}

type Account struct {
	AccountID    string `json:"id"`
	AccountAlias string `json:"alias"`
}

type Accounts struct {
	Status string    `json:"status"`
	Data   []Account `json:"data"`
}

func listAccounts() []Account {
	data := []byte(`{}`)
	body := request(listAccountsURL, data)

	accounts := new(Accounts)
	if err := json.Unmarshal(body, accounts); err != nil {
		fmt.Println(err)
	}
	return accounts.Data
}

type Address struct {
	AccountAlias   string `json:"account_alias"`
	AccountID      string `json:"account_id"`
	Address        string `json:"address"`
	ControlProgram string `json:"control_program"`
	Change         bool   `json:"change"`
	KeyIndex       uint64 `json:"key_index"`
}

type Addresses struct {
	Status string    `json:"status"`
	Data   []Address `json:"data"`
}

func listAddresses(accountAlias string) []Address {
	data := []byte(`{"account_alias": "` + accountAlias + `"}`)
	body := request(listAddressesURL, data)

	addresses := new(Addresses)
	if err := json.Unmarshal(body, addresses); err != nil {
		fmt.Println(err)
	}
	return addresses.Data
}

type Balance struct {
	AccountID string `json:"account_id"`
	Amount    uint64 `json:"amount"`
}

type Balances struct {
	Status string    `json:"status"`
	Data   []Balance `json:"data"`
}

func listBalances(accountAlias string) []Balance {
	data := []byte(`{"account_alias": "` + accountAlias + `"}`)
	body := request(listBalancesURL, data)

	balances := new(Balances)
	if err := json.Unmarshal(body, balances); err != nil {
		fmt.Println(err)
	}
	return balances.Data
}

type PubkeyInfo struct {
	Pubkey string   `json:"pubkey"`
	Path   []string `json:"derivation_path"`
}

type KeyInfo struct {
	XPubkey     string       `json:"root_xpub"`
	PubkeyInfos []PubkeyInfo `json:"pubkey_infos"`
}

type Pubkeys struct {
	Status string  `json:"status"`
	Data   KeyInfo `json:"data"`
}

func listPubkeys(accountAlias string) KeyInfo {
	data := []byte(`{"account_alias": "` + accountAlias + `"}`)
	body := request(listPubkeysURL, data)

	pubkeys := new(Pubkeys)
	if err := json.Unmarshal(body, pubkeys); err != nil {
		fmt.Println(err)
	}
	return pubkeys.Data
}
