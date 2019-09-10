package swap

import (
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

var (
	errFailedGetContractUTXOID   = errors.New("Failed to get contract UTXO ID")
	errMarshal                   = errors.New("Failed to marshal")
	errListUnspentOutputs        = errors.New("Failed to list unspent outputs")
	errTradeOffParametersInvalid = errors.New("Trade off parameters invalid")
)

type compileLockContractResp struct {
	Program string `json:"program"`
}

var compileLockContractReq = `{
	"contract":"contract TradeOffer(assetRequested: Asset, amountRequested: Amount, seller: Program, cancelKey: PublicKey) locks valueAmount of valueAsset { clause trade() { lock amountRequested of assetRequested with seller unlock valueAmount of valueAsset } clause cancel(sellerSig: Signature) { verify checkTxSig(cancelKey, sellerSig) unlock valueAmount of valueAsset}}",
	"args":[
		{
			"string":"%s"
		},
		{
			"integer":%s
		},
		{
			"string":"%s"
		},
		{
			"string":"%s"
		}
	]
}`

// compileLockContract return contract control program
func compileLockContract(contractArgs ContractArgs) (string, error) {
	payload := []byte(fmt.Sprintf(compileLockContractReq,
		contractArgs.Asset,
		strconv.FormatUint(contractArgs.Amount, 10),
		contractArgs.Seller,
		contractArgs.CancelKey,
	))
	res := new(compileLockContractResp)
	if err := request(compileURL, payload, res); err != nil {
		return "", err
	}
	return res.Program, nil
}

var buildLockTxReq = `{
	"actions":[
		{
			"account_id":"%s",
			"amount":%s,
			"asset_id":"%s",
			"use_unconfirmed":true,
			"type":"spend_account"
		},
		{
			"amount":%s,
			"asset_id":"%s",
			"control_program":"%s",
			"type":"control_program"
		},
		{
			"account_id":"%s",
			"amount":%s,
			"asset_id":"ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff",
			"use_unconfirmed":true,
			"type":"spend_account"
		}
	],
	"ttl":0,
	"base_transaction":null
}`

// buildLockTransaction build locked contract transaction.
func buildLockTransaction(accountInfo AccountInfo, contractValue AssetAmount, contractControlProgram string) (interface{}, error) {
	payload := []byte(fmt.Sprintf(buildLockTxReq,
		accountInfo.AccountID,
		strconv.FormatUint(contractValue.Amount, 10),
		contractValue.Asset,
		strconv.FormatUint(contractValue.Amount, 10),
		contractValue.Asset,
		contractControlProgram,
		accountInfo.AccountID,
		strconv.FormatUint(accountInfo.TxFee, 10),
	))
	res := new(interface{})
	if err := request(buildTransactionURL, payload, res); err != nil {
		return "", err
	}
	return res, nil
}

type signTxReq struct {
	Password    string      `json:"password"`
	Transaction interface{} `json:"transaction"`
}

type Transaction struct {
	RawTransaction string `json:"raw_transaction"`
}

type signTxResp struct {
	Tx Transaction `json:"transaction"`
}

// signTransaction sign built contract transaction.
func signTransaction(password string, transaction interface{}) (string, error) {
	payload, err := json.Marshal(signTxReq{Password: password, Transaction: transaction})
	if err != nil {
		return "", err
	}

	res := new(signTxResp)
	if err := request(signTransactionURL, payload, res); err != nil {
		return "", err
	}

	return res.Tx.RawTransaction, nil
}

type submitTxReq struct {
	RawTransaction string `json:"raw_transaction"`
}

type submitTxResp struct {
	TransactionID string `json:"tx_id"`
}

// submitTransaction submit raw singed contract transaction.
func submitTransaction(rawTransaction string) (string, error) {
	payload, err := json.Marshal(submitTxReq{RawTransaction: rawTransaction})
	if err != nil {
		return "", err
	}

	res := new(submitTxResp)
	if err := request(submitTransactionURL, payload, res); err != nil {
		return "", err
	}

	return res.TransactionID, nil
}

type getContractUTXOIDReq struct {
	TransactionID string `json:"tx_id"`
}

type TransactionOutput struct {
	TransactionOutputID string `json:"id"`
	ControlProgram      string `json:"control_program"`
}

type getContractUTXOIDResp struct {
	TransactionOutputs []TransactionOutput `json:"outputs"`
}

// getContractUTXOID get contract UTXO ID by transaction ID and contract control program.
func getContractUTXOID(transactionID, controlProgram string) (string, error) {
	payload, err := json.Marshal(getContractUTXOIDReq{TransactionID: transactionID})
	if err != nil {
		return "", err
	}

	res := new(getContractUTXOIDResp)
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

var buildUnlockContractTxReq = `{
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
			"output_id":"%s"
		},
		{
			"amount":%s,
			"asset_id":"%s",
			"control_program":"%s",
			"type":"control_program"
		},
		{
			"account_id":"%s",
			"amount":%s,
			"asset_id":"%s",
			"use_unconfirmed":true,
			"type":"spend_account"
		},
		{
			"account_id":"%s",
			"amount":%s,
			"asset_id":"ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff",
			"use_unconfirmed":true,
			"type":"spend_account"
		},
		{
			"amount":%s,
			"asset_id":"%s",
			"control_program":"%s",
			"type":"control_program"
		}
	],
	"ttl":0,
	"base_transaction":null
}`

// buildUnlockContractTransaction build unlocked contract transaction.
func buildUnlockContractTransaction(accountInfo AccountInfo, contractUTXOID string, contractArgs ContractArgs, contractValue AssetAmount) (interface{}, error) {
	payload := []byte(fmt.Sprintf(buildUnlockContractTxReq,
		contractUTXOID,
		strconv.FormatUint(contractArgs.Amount, 10),
		contractArgs.Asset,
		contractArgs.Seller,
		accountInfo.AccountID,
		strconv.FormatUint(contractArgs.Amount, 10),
		contractArgs.Asset,
		accountInfo.AccountID,
		strconv.FormatUint(accountInfo.TxFee, 10),
		strconv.FormatUint(contractValue.Amount, 10),
		contractValue.Asset,
		accountInfo.Receiver,
	))
	res := new(interface{})
	if err := request(buildTransactionURL, payload, res); err != nil {
		return "", err
	}
	return res, nil
}

type listUnspentOutputsResp struct {
	AssetID     string `json:"asset_id"`
	AssetAmount uint64 `json:"amount"`
	Program     string `json:"program"`
}

var listUnspentOutputsReq = `{
	"id": "%s",
	"unconfirmed": true,
	"smart_contract": true
}`

func ListUnspentOutputs(contractUTXOID string) (string, *AssetAmount, error) {
	payload := []byte(fmt.Sprintf(listUnspentOutputsReq, contractUTXOID))
	var res []listUnspentOutputsResp
	if err := request(listUnspentOutputsURL, payload, &res); err != nil {
		return "", nil, err
	}

	contractLockedValue := new(AssetAmount)
	if len(res) != 0 {
		contractLockedValue.Asset = res[0].AssetID
		contractLockedValue.Amount = res[0].AssetAmount
		return res[0].Program, contractLockedValue, nil
	}
	return "", nil, errListUnspentOutputs
}

type decodeProgramResp struct {
	Instructions string `json:"instructions"`
}

var decodeProgramPayload = `{
	"program": "%s"
}`

func DecodeProgram(program string) (*ContractArgs, error) {
	payload := []byte(fmt.Sprintf(decodeProgramPayload, program))
	res := new(decodeProgramResp)
	if err := request(decodeProgramURL, payload, res); err != nil {
		return nil, err
	}

	instructions := strings.Fields(res.Instructions)
	contractArgs := new(ContractArgs)
	contractArgs.CancelKey = instructions[1]
	contractArgs.Seller = instructions[3]
	contractArgs.AssetAmount.Asset = instructions[7]
	if len(contractArgs.CancelKey) != 64 || len(contractArgs.AssetAmount.Asset) != 64 {
		return nil, errTradeOffParametersInvalid
	}

	amount, err := parseUint64(instructions[5])
	if err != nil {
		return nil, err
	}

	contractArgs.AssetAmount.Amount = amount
	return contractArgs, nil
}

func parseUint64(s string) (uint64, error) {
	data, err := hex.DecodeString(s)
	if err != nil {
		return 0, err
	}

	var padded [8]byte
	copy(padded[:], data)
	num := binary.LittleEndian.Uint64(padded[:])

	return num, nil
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
