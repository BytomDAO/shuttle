package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

var (
	localURL        = "http://127.0.0.1:9888/"
	listAccountsURL = localURL + "list-accounts"
	listBalancesURL = localURL + "list-balances"
)

type Balance struct {
	AccountID string `json:"account_id"`
	Amount    uint64 `json:"amount"`
}

type Balances struct {
	Status string    `json:"status"`
	Data   []Balance `json:"data"`
}

func main() {
	listBalances()
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

	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("response Body:", string(body))
	return body
}

func listAccounts() {

}

func listBalances() {
	var data = []byte(`{
		"account_alias": "a1"
		}`)
	body := request(listBalancesURL, data)

	balances := new(Balances)
	if err := json.Unmarshal(body, balances); err != nil {
		fmt.Println(err)
	}
	fmt.Println("balances:", balances)
	fmt.Println("status: ", balances.Status)
	fmt.Println("amount: ", balances.Data[0].Amount)
}
