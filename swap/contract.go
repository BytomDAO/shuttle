package swap

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
)

var (
	errFailedGetContractUTXOID = errors.New("Failed to get contract UTXO ID")
)

type Account struct {
	AccountID    string `json:"id"`
	AccountAlias string `json:"alias"`
}

type AccountsResponse struct {
	Status string    `json:"status"`
	Data   []Account `json:"data"`
}

func ListAccounts() []Account {
	data := []byte(`{}`)
	body := request(listAccountsURL, data)

	accountsResp := new(AccountsResponse)
	if err := json.Unmarshal(body, accountsResp); err != nil {
		fmt.Println(err)
	}
	return accountsResp.Data
}

type Address struct {
	AccountAlias   string `json:"account_alias"`
	AccountID      string `json:"account_id"`
	Address        string `json:"address"`
	ControlProgram string `json:"control_program"`
	Change         bool   `json:"change"`
	KeyIndex       uint64 `json:"key_index"`
}

type AddressesResponse struct {
	Status string    `json:"status"`
	Data   []Address `json:"data"`
}

func ListAddresses(accountAlias string) []Address {
	data := []byte(`{"account_alias": "` + accountAlias + `"}`)
	body := request(listAddressesURL, data)

	addresses := new(AddressesResponse)
	if err := json.Unmarshal(body, addresses); err != nil {
		fmt.Println(err)
	}
	return addresses.Data
}

type Balance struct {
	AccountID string `json:"account_id"`
	Amount    uint64 `json:"amount"`
}

type BalancesResponse struct {
	Status string    `json:"status"`
	Data   []Balance `json:"data"`
}

func ListBalances(accountAlias string) []Balance {
	data := []byte(`{"account_alias": "` + accountAlias + `"}`)
	body := request(listBalancesURL, data)

	balancesResp := new(BalancesResponse)
	if err := json.Unmarshal(body, balancesResp); err != nil {
		fmt.Println(err)
	}
	return balancesResp.Data
}

type PubkeyInfo struct {
	Pubkey string   `json:"pubkey"`
	Path   []string `json:"derivation_path"`
}

type KeyInfo struct {
	XPubkey     string       `json:"root_xpub"`
	PubkeyInfos []PubkeyInfo `json:"pubkey_infos"`
}

type PubkeysResponse struct {
	Status string  `json:"status"`
	Data   KeyInfo `json:"data"`
}

func ListPubkeys(accountAlias string) KeyInfo {
	data := []byte(`{"account_alias": "` + accountAlias + `"}`)
	body := request(listPubkeysURL, data)

	pubkeysResp := new(PubkeysResponse)
	if err := json.Unmarshal(body, pubkeysResp); err != nil {
		fmt.Println(err)
	}
	return pubkeysResp.Data
}

type ContractInfo struct {
	Program string `json:"program"`
}

type ContractResponse struct {
	Status string       `json:"status"`
	Data   ContractInfo `json:"data"`
}

