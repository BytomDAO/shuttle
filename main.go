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
)

type Balance struct {
	AccountID string `json:"account_id"`
	Amount    uint64 `json:"amount"`
	Other     interface{}
}

type Balances struct {
	Status string    `json:"status"`
	Data   []Balance `json:"data"`
}

func main() {
	url := localURL + "list-balances"
	fmt.Println("URL:>", url)

	var jsonStr = []byte(`{
		"account_alias": "a1"
		}`)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	// req.Header.Set("X-Custom-Header", "myvalue")
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	// fmt.Println("response Status:", resp.Status)
	// fmt.Println("response Headers:", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("response Body:", string(body))

	newBalances := new(Balances)
	if err := json.Unmarshal(body, newBalances); err != nil {
		fmt.Println(err)
	}
	fmt.Println("newBalances:", newBalances)
	fmt.Println("status: ", newBalances.Status)
	fmt.Println("amount: ", newBalances.Data[0].Amount)
}
