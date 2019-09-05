package swap

import (
	"fmt"
	"strconv"
)

type HTLCAccount struct {
	AccountID string
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

var buildLockHTLCContractPayload = `{
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

func buildLockHTLCContract(account HTLCAccount, contractValue AssetAmount, contractControlProgram string) (interface{}, error) {
	payload := []byte(fmt.Sprintf(
		buildLockHTLCContractPayload,
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

// func DeployHTLCContract(contractArgs HTLCContractArgs) (string, error) {
// 	// compile locked HTLC cotnract
// 	HTLCContractControlProgram, err := compileLockHTLCContract(contractArgs)
// 	if err != nil {
// 		return "", err
// 	}

// 	return
// }
