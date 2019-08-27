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

// DeployContract deploy contract.
func DeployContract(assetRequested, seller, cancelKey, accountIDLocked, assetLocked, accountPasswordLocked string, amountRequested, amountLocked, txFee uint64) string {
	// compile locked contract
	contractInfo := CompileLockContract(assetRequested, seller, cancelKey, amountRequested)
	fmt.Println("--> contract info:", contractInfo)

	// build locked contract
	txLocked := BuildLockTransaction(accountIDLocked, assetLocked, contractInfo.Program, amountLocked, txFee)
	fmt.Println("--> txLocked:", string(txLocked))

	// sign locked contract transaction
	signedTransaction := SignTransaction(accountPasswordLocked, string(txLocked))
	fmt.Println("--> signedTransaction:", signedTransaction)

	// submit signed transaction
	txID := SubmitTransaction(signedTransaction)
	fmt.Println("--> txID:", txID)

	// get contract output ID
	contractUTXOID, err := GetContractUTXOID(txID, contractInfo.Program)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("--> contractUTXOID:", contractUTXOID)
	return contractUTXOID
}

// CallContract call contract.
func CallContract(accountIDUnlocked, contractUTXOID, seller, assetIDLocked, assetRequested, buyerContolProgram, accountPasswordUnlocked string, amountRequested, amountLocked, txFee uint64) string {
	// build unlocked contract transaction
	txUnlocked := BuildUnlockContractTransaction(accountIDUnlocked, contractUTXOID, seller, assetIDLocked, assetRequested, buyerContolProgram, amountRequested, amountLocked, txFee)
	fmt.Println("--> txUnlocked:", string(txUnlocked))

	// sign unlocked contract transaction
	signedTransaction := SignTransaction(accountPasswordUnlocked, string(txUnlocked))
	fmt.Println("--> signedTransaction:", signedTransaction)

	// submit signed unlocked contract transaction
	txID := SubmitTransaction(signedTransaction)
	fmt.Println("--> txID:", txID)
	return txID
}
