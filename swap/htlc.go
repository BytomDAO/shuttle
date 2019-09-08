package swap

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
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

// type SigningInstruction struct {
// 	Position          uint64        `json:"position"`
// 	WitnessComponents []interface{} `json:"witness_components"`
// }

type buildUnlockHTLCContractTransactionResponse struct {
	RawTransaction         string        `json:"raw_transaction"`
	SigningInstructions    []interface{} `json:"signing_instructions"`
	TxFee                  uint64        `json:"fee"`
	AllowAdditionalActions bool          `json:"allow_additional_actions"`
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

func buildUnlockHTLCContractTransaction(account HTLCAccount, contractUTXOID string, contractValue AssetAmount) (*buildUnlockHTLCContractTransactionResponse, error) {
	payload := []byte(fmt.Sprintf(
		buildUnlockHTLCContractTransactionPayload,
		contractUTXOID,
		account.AccountID,
		strconv.FormatUint(account.TxFee, 10),
		strconv.FormatUint(contractValue.Amount, 10),
		contractValue.Asset,
		account.Receiver,
	))
	res := new(buildUnlockHTLCContractTransactionResponse)
	if err := request(buildTransactionURL, payload, res); err != nil {
		return nil, err
	}
	// signingInst, err := json.Marshal(res.SigningInstructions[1])
	// if err != nil {
	// 	return nil, err
	// }
	// fmt.Println("signingInst:", string(signingInst))
	return res, nil
}

type TransactionInput struct {
	AssetID        string `json:"asset_id"`
	ControlProgram string `json:"control_program"`
	SignData       string `json:"sign_data"`
}

type decodeRawTransactionResponse struct {
	TransactionInputs []TransactionInput `json:"inputs"`
}

var decodeRawTransactionPayload = `{
	"raw_transaction":"%s"
}`

func decodeRawTransaction(rawTransaction string, contractValue AssetAmount) (string, string, error) {
	payload := []byte(fmt.Sprintf(
		decodeRawTransactionPayload,
		rawTransaction,
	))
	res := new(decodeRawTransactionResponse)
	if err := request(decodeRawTransactionURL, payload, res); err != nil {
		return "", "", err
	}

	for _, v := range res.TransactionInputs {
		if v.AssetID == contractValue.Asset {
			return v.ControlProgram, v.SignData, nil
		}
	}
	return "", "", errFailedGetSignData
}

func getRecipientPublicKey(contractControlProgram string) (string, error) {
	payload := []byte(fmt.Sprintf(
		decodeProgramPayload,
		contractControlProgram,
	))
	res := new(decodeProgramResponse)
	if err := request(decodeProgramURL, payload, res); err != nil {
		return "", err
	}

	publicKey := strings.Fields(res.Instructions)[5]
	return publicKey, nil
}

type AddressInfo struct {
	AccountAlias   string `json:"account_alias"`
	AccountID      string `json:"account_id"`
	Address        string `json:"address"`
	ControlProgram string `json:"control_program"`
}

var listAddressesPayload = `{
	"account_id":"%s"
}`

func listAddresses(accountID string) ([]AddressInfo, error) {
	payload := []byte(fmt.Sprintf(
		listAddressesPayload,
		accountID,
	))
	res := new([]AddressInfo)
	if err := request(listAddressesURL, payload, res); err != nil {
		return nil, err
	}

	return *res, nil
}

type signUnlockHTLCContractTransactionRequest struct {
	Password    string                                     `json:"password"`
	Transaction buildUnlockHTLCContractTransactionResponse `json:"transaction"`
}

var signUnlockHTLCContractTransactionPayload = `{
    "password": "%s",
    "transaction": {
        "raw_transaction": "%s",
        "signing_instructions": [
            {
                "position": 0,
                "witness_components": [
                    {
                        "type": "data",
                        "value": "%s"
                    },
                    {
                        "type": "data",
                        "value": "%s"
                    },
                    {
                        "type": "data",
                        "value": ""
                    }
                ]
            },
            %s
        ],
        "fee": %s,
        "allow_additional_actions": false
    }
}`

// func signUnlockHTLCContractTransaction(account HTLCAccount, preimage, recipientSig, rawTransaction string) (string, error) {
// 	rawSigningInstruction, err := json.Marshal(buildTxResp.SigningInstructions[1])
// 	if err != nil {
// 		return "", err
// 	}

// 	fmt.Println("rawSigningInstruction string:", string(rawSigningInstruction))

// 	payload := []byte(fmt.Sprintf(
// 		signUnlockHTLCContractTransactionPayload,
// 		account.Password,
// 		rawTransaction,
// 		preimage,
// 		recipientSig,
// 		string(rawSigningInstruction),
// 		strconv.FormatUint(account.TxFee, 10),
// 	))
// 	res := new(signTransactionResponse)
// 	if err := request(signTransactionURL, payload, res); err != nil {
// 		return "", err
// 	}

// 	return res.Tx.RawTransaction, nil
// }

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

// // CallHTLCContract call HTLC contract.
// func CallHTLCContract(account HTLCAccount, contractUTXOID, preimage, recipientSig string, contractArgs HTLCContractArgs, contractValue AssetAmount) (string, string, string, error) {
// 	// build unlocked contract transaction
// 	buildTxResp, err := buildUnlockHTLCContractTransaction(account, contractUTXOID, contractValue)
// 	if err != nil {
// 		return "", "", "", err
// 	}

// 	fmt.Println("raw transaction:", buildTxResp.RawTransaction)
// 	contractControlProgram, signData, err := decodeRawTransaction(buildTxResp.RawTransaction, contractValue)
// 	if err != nil {
// 		fmt.Println(err)
// 	}
// 	fmt.Println("contractControlProgram:", contractControlProgram)
// 	fmt.Println("signData:", signData)

// 	return buildTxResp.RawTransaction, contractControlProgram, signData, nil

// 	// // sign unlocked HTLC contract transaction
// 	// signedTransaction, err := signUnlockHTLCContractTransaction(account, preimage, recipientSig, *buildTxResp)
// 	// if err != nil {
// 	// 	return "", "", err
// 	// }

// 	// // submit signed HTLC contract transaction
// 	// txID, err := submitTransaction(signedTransaction)
// 	// if err != nil {
// 	// 	return "", "", err
// 	// }

// 	// return "", txID, nil
// }
