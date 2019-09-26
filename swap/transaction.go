package swap

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/bytom/crypto/ed25519/chainkd"
)

var (
	errXPrvLength = errors.New("XPrv length is invalid.")
)

type TxOutput struct {
	UTXOID      string `json:"utxo_id"`
	Script      string `json:"script"`
	Address     string `json:"address"`
	AssetID     string `json:"asset"`
	AssetAmount uint64 `json:"amount"`
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

// buildTx build tx.
func buildTx(s *Server, guid, outputID, lockedAsset, contractProgram string, fee, confirmations, lockedAmount uint64) (*buildTxResp, error) {
	// inputs:
	spendUTXOInput := SpendUTXOInput{
		Type:     "spend_utxo",
		OutputID: outputID,
	}
	spendWalletInput := SpendWalletInput{
		Type:    "spend_wallet",
		AssetID: BTMAssetID,
		Amount:  fee,
	}

	// outputs:
	controlProgramOutput := ControlProgramOutput{
		Type:           "control_program",
		Amount:         lockedAmount,
		AssetID:        lockedAsset,
		ControlProgram: contractProgram,
	}

	var inputs, outputs []interface{}
	inputs = append(inputs, spendUTXOInput, spendWalletInput)
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

	res := new(buildTxResp)
	if err := s.request(buildTransactionURL, payload, res); err != nil {
		return nil, err
	}

	return res, nil
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

// submitPayment submit raw transaction and return transaction ID.
func submitPayment(s *Server, guid, rawTx, memo, spendWalletSig string, spendUTXOSignatures []string) (string, error) {
	// spendUTXOSignatures := append([]string{}, spendUTXOSig, spendUTXOPublicKey)
	spendWalletSignatures := append([]string{}, spendWalletSig)
	sigs := append([][]string{}, spendUTXOSignatures, spendWalletSignatures)

	payload, err := json.Marshal(submitPaymentReq{
		GUID:       guid,
		RawTx:      rawTx,
		Signatures: sigs,
		Memo:       memo,
	})
	if err != nil {
		return "", err
	}

	res := new(submitPaymentResp)
	if err := s.request(submitTransactionURL, payload, res); err != nil {
		return "", err
	}

	return res.TxID, nil
}

// signMessage sign message, return sig.
func signMsg(signData, xprv string) (string, error) {
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

// buildUnlockedTx build unlocked contract tx.
func buildUnlockedTx(s *Server, guid, contractUTXOID, contractAsset, receiver string, fee, spendWalletAmount, confirmations, contractAmount uint64) (*buildTxResp, error) {
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

	// outputs:
	controlAddressOutput := ControlAddressOutput{
		Type:    "control_address",
		Amount:  contractAmount,
		AssetID: contractAsset,
		Address: receiver,
	}

	var inputs, outputs []interface{}
	inputs = append(inputs, spendUTXOInput, spendWalletInput)
	outputs = append(outputs, controlAddressOutput)
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

// submitUnlockedPayment submit raw transaction and return transaction ID.
func submitUnlockedPayment(s *Server, guid, rawTx, memo, spendWalletSig string, spendUTXOSignatures []string) (string, error) {
	// spendUTXOSignatures := append([]string{}, preimage, spendUTXOSig, "")
	spendWalletSignatures := append([]string{}, spendWalletSig)
	sigs := append([][]string{}, spendUTXOSignatures, spendWalletSignatures)

	payload, err := json.Marshal(submitPaymentReq{
		GUID:       guid,
		RawTx:      rawTx,
		Signatures: sigs,
		Memo:       memo,
	})
	if err != nil {
		return "", err
	}

	fmt.Println("submitUnlockedPayment:", string(payload))

	res := new(submitPaymentResp)
	if err := s.request(submitTransactionURL, payload, res); err != nil {
		return "", err
	}

	return res.TxID, nil
}
