package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
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

func main() {
	balances := listBalances("a1")
	fmt.Println(balances)

	accounts := listAccounts()
	fmt.Println(accounts)

	addresses := listAddresses("a1")
	fmt.Println(addresses)

	pubkeyInfo := listPubkeys("a1")
	fmt.Println(pubkeyInfo)

	seller := "00145dd7b82556226d563b6e7d573fe61d23bd461c1f"
	cancelKey := "3e5d7d52d334964eef173021ef6a04dc0807ac8c41700fe718f5a80c2109f79e"
	contractInfo := compile(seller, cancelKey)
	fmt.Println(contractInfo)

	assetID := "bae7e17bb8f5d0cfbfd87a92f3204da082d388d4c9b10e8dcd36b3d0a18ceb3a"
	amount := uint64(20000000000)
	controlProgram := "203e5d7d52d334964eef173021ef6a04dc0807ac8c41700fe718f5a80c2109f79e1600145dd7b82556226d563b6e7d573fe61d23bd461c1f0400ca9a3b20ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff741a547a6413000000007b7b51547ac1631a000000547a547aae7cac00c0"
	tx := buildTransaction(assetID, controlProgram, amount)
	fmt.Println(tx)

	password := "12345"
	signedTransaction := signTransaction(password, string(tx))
	fmt.Println("signedTransaction:", signedTransaction)

	txID := submitTransaction(signedTransaction)
	fmt.Println("txID:", txID)

	txID = "5ebf9cb8d3e98450ba0c41a5cf60fd04da388c6da06addc2eb8e07265305a30a"
	contractUTXOID := getContractUTXOID(txID, controlProgram)
	fmt.Println("contractUTXOID:", contractUTXOID)
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

type ContractInfo struct {
	Program string `json:"program"`
}

type Contract struct {
	Status string       `json:"status"`
	Data   ContractInfo `json:"data"`
}

func compile(seller, cancelKey string) ContractInfo {
	data := []byte(`{
		"contract":"contract TradeOffer(assetRequested: Asset, amountRequested: Amount, seller: Program, cancelKey: PublicKey) locks valueAmount of valueAsset { clause trade() { lock amountRequested of assetRequested with seller unlock valueAmount of valueAsset } clause cancel(sellerSig: Signature) { verify checkTxSig(cancelKey, sellerSig) unlock valueAmount of valueAsset}}",
		"args":[
			{
				"string":"ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff"
			},
			{
				"integer":1000000000
			},
			{
				"string":"` + seller + `"
			},
			{
				"string":"` + cancelKey + `"
			}
		]
	}`)
	body := request(compileURL, data)

	contract := new(Contract)
	if err := json.Unmarshal(body, contract); err != nil {
		fmt.Println(err)
	}
	return contract.Data
}

type BuildTransaction struct {
	Status string      `json:"status"`
	Data   interface{} `json:"data"`
}

func buildTransaction(assetID, controlProgram string, amount uint64) []byte {
	data := []byte(`{
		"actions":[
			{
				"account_id":"10CJPO1HG0A02",
				"amount":` + strconv.FormatUint(amount, 10) + `,
				"asset_id":"` + assetID + `",
				"type":"spend_account"
			},
			{
				"amount":` + strconv.FormatUint(amount, 10) + `,
				"asset_id":"` + assetID + `",
				"control_program":"` + controlProgram + `",
				"type":"control_program"
			},
			{
				"account_id":"10CJPO1HG0A02",
				"amount":100000000,
				"asset_id":"ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff",
				"type":"spend_account"
			}
		],
		"ttl":0,
		"base_transaction":null
	}`)
	body := request(buildTransactionURL, data)
	return body
}

type SignedTransaction struct {
	RawTransaction string `json:"raw_transaction"`
}

type TransactionData struct {
	SignedTransaction SignedTransaction `json:"transaction"`
}

type signedTransaction struct {
	Status string          `json:"status"`
	Data   TransactionData `json:"data"`
}

func signTransaction(password, transaction string) string {
	data := []byte(`{
		"password": "` + password + `",
		"transaction` + transaction[25:])
	body := request(signTransactionURL, data)

	signedTransaction := new(signedTransaction)
	if err := json.Unmarshal(body, signedTransaction); err != nil {
		fmt.Println(err)
	}
	return signedTransaction.Data.SignedTransaction.RawTransaction
}

type TransactionID struct {
	TxID string `json:"tx_id"`
}

type submitedTransaction struct {
	Status string        `json:"status"`
	Data   TransactionID `json:"data"`
}

func submitTransaction(rawTransaction string) string {
	data := []byte(`{"raw_transaction": "` + rawTransaction + `"}`)
	body := request(submitTransactionURL, data)

	submitedTransaction := new(submitedTransaction)
	if err := json.Unmarshal(body, submitedTransaction); err != nil {
		fmt.Println(err)
	}
	return submitedTransaction.Data.TxID
}

type TransactionOutput struct {
	TransactionOutputID string `json:"id"`
	ControlProgram      string `json:"control_program"`
}

type GotTransactionInfo struct {
	TransactionOutputs []TransactionOutput `json:"outputs"`
}

type getTransactionResponse struct {
	Status string             `json:"status"`
	Data   GotTransactionInfo `json:"data"`
}

// getContractUTXOID get contract UTXO ID by transaction ID and control program.
func getContractUTXOID(transactionID, controlProgram string) string {
	data := []byte(`{"tx_id":"` + transactionID + `"}`)
	body := request(getTransactionURL, data)

	getTransactionResponse := new(getTransactionResponse)
	if err := json.Unmarshal(body, getTransactionResponse); err != nil {
		fmt.Println(err)
	}

	for _, v := range getTransactionResponse.Data.TransactionOutputs {
		if v.ControlProgram == controlProgram {
			return v.TransactionOutputID
		}
	}

	return ""
}
