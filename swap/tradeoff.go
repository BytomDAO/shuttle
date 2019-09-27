package swap

import (
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

var (
	errFailedGetContractUTXOID   = errors.New("Failed to get contract UTXO ID")
	errMarshal                   = errors.New("Failed to marshal")
	errListUnspentOutputs        = errors.New("Failed to list unspent outputs")
	errTradeOffParametersInvalid = errors.New("Trade off parameters invalid")
	errFailedSignTx              = errors.New("Failed to sign transaction")
	errFailedGetPublicKey        = errors.New("Failed to get public key")
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
			"integer":%d
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
func compileLockContract(s *Server, contractArgs ContractArgs) (string, error) {
	payload := []byte(fmt.Sprintf(compileLockContractReq,
		contractArgs.Asset,
		contractArgs.Amount,
		contractArgs.Seller,
		contractArgs.CancelKey,
	))
	res := new(compileLockContractResp)
	if err := s.request(compileURL, payload, res); err != nil {
		return "", err
	}
	return res.Program, nil
}

var buildLockTxReq = `{
	"actions":[
		{
			"account_id":"%s",
			"amount":%d,
			"asset_id":"%s",
			"use_unconfirmed":true,
			"type":"spend_account"
		},
		{
			"amount":%d,
			"asset_id":"%s",
			"control_program":"%s",
			"type":"control_program"
		},
		{
			"account_id":"%s",
			"amount":%d,
			"asset_id":"ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff",
			"use_unconfirmed":true,
			"type":"spend_account"
		}
	],
	"ttl":0,
	"base_transaction":null
}`

// buildLockTransaction build locked contract transaction.
func buildLockTransaction(s *Server, accountInfo AccountInfo, contractValue AssetAmount, contractControlProgram string) (interface{}, error) {
	payload := []byte(fmt.Sprintf(buildLockTxReq,
		accountInfo.AccountID,
		contractValue.Amount,
		contractValue.Asset,
		contractValue.Amount,
		contractValue.Asset,
		contractControlProgram,
		accountInfo.AccountID,
		accountInfo.TxFee,
	))
	res := new(interface{})
	if err := s.request(buildTransactionURL, payload, res); err != nil {
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
	Tx           Transaction `json:"transaction"`
	SignComplete bool        `json:"sign_complete"`
}

// signTransaction sign built contract transaction.
func signTransaction(s *Server, password string, transaction interface{}) (string, error) {
	payload, err := json.Marshal(signTxReq{Password: password, Transaction: transaction})
	if err != nil {
		return "", err
	}

	res := new(signTxResp)
	if err := s.request(signTransactionURL, payload, res); err != nil {
		return "", err
	}

	if !res.SignComplete {
		return "", errFailedSignTx
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
func submitTransaction(s *Server, rawTransaction string) (string, error) {
	payload, err := json.Marshal(submitTxReq{RawTransaction: rawTransaction})
	if err != nil {
		return "", err
	}

	res := new(submitTxResp)
	if err := s.request(submitTransactionURL, payload, res); err != nil {
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
func getContractUTXOID(s *Server, transactionID, controlProgram string) (string, error) {
	payload, err := json.Marshal(getContractUTXOIDReq{TransactionID: transactionID})
	if err != nil {
		return "", err
	}

	res := new(getContractUTXOIDResp)
	if err := s.request(getTransactionURL, payload, res); err != nil {
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
			"amount":%d,
			"asset_id":"%s",
			"control_program":"%s",
			"type":"control_program"
		},
		{
			"account_id":"%s",
			"amount":%d,
			"asset_id":"%s",
			"use_unconfirmed":true,
			"type":"spend_account"
		},
		{
			"account_id":"%s",
			"amount":%d,
			"asset_id":"ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff",
			"use_unconfirmed":true,
			"type":"spend_account"
		},
		{
			"amount":%d,
			"asset_id":"%s",
			"control_program":"%s",
			"type":"control_program"
		}
	],
	"ttl":0,
	"base_transaction":null
}`

// buildUnlockContractTransaction build unlocked contract transaction.
func buildUnlockContractTransaction(s *Server, accountInfo AccountInfo, contractUTXOID string) (interface{}, error) {
	program, contractValue, err := ListUnspentOutputs(s, contractUTXOID)
	if err != nil {
		return "", err
	}

	contractArgs, err := decodeProgram(s, program)
	if err != nil {
		return "", err
	}

	payload := []byte(fmt.Sprintf(buildUnlockContractTxReq,
		contractUTXOID,
		contractArgs.Amount,
		contractArgs.Asset,
		contractArgs.Seller,
		accountInfo.AccountID,
		contractArgs.Amount,
		contractArgs.Asset,
		accountInfo.AccountID,
		accountInfo.TxFee,
		contractValue.Amount,
		contractValue.Asset,
		accountInfo.Receiver,
	))
	res := new(interface{})
	if err := s.request(buildTransactionURL, payload, res); err != nil {
		return "", err
	}
	return res, nil
}

type listUnspentOutputsResp struct {
	AssetID     string `json:"asset_id"`
	AssetAmount uint64 `json:"amount"`
	Program     string `json:"program"`
}

type listUnspentOutputsReq struct {
	UTXOID        string `json:"id"`
	Unconfirmed   bool   `json:"unconfirmed"`
	SmartContract bool   `json:"smart_contract"`
}

func ListUnspentOutputs(s *Server, contractUTXOID string) (string, *AssetAmount, error) {
	payload, err := json.Marshal(listUnspentOutputsReq{
		UTXOID:        contractUTXOID,
		Unconfirmed:   true,
		SmartContract: true,
	})
	if err != nil {
		return "", nil, err
	}

	var res []listUnspentOutputsResp
	if err := s.request(listUnspentOutputsURL, payload, &res); err != nil {
		return "", nil, err
	}

	if len(res) == 0 {
		return "", nil, errListUnspentOutputs
	}

	contractLockedValue := new(AssetAmount)
	contractLockedValue.Asset = res[0].AssetID
	contractLockedValue.Amount = res[0].AssetAmount
	return res[0].Program, contractLockedValue, nil
}

type decodeProgramResp struct {
	Instructions string `json:"instructions"`
}

type decodeProgramReq struct {
	Program string `json:"program"`
}

func decodeProgram(s *Server, program string) (*ContractArgs, error) {
	payload, err := json.Marshal(decodeProgramReq{Program: program})
	if err != nil {
		return nil, err
	}

	res := new(decodeProgramResp)
	if err := s.request(decodeProgramURL, payload, res); err != nil {
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

	amount, err := ParseUint64(instructions[5])
	if err != nil {
		return nil, err
	}

	contractArgs.AssetAmount.Amount = amount
	return contractArgs, nil
}

func ParseUint64(s string) (uint64, error) {
	data, err := hex.DecodeString(s)
	if err != nil {
		return 0, err
	}

	var padded [8]byte
	copy(padded[:], data)
	num := binary.LittleEndian.Uint64(padded[:])

	return num, nil
}

type listPublicKeysReq struct {
	AccountID string `json:"account_id"`
}

type PubkeyInfo struct {
	PublicKey      string   `json:"pubkey"`
	DerivationPath []string `json:"derivation_path"`
}

type listPublicKeysResp struct {
	RootXPub    string       `json:"root_xpub"`
	PubkeyInfos []PubkeyInfo `json:"pubkey_infos"`
}

type XPubKeyInfo struct {
	XPubKey        string   `json:"xpub"`
	DerivationPath []string `json:"derivation_path"`
}

func getXPubKeyInfo(s *Server, accountID, publicKey string) (*XPubKeyInfo, error) {
	payload, err := json.Marshal(listPublicKeysReq{AccountID: accountID})
	if err != nil {
		return nil, err
	}

	res := new(listPublicKeysResp)
	if err := s.request(listPubkeysURL, payload, res); err != nil {
		return nil, err
	}

	xpubKeyInfo := new(XPubKeyInfo)
	xpubKeyInfo.XPubKey = res.RootXPub
	for _, PubkeyInfo := range res.PubkeyInfos {
		if PubkeyInfo.PublicKey == publicKey {
			xpubKeyInfo.DerivationPath = PubkeyInfo.DerivationPath
			return xpubKeyInfo, nil
		}
	}
	return nil, errFailedGetPublicKey
}

var buildCancelContractTxReq = `{
    "actions": [
        {
            "type": "spend_account_unspent_output",
            "arguments": [
                {
                    "type": "raw_tx_signature",
                    "raw_data": %s
                },
                {
                	"type":"integer",
                	"raw_data":{
                		"value":1
                	}
                }
			],
			"use_unconfirmed":true,
            "output_id": "%s"
        },
        {
            "account_id": "%s",
            "amount": %d,
            "asset_id": "ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff",
			"use_unconfirmed":true,
			"type": "spend_account"
        },
        {
            "amount": %d,
            "asset_id": "%s",
            "control_program": "%s",
            "type": "control_program"
        }
    ],
    "ttl": 0,
    "base_transaction": null
}`

func buildCancelContractTransaction(s *Server, accountInfo AccountInfo, contractUTXOID string, xpubKeyInfo *XPubKeyInfo, contractValue *AssetAmount) (interface{}, error) {
	xpubKeyInfoStr, err := json.Marshal(xpubKeyInfo)
	if err != nil {
		return "", err
	}
	payload := []byte(fmt.Sprintf(buildCancelContractTxReq,
		xpubKeyInfoStr,
		contractUTXOID,
		accountInfo.AccountID,
		accountInfo.TxFee,
		contractValue.Amount,
		contractValue.Asset,
		accountInfo.Receiver,
	))
	res := new(interface{})
	if err := s.request(buildTransactionURL, payload, res); err != nil {
		return "", err
	}
	return res, nil
}

// DeployTradeoffContract deploy contract.
func DeployTradeoffContract(s *Server, accountInfo AccountInfo, contractArgs ContractArgs, contractValue AssetAmount) (string, error) {
	// compile locked contract
	contractControlProgram, err := compileLockContract(s, contractArgs)
	if err != nil {
		return "", err
	}

	// build locked contract
	txLocked, err := buildLockTransaction(s, accountInfo, contractValue, contractControlProgram)
	if err != nil {
		return "", err
	}

	// sign locked contract transaction
	signedTransaction, err := signTransaction(s, accountInfo.Password, txLocked)
	if err != nil {
		return "", err
	}

	// submit signed transaction
	txID, err := submitTransaction(s, signedTransaction)
	if err != nil {
		return "", err
	}

	// get contract output ID
	contractUTXOID, err := getContractUTXOID(s, txID, contractControlProgram)
	if err != nil {
		return "", err
	}
	return contractUTXOID, nil
}

// CallTradeoffContract call contract.
func CallTradeoffContract(s *Server, accountInfo AccountInfo, contractUTXOID string) (string, error) {
	// build unlocked contract transaction
	txUnlocked, err := buildUnlockContractTransaction(s, accountInfo, contractUTXOID)
	if err != nil {
		return "", err
	}

	// sign unlocked contract transaction
	signedTransaction, err := signTransaction(s, accountInfo.Password, txUnlocked)
	if err != nil {
		return "", err
	}

	// submit signed unlocked contract transaction
	txID, err := submitTransaction(s, signedTransaction)
	if err != nil {
		return "", err
	}

	return txID, nil
}

// CancelTradeoffContract cancel tradeoff contract.
func CancelTradeoffContract(s *Server, accountInfo AccountInfo, contractUTXOID string) (string, error) {
	// get contract control program by contract UTXOID
	contractControlProgram, contractValue, err := ListUnspentOutputs(s, contractUTXOID)
	if err != nil {
		return "", err
	}

	// get public key by contract control program
	contractArgs, err := decodeProgram(s, contractControlProgram)
	if err != nil {
		return "", err
	}

	// get public key path and root xpub by contract args
	xpubInfo, err := getXPubKeyInfo(s, accountInfo.AccountID, contractArgs.CancelKey)
	if err != nil {
		return "", err
	}

	// build cancel contract transaction
	builtTx, err := buildCancelContractTransaction(s, accountInfo, contractUTXOID, xpubInfo, contractValue)
	if err != nil {
		return "", err
	}

	// sign cancel contract transaction
	signedTx, err := signTransaction(s, accountInfo.Password, builtTx)
	if err != nil {
		return "", err
	}

	// submit signed unlocked contract transaction
	txID, err := submitTransaction(s, signedTx)
	if err != nil {
		return "", err
	}

	return txID, nil
}
