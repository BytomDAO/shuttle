package swap

import (
	"errors"
	"fmt"
	"strconv"
)

var (
	errFailedGetSignData = errors.New("Failed to get sign data")
)

type HTLCAccount struct {
	AccountID string
	Password  string
	Receiver  string
	TxFee     uint64
}

type HTLCContractArgs struct {
	SenderPublicKey    string
	RecipientPublicKey string
	BlockHeight        uint64
	Hash               string
}

type compileLockHTLCContractResponse struct {
	Program string `json:"program"`
}

var compileLockHTLCContractPayload = `{
    "contract":"contract HTLC(sender: PublicKey, recipient: PublicKey, blockHeight: Integer, hash: Hash) locks valueAmount of valueAsset { clause complete(preimage: String, sig: Signature) {verify sha256(preimage) == hash verify checkTxSig(recipient, sig) unlock valueAmount of valueAsset} clause cancel(sig: Signature) {verify above(blockHeight) verify checkTxSig(sender, sig) unlock valueAmount of valueAsset}}",
    "args":[
        {
            "string":"%s"
        },
        {
            "string":"%s"
        },
        {
            "integer":%s
        },
        {
            "string":"%s"
        }
    ]
}`

func compileLockHTLCContract(contractArgs HTLCContractArgs) (string, error) {
	payload := []byte(fmt.Sprintf(
		compileLockHTLCContractPayload,
		contractArgs.SenderPublicKey,
		contractArgs.RecipientPublicKey,
		strconv.FormatUint(contractArgs.BlockHeight, 10),
		contractArgs.Hash,
	))
	res := new(compileLockHTLCContractResponse)
	if err := request(compileURL, payload, res); err != nil {
		return "", err
	}
	return res.Program, nil
}

var buildLockHTLCContractTransactionPayload = `{
    "actions": [
        {
            "account_id": "%s",
            "amount": %s,
            "asset_id": "%s",
            "use_unconfirmed":true,
            "type": "spend_account"
        },
        {
            "amount": %s,
            "asset_id": "%s",
            "control_program": "%s",
            "type": "control_program"
        },
        {
            "account_id": "%s",
            "amount": %s,
            "asset_id": "ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff",
            "use_unconfirmed":true,
            "type": "spend_account"
        }
    ],
    "ttl": 0,
    "base_transaction": null
}`

func buildLockHTLCContractTransaction(account HTLCAccount, contractValue AssetAmount, contractControlProgram string) (interface{}, error) {
	payload := []byte(fmt.Sprintf(
		buildLockHTLCContractTransactionPayload,
		account.AccountID,
		strconv.FormatUint(contractValue.Amount, 10),
		contractValue.Asset,
		strconv.FormatUint(contractValue.Amount, 10),
		contractValue.Asset,
		contractControlProgram,
		account.AccountID,
		strconv.FormatUint(account.TxFee, 10),
	))
	res := new(interface{})
	if err := request(buildTransactionURL, payload, res); err != nil {
		return "", err
	}
	return res, nil
}

var buildUnlockHTLCContractTransactionPayload = `{
    "actions": [
        {
            "type": "spend_account_unspent_output",
            "use_unconfirmed":true,
            "arguments": [],
            "output_id": "%s"
        },
        {
            "account_id": "%s",
            "amount": %s,
            "asset_id": "ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff",
            "use_unconfirmed":true,
            "type": "spend_account"
        },
        {
            "amount": %s,
            "asset_id": "%s",
            "control_program": "%s",
            "type": "control_program"
        }
    ],
    "ttl": 0,
    "base_transaction": null
}`

func buildUnlockHTLCContractTransaction(account HTLCAccount, contractUTXOID string, contractArgs ContractArgs, contractValue AssetAmount) (interface{}, error) {
	payload := []byte(fmt.Sprintf(
		buildUnlockHTLCContractTransactionPayload,
		contractUTXOID,
		account.AccountID,
		strconv.FormatUint(account.TxFee, 10),
		strconv.FormatUint(contractValue.Amount, 10),
		contractValue.Asset,
		account.Receiver,
	))
	res := new(interface{})
	if err := request(buildTransactionURL, payload, res); err != nil {
		return "", err
	}
	return res, nil
}

type TransactionInput struct {
	ControlProgram string `json:"control_program"`
	SignData       string `json:"sign_data"`
}

type decodeRawTransactionResponse struct {
	TransactionInputs []TransactionInput `json:"inputs"`
}

var decodeRawTransactionPayload = `{
	"raw_transaction":"%s"
}`

func decodeRawTransaction(rawTransaction, controlProgram string) (string, error) {
	payload := []byte(fmt.Sprintf(
		decodeRawTransactionPayload,
		rawTransaction,
	))
	res := new(decodeRawTransactionResponse)
	if err := request(decodeRawTransactionPayload, payload, res); err != nil {
		return "", err
	}

	for _, v := range res.TransactionInputs {
		if v.ControlProgram == controlProgram {
			return v.SignData, nil
		}
	}
	return "", errFailedGetSignData
}

// DeployHTLCContract deploy HTLC contract.
func DeployHTLCContract(account HTLCAccount, contractValue AssetAmount, contractArgs HTLCContractArgs) (string, error) {
	// compile locked HTLC cotnract
	HTLCContractControlProgram, err := compileLockHTLCContract(contractArgs)
	if err != nil {
		return "", err
	}

	// build locked HTLC contract
	txLocked, err := buildLockHTLCContractTransaction(account, contractValue, HTLCContractControlProgram)
	if err != nil {
		return "", err
	}

	// sign locked HTLC contract transaction
	signedTransaction, err := signTransaction(account.Password, txLocked)
	if err != nil {
		return "", err
	}

	// submit signed HTLC contract transaction
	txID, err := submitTransaction(signedTransaction)
	if err != nil {
		return "", err
	}

	// get HTLC contract output ID
	contractUTXOID, err := getContractUTXOID(txID, HTLCContractControlProgram)
	if err != nil {
		return "", err
	}
	return contractUTXOID, nil
}