func CompileLockContract(assetRequested, seller, cancelKey string, amountRequested uint64) ContractInfo {
	data := []byte(`{
		"contract":"contract TradeOffer(assetRequested: Asset, amountRequested: Amount, seller: Program, cancelKey: PublicKey) locks valueAmount of valueAsset { clause trade() { lock amountRequested of assetRequested with seller unlock valueAmount of valueAsset } clause cancel(sellerSig: Signature) { verify checkTxSig(cancelKey, sellerSig) unlock valueAmount of valueAsset}}",
		"args":[
			{
				"string":"` + assetRequested + `"
			},
			{
				"integer":` + strconv.FormatUint(amountRequested, 10) + `
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

	contract := new(ContractResponse)
	if err := json.Unmarshal(body, contract); err != nil {
		fmt.Println(err)
	}
	return contract.Data
}

// type BuildTransactionResponse struct {
// 	Status string      `json:"status"`
// 	Data   interface{} `json:"data"`
// }

// BuildLockTransaction build locked contract transaction.
func BuildLockTransaction(accountIDLocked, assetIDLocked, contractControlProgram string, amountLocked, txFee uint64) []byte {
	data := []byte(`{
		"actions":[
			{
				"account_id":"` + accountIDLocked + `",
				"amount":` + strconv.FormatUint(amountLocked, 10) + `,
				"asset_id":"` + assetIDLocked + `",
				"type":"spend_account"
			},
			{
				"amount":` + strconv.FormatUint(amountLocked, 10) + `,
				"asset_id":"` + assetIDLocked + `",
				"control_program":"` + contractControlProgram + `",
				"type":"control_program"
			},
			{
				"account_id":"` + accountIDLocked + `",
				"amount":` + strconv.FormatUint(txFee, 10) + `,
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

type signedTransactionResponse struct {
	Status string          `json:"status"`
	Data   TransactionData `json:"data"`
}

// SignTransaction sign built contract transaction.
func SignTransaction(password, transaction string) string {
	data := []byte(`{
		"password": "` + password + `",
		"transaction` + transaction[25:])
	body := request(signTransactionURL, data)

	signedTransaction := new(signedTransactionResponse)
	if err := json.Unmarshal(body, signedTransaction); err != nil {
		fmt.Println(err)
	}
	return signedTransaction.Data.SignedTransaction.RawTransaction
}

type TransactionID struct {
	TxID string `json:"tx_id"`
}

type submitedTransactionResponse struct {
	Status string        `json:"status"`
	Data   TransactionID `json:"data"`
}

// SubmitTransaction submit raw singed contract transaction.
func SubmitTransaction(rawTransaction string) string {
	data := []byte(`{"raw_transaction": "` + rawTransaction + `"}`)
	body := request(submitTransactionURL, data)

	submitedTransaction := new(submitedTransactionResponse)
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

// GetContractUTXOID get contract UTXO ID by transaction ID and contract control program.
func GetContractUTXOID(transactionID, controlProgram string) (string, error) {
	data := []byte(`{"tx_id":"` + transactionID + `"}`)
	body := request(getTransactionURL, data)

	getTransactionResponse := new(getTransactionResponse)
	if err := json.Unmarshal(body, getTransactionResponse); err != nil {
		fmt.Println(err)
	}

	for _, v := range getTransactionResponse.Data.TransactionOutputs {
		if v.ControlProgram == controlProgram {
			return v.TransactionOutputID, nil
		}
	}

	return "", errFailedGetContractUTXOID
}

// BuildUnlockContractTransaction build unlocked contract transaction.
func BuildUnlockContractTransaction(accountIDUnlocked, contractUTXOID, seller, assetIDLocked, assetRequested, buyerContolProgram string, amountRequested, amountLocked, txFee uint64) []byte {
	data := []byte(`{
		"actions":[
			{
				"type":"spend_account_unspent_output",
				"arguments":[
					{
						"type":"integer",
						"raw_data":{
							"value":0
						}
					}
				],
				"use_unconfirmed":true,
				"output_id":"` + contractUTXOID + `"
			},
			{
				"amount":` + strconv.FormatUint(amountRequested, 10) + `,
				"asset_id":"` + assetRequested + `",
				"control_program":"` + seller + `",
				"type":"control_program"
			},
			{
				"account_id":"` + accountIDUnlocked + `",
				"amount":` + strconv.FormatUint(amountRequested, 10) + `,
				"asset_id":"` + assetRequested + `",
				"use_unconfirmed":true,
				"type":"spend_account"
			},
			{
				"account_id":"` + accountIDUnlocked + `",
				"amount":` + strconv.FormatUint(txFee, 10) + `,
				"asset_id":"ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff",
				"use_unconfirmed":true,
				"type":"spend_account"
			},
			{
				"amount":` + strconv.FormatUint(amountLocked, 10) + `,
				"asset_id":"` + assetIDLocked + `",
				"control_program":"` + buyerContolProgram + `",
				"type":"control_program"
			}
		],
		"ttl":0,
		"base_transaction":null
	}`)
	body := request(buildTransactionURL, data)
	return body
}
