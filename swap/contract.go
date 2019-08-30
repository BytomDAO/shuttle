package swap

import (
	"encoding/json"
	"errors"
	"strconv"
)

var (
	errFailedGetContractUTXOID = errors.New("Failed to get contract UTXO ID")
	errMarshal                 = errors.New("Failed to marshal")
)

type AccountInfo struct {
	AccountID string
	Password  string
	Receiver  string
	TxFee     uint64
}

type AssetAmount struct {
	Asset  string
	Amount uint64
}

type ContractArgs struct {
	AssetAmount
	Seller    string
	CancelKey string
}

type compileLockContractResponse struct {
	Program string `json:"program"`
}

// compileLockContract return contract control program
func compileLockContract(contractArgs ContractArgs) (string, error) {
	payload := []byte(`{
		"contract":"contract TradeOffer(assetRequested: Asset, amountRequested: Amount, seller: Program, cancelKey: PublicKey) locks valueAmount of valueAsset { clause trade() { lock amountRequested of assetRequested with seller unlock valueAmount of valueAsset } clause cancel(sellerSig: Signature) { verify checkTxSig(cancelKey, sellerSig) unlock valueAmount of valueAsset}}",
		"args":[
			{
				"string":"` + contractArgs.Asset + `"
			},
			{
				"integer":` + strconv.FormatUint(contractArgs.Amount, 10) + `
			},
			{
				"string":"` + contractArgs.Seller + `"
			},
			{
				"string":"` + contractArgs.CancelKey + `"
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
func buildLockTransaction(accountInfo AccountInfo, contractValue AssetAmount, contractControlProgram string) (interface{}, error) {
	payload := []byte(`{
		"actions":[
			{
				"account_id":"` + accountInfo.AccountID + `",
				"amount":` + strconv.FormatUint(contractValue.Amount, 10) + `,
				"asset_id":"` + contractValue.Asset + `",
				"type":"spend_account"
			},
			{
				"amount":` + strconv.FormatUint(contractValue.Amount, 10) + `,
				"asset_id":"` + contractValue.Asset + `",
				"control_program":"` + contractControlProgram + `",
				"type":"control_program"
			},
			{
				"account_id":"` + accountInfo.AccountID + `",
				"amount":` + strconv.FormatUint(accountInfo.TxFee, 10) + `,
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

// buildUnlockContractTransaction build unlocked contract transaction.
func buildUnlockContractTransaction(accountInfo AccountInfo, contractUTXOID string, contractArgs ContractArgs, contractValue AssetAmount) (interface{}, error) {
	payload := []byte(`{
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
				"amount":` + strconv.FormatUint(contractArgs.Amount, 10) + `,
				"asset_id":"` + contractArgs.Asset + `",
				"control_program":"` + contractArgs.Seller + `",
				"type":"control_program"
			},
			{
				"account_id":"` + accountInfo.AccountID + `",
				"amount":` + strconv.FormatUint(contractArgs.Amount, 10) + `,
				"asset_id":"` + contractArgs.Asset + `",
				"use_unconfirmed":true,
				"type":"spend_account"
			},
			{
				"account_id":"` + accountInfo.AccountID + `",
				"amount":` + strconv.FormatUint(accountInfo.TxFee, 10) + `,
				"asset_id":"ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff",
				"use_unconfirmed":true,
				"type":"spend_account"
			},
			{
				"amount":` + strconv.FormatUint(contractValue.Amount, 10) + `,
				"asset_id":"` + contractValue.Asset + `",
				"control_program":"` + accountInfo.Receiver + `",
				"type":"control_program"
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

// DeployContract deploy contract.
func DeployContract(accountInfo AccountInfo, contractArgs ContractArgs, contractValue AssetAmount) (string, error) {
	// compile locked contract
	contractControlProgram, err := compileLockContract(contractArgs)
	if err != nil {
		return "", err
	}

	// build locked contract
	txLocked, err := buildLockTransaction(accountInfo, contractValue, contractControlProgram)
	if err != nil {
		return "", err
	}

	// sign locked contract transaction
	signedTransaction, err := signTransaction(accountInfo.Password, txLocked)
	if err != nil {
		return "", err
	}

	// submit signed transaction
	txID, err := submitTransaction(signedTransaction)
	if err != nil {
		return "", err
	}

	// get contract output ID
	contractUTXOID, err := getContractUTXOID(txID, contractControlProgram)
	if err != nil {
		return "", err
	}
	return contractUTXOID, nil
}

// CallContract call contract.
func CallContract(accountInfo AccountInfo, contractUTXOID string, contractArgs ContractArgs, contractValue AssetAmount) (string, error) {
	// build unlocked contract transaction
	txUnlocked, err := buildUnlockContractTransaction(accountInfo, contractUTXOID, contractArgs, contractValue)
	if err != nil {
		return "", err
	}

	// sign unlocked contract transaction
	signedTransaction, err := signTransaction(accountInfo.Password, txUnlocked)
	if err != nil {
		return "", err
	}

	// submit signed unlocked contract transaction
	txID, err := submitTransaction(signedTransaction)
	if err != nil {
		return "", err
	}

	return txID, nil
}
