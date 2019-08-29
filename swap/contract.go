package swap

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
)

var (
	errFailedGetContractUTXOID = errors.New("Failed to get contract UTXO ID")
	errMarshal                 = errors.New("Failed to marshal")
)

type compileLockContractResponse struct {
	Program string `json:"program"`
}

// compileLockContract return contract control program
func compileLockContract(assetRequested, seller, cancelKey string, amountRequested uint64) (string, error) {
	payload := []byte(`{
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
	res := new(compileLockContractResponse)
	if err := request(compileURL, payload, res); err != nil {
		return "", err
	}
	return res.Program, nil
}

// buildLockTransaction build locked contract transaction.
func buildLockTransaction(accountIDLocked, assetIDLocked, contractControlProgram string, amountLocked, txFee uint64) (interface{}, error) {
	payload := []byte(`{
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
	res := new(interface{})
	if err := request(buildTransactionURL, payload, res); err != nil {
		return "", err
	}
	return res, nil
}

type signTransactionRequest struct {
	Password    string      `json:"password"`
	Transaction interface{} `json:"transaction"`
}

type Transaction struct {
	RawTransaction string `json:"raw_transaction"`
}

type signTransactionResponse struct {
	Tx Transaction `json:"transaction"`
}

// signTransaction sign built contract transaction.
func signTransaction(password string, transaction interface{}) (string, error) {
	payload, err := json.Marshal(signTransactionRequest{Password: password, Transaction: transaction})
	if err != nil {
		return "", err
	}

	res := new(signTransactionResponse)
	if err := request(signTransactionURL, payload, res); err != nil {
		return "", err
	}

	return res.Tx.RawTransaction, nil
}

type submitTransactionRequest struct {
	RawTransaction string `json:"raw_transaction"`
}

type submitTransactionResponse struct {
	TransactionID string `json:"tx_id"`
}

// submitTransaction submit raw singed contract transaction.
func submitTransaction(rawTransaction string) (string, error) {
	payload, err := json.Marshal(submitTransactionRequest{RawTransaction: rawTransaction})
	if err != nil {
		return "", err
	}

	res := new(submitTransactionResponse)
	if err := request(submitTransactionURL, payload, res); err != nil {
		return "", err
	}

	return res.TransactionID, nil
}

type getContractUTXOIDRequest struct {
	TransactionID string `json:"tx_id"`
}

type TransactionOutput struct {
	TransactionOutputID string `json:"id"`
	ControlProgram      string `json:"control_program"`
}

type getContractUTXOIDResponse struct {
	TransactionOutputs []TransactionOutput `json:"outputs"`
}

// getContractUTXOID get contract UTXO ID by transaction ID and contract control program.
func getContractUTXOID(transactionID, controlProgram string) (string, error) {
	payload, err := json.Marshal(getContractUTXOIDRequest{TransactionID: transactionID})
	if err != nil {
		return "", err
	}

	res := new(getContractUTXOIDResponse)
	if err := request(getTransactionURL, payload, res); err != nil {
		return "", err
	}

	for _, v := range res.TransactionOutputs {
		if v.ControlProgram == controlProgram {
			return v.TransactionOutputID, nil
		}
	}

	return "", errFailedGetContractUTXOID
}

// // BuildUnlockContractTransaction build unlocked contract transaction.
// func BuildUnlockContractTransaction(accountIDUnlocked, contractUTXOID, seller, assetIDLocked, assetRequested, buyerContolProgram string, amountRequested, amountLocked, txFee uint64) []byte {
// 	data := []byte(`{
// 		"actions":[
// 			{
// 				"type":"spend_account_unspent_output",
// 				"arguments":[
// 					{
// 						"type":"integer",
// 						"raw_data":{
// 							"value":0
// 						}
// 					}
// 				],
// 				"use_unconfirmed":true,
// 				"output_id":"` + contractUTXOID + `"
// 			},
// 			{
// 				"amount":` + strconv.FormatUint(amountRequested, 10) + `,
// 				"asset_id":"` + assetRequested + `",
// 				"control_program":"` + seller + `",
// 				"type":"control_program"
// 			},
// 			{
// 				"account_id":"` + accountIDUnlocked + `",
// 				"amount":` + strconv.FormatUint(amountRequested, 10) + `,
// 				"asset_id":"` + assetRequested + `",
// 				"use_unconfirmed":true,
// 				"type":"spend_account"
// 			},
// 			{
// 				"account_id":"` + accountIDUnlocked + `",
// 				"amount":` + strconv.FormatUint(txFee, 10) + `,
// 				"asset_id":"ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff",
// 				"use_unconfirmed":true,
// 				"type":"spend_account"
// 			},
// 			{
// 				"amount":` + strconv.FormatUint(amountLocked, 10) + `,
// 				"asset_id":"` + assetIDLocked + `",
// 				"control_program":"` + buyerContolProgram + `",
// 				"type":"control_program"
// 			}
// 		],
// 		"ttl":0,
// 		"base_transaction":null
// 	}`)
// 	body := request(buildTransactionURL, data)
// 	return body
// }

// DeployContract deploy contract.
func DeployContract(assetRequested, seller, cancelKey, accountIDLocked, assetLocked, accountPasswordLocked string, amountRequested, amountLocked, txFee uint64) string {
	// compile locked contract
	contractControlProgram, err := compileLockContract(assetRequested, seller, cancelKey, amountRequested)
	if err != nil {
		panic(err)
	}
	fmt.Println("--> contractControlProgram:", contractControlProgram)

	// build locked contract
	txLocked, err := buildLockTransaction(accountIDLocked, assetLocked, contractControlProgram, amountLocked, txFee)
	if err != nil {
		panic(err)
	}
	fmt.Println("--> txLocked:", txLocked)

	// sign locked contract transaction
	signedTransaction, err := signTransaction(accountPasswordLocked, txLocked)
	if err != nil {
		panic(err)
	}
	fmt.Println("--> signedTransaction:", signedTransaction)

	// submit signed transaction
	txID, err := submitTransaction(signedTransaction)
	if err != nil {
		panic(err)
	}
	fmt.Println("--> txID:", txID)

	// get contract output ID
	contractUTXOID, err := getContractUTXOID(txID, contractControlProgram)
	if err != nil {
		panic(err)
	}
	fmt.Println("--> contractUTXOID:", contractUTXOID)
	return contractUTXOID
}

// // CallContract call contract.
// func CallContract(accountIDUnlocked, contractUTXOID, seller, assetIDLocked, assetRequested, buyerContolProgram, accountPasswordUnlocked string, amountRequested, amountLocked, txFee uint64) string {
// 	// build unlocked contract transaction
// 	txUnlocked := BuildUnlockContractTransaction(accountIDUnlocked, contractUTXOID, seller, assetIDLocked, assetRequested, buyerContolProgram, amountRequested, amountLocked, txFee)
// 	fmt.Println("--> txUnlocked:", string(txUnlocked))

// 	// sign unlocked contract transaction
// 	signedTransaction := SignTransaction(accountPasswordUnlocked, string(txUnlocked))
// 	fmt.Println("--> signedTransaction:", signedTransaction)

// 	// submit signed unlocked contract transaction
// 	txID := SubmitTransaction(signedTransaction)
// 	fmt.Println("--> txID:", txID)
// 	return txID
// }
