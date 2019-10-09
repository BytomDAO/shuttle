package swap

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/bytom/crypto/ed25519/chainkd"
)

var (
	errXPrvLength              = errors.New("XPrv length is invalid.")
	errFailedGetContractUTXOID = errors.New("Failed to get contract UTXO ID")
)

type TxOutput struct {
	UTXOID string `json:"utxo_id"`
	Script string `json:"script"`
}

type getTxReq struct {
	TxID string `json:"tx_id"`
}

type getTxResp struct {
	TxOutputs []TxOutput `json:"outputs"`
}

// getUTXOID get UTXO ID by transaction ID.
func getUTXOID(s *Server, txID, controlProgram string) (string, error) {
	payload, err := json.Marshal(getTxReq{TxID: txID})
	if err != nil {
		return "", err
	}

	res := new(getTxResp)
	if err := s.request(getTransactionURL, payload, res); err != nil {
		return "", err
	}

	for _, v := range res.TxOutputs {
		if v.Script == controlProgram {
			return v.UTXOID, nil
		}
	}

	return "", errFailedGetContractUTXOID
}

type SigningInstruction struct {
	DerivationPath []string `json:"derivation_path"`
	SignData       []string `json:"sign_data"`
	DataWitness    []byte   `json:"-"`

	// only shown for a single-signature tx
	Pubkey string `json:"pubkey,omitempty"`
}

type SpendUTXOInput struct {
	Type     string `json:"type"`
	OutputID string `json:"output_id"`
}

type SpendWalletInput struct {
	Type    string `json:"type"`
	AssetID string `json:"asset"`
	Amount  uint64 `json:"amount"`
}

type ControlAddressOutput struct {
	Type    string `json:"type"`
	Amount  uint64 `json:"amount"`
	AssetID string `json:"asset"`
	Address string `json:"address"`
}

type ControlProgramOutput struct {
	Type           string `json:"type"`
	Amount         uint64 `json:"amount"`
	AssetID        string `json:"asset"`
	ControlProgram string `json:"control_program"`
}

type buildTxReq struct {
	GUID          string        `json:"guid"`
	Fee           uint64        `json:"fee"`
	Confirmations uint64        `json:"confirmations"`
	Inputs        []interface{} `json:"inputs"`
	Outputs       []interface{} `json:"outputs"`
}

type buildTxResp struct {
	RawTx               string                `json:"raw_transaction"`
	SigningInstructions []*SigningInstruction `json:"signing_instructions"`
	Fee                 uint64                `json:"fee"`
}

// BuildTx build tx.
func BuildTx(s *Server, guid, lockedAsset, contractProgram string, lockedAmount uint64) (string, error) {
	// inputs:
	spendWalletInput := SpendWalletInput{
		Type:    "spend_wallet",
		AssetID: lockedAsset,
		Amount:  lockedAmount,
	}

	// outputs:
	controlProgramOutput := ControlProgramOutput{
		Type:           "control_program",
		Amount:         lockedAmount,
		AssetID:        lockedAsset,
		ControlProgram: contractProgram,
	}

	var inputs, outputs []interface{}
	inputs = append(inputs, spendWalletInput)
	outputs = append(outputs, controlProgramOutput)
	payload, err := json.Marshal(buildTxReq{
		GUID:          guid,
		Fee:           fee,
		Confirmations: confirmations,
		Inputs:        inputs,
		Outputs:       outputs,
	})
	if err != nil {
		return "", err
	}

	fmt.Println("buildTx:", string(payload))

	res := new(buildTxResp)
	if err := s.request(buildTransactionURL, payload, res); err != nil {
		return "", err
	}

	r, err := json.MarshalIndent(res, "", "\t")
	if err != nil {
		return "", err
	}

	return string(r), nil
}

type submitPaymentReq struct {
	GUID       string     `json:"guid"`
	RawTx      string     `json:"raw_transaction"`
	Signatures [][]string `json:"signatures"`
	Memo       string     `json:"memo"`
}

type submitPaymentResp struct {
	TxID string `json:"transaction_hash"`
}

