package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

var (
	localURL         = "http://127.0.0.1:9888/"
	compileURL       = localURL + "compile"
	listAccountsURL  = localURL + "list-accounts"
	listAddressesURL = localURL + "list-addresses"
	listBalancesURL  = localURL + "list-balances"
)

func main() {
	listBalances()
	listAccounts()
	addresses := listAddresses("a1")
	fmt.Println(addresses)
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

func listBalances() {
	data := []byte(`{
		"account_alias": "a1"
		}`)
	body := request(listBalancesURL, data)

	balances := new(Balances)
	if err := json.Unmarshal(body, balances); err != nil {
		fmt.Println(err)
	}
	// fmt.Println("balances:", balances)
	// fmt.Println("status: ", balances.Status)
	// fmt.Println("amount: ", balances.Data[0].Amount)
}

func compile() {
	// data := []byte(`{
	// 	"contract":"contract TradeOffer(assetRequested: Asset, amountRequested: Amount, seller: Program, cancelKey: PublicKey) locks valueAmount of valueAsset { clause trade() { lock amountRequested of assetRequested with seller unlock valueAmount of valueAsset } clause cancel(sellerSig: Signature) { verify checkTxSig(cancelKey, sellerSig) unlock valueAmount of valueAsset}}",
	// 	"args":[
	// 		{
	// 			"string":"ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff"
	// 		},
	// 		{
	// 			"integer":1000000000
	// 		},
	// 		{
	// 			"string":"00145dd7b82556226d563b6e7d573fe61d23bd461c1f"
	// 		},
	// 		{
	// 			"string":"3e5d7d52d334964eef173021ef6a04dc0807ac8c41700fe718f5a80c2109f79e"
	// 		}
	// 	]
	// }`)
	// body := request(compileURL, data)

}

func builTransaction() {

}
