package swap

import (
	"encoding/json"
	"fmt"
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
func buildTx(s *Server, guid, outputID, lockedAsset, controlAddress, contractProgram string, fee, confirmations, lockedAmount uint64) (*buildTxResp, error) {
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
	// controlAddressOutput := ControlAddressOutput{
	// 	Type:    "control_address",
	// 	Amount:  lockedAmount,
	// 	AssetID: lockedAsset,
	// 	Address: controlAddress,
	// }
	controlProgramOutput := ControlProgramOutput{
		Type:           "control_program",
		Amount:         lockedAmount,
		AssetID:        lockedAsset,
		ControlProgram: contractProgram,
	}

	var inputs, outputs []interface{}
	inputs = append(inputs, spendUTXOInput, spendWalletInput)
	// outputs = append(outputs, controlAddressOutput, controlProgramOutput)
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
	if err := s.request(getTransactionURL, payload, res); err != nil {
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
func submitPayment(s *Server, guid, rawTx, memo string) (string, error) {
	payload, err := json.Marshal(submitPaymentReq{
		GUID:       guid,
		RawTx:      rawTx,
		Signatures: [][]string{},
		Memo:       memo,
	})
	if err != nil {
		return "", err
	}

	fmt.Println("payload:", string(payload))

	res := new(submitPaymentResp)
	if err := s.request(submitTransactionURL, payload, res); err != nil {
		return "", err
	}

	return res.TxID, nil
}

// func signMessage(s *Server, signData, xprv string) (string, error) {

// }