// SubmitPayment submit raw transaction and return transaction ID.
func SubmitPayment(s *Server, guid, rawTx, memo string, sigs [][]string) (string, error) {
	payload, err := json.MarshalIndent(submitPaymentReq{
		GUID:       guid,
		RawTx:      rawTx,
		Signatures: sigs,
		Memo:       memo,
	}, "", "\t")
	if err != nil {
		return "", err
	}

	fmt.Println("\nsubmitPayment:", string(payload))

	res := new(submitPaymentResp)
	if err := s.request(submitTransactionURL, payload, res); err != nil {
		return "", err
	}

	return res.TxID, nil
}

// SignMessage sign message, return sig.
func SignMessage(signData, xprv string) (string, error) {
	xprvBytes, err := hex.DecodeString(xprv)
	if err != nil {
		return "", err
	}
	if len(xprvBytes) != 64 {
		return "", errXPrvLength
	}

	var newXPrv chainkd.XPrv
	copy(newXPrv[:], xprvBytes[:])

	msg, err := hex.DecodeString(signData)
	if err != nil {
		return "", err
	}
	sig := newXPrv.Sign(msg)
	return hex.EncodeToString(sig), nil
}

// BuildUnlockedTx build unlocked contract tx.
func BuildUnlockedTx(s *Server, guid, contractUTXOID, contractAsset, receiver string, contractAmount uint64) (string, error) {
	// inputs:
	spendUTXOInput := SpendUTXOInput{
		Type:     "spend_utxo",
		OutputID: contractUTXOID,
	}

	// spendWalletInput := SpendWalletInput{
	// 	Type:    "spend_wallet",
	// 	AssetID: contractAsset,
	// 	Amount:  contractAmount,
	// }

	// outputs:
	controlAddressOutput := ControlAddressOutput{
		Type:    "control_address",
		Amount:  contractAmount,
		AssetID: contractAsset,
		Address: receiver,
	}

	var inputs, outputs []interface{}
	// inputs = append(inputs, spendUTXOInput, spendWalletInput)
	inputs = append(inputs, spendUTXOInput)
	outputs = append(outputs, controlAddressOutput)
	payload, err := json.Marshal(buildTxReq{
		GUID:          guid,
		Fee:           fee,
		Confirmations: confirmations,
		Inputs:        inputs,
		Outputs:       outputs,
	})
	if err != nil {
		return "", err
	}

	fmt.Println("build unlocked contract tx:", string(payload))

	res := new(buildTxResp)
	if err := s.request(buildTransactionURL, payload, res); err != nil {
		return "", err
	}

	r, err := json.MarshalIndent(res, "", "\t")
	if err != nil {
		return "", err
	}

	return string(r), nil
}

// BuildCallTradeoffTx build unlocked tradeoff contract tx.
func BuildCallTradeoffTx(s *Server, guid, contractUTXOID, seller, assetRequested string, spendWalletAmount, contractAmount, amountRequested uint64) (*buildTxResp, error) {
	// inputs:
	spendUTXOInput := SpendUTXOInput{
		Type:     "spend_utxo",
		OutputID: contractUTXOID,
	}

	spendWalletInput := SpendWalletInput{
		Type:    "spend_wallet",
		AssetID: BTMAssetID,
		Amount:  spendWalletAmount,
	}

	spendWalletUnlockTradeoffInput := SpendWalletInput{
		Type:    "spend_wallet",
		AssetID: assetRequested,
		Amount:  amountRequested,
	}

	// outputs:
	controlProgramOutput := ControlProgramOutput{
		Type:           "control_program",
		Amount:         amountRequested,
		AssetID:        assetRequested,
		ControlProgram: seller,
	}

	var inputs, outputs []interface{}
	inputs = append(inputs, spendUTXOInput, spendWalletInput, spendWalletUnlockTradeoffInput)
	outputs = append(outputs, controlProgramOutput)
	payload, err := json.Marshal(buildTxReq{
		GUID:          guid,
		Fee:           fee,
		Confirmations: confirmations,
		Inputs:        inputs,
		Outputs:       outputs,
	})
	if err != nil {
		return nil, err
	}

	fmt.Println("build unlocked contract tx:", string(payload))

	res := new(buildTxResp)
	if err := s.request(buildTransactionURL, payload, res); err != nil {
		return nil, err
	}

	return res, nil
}
