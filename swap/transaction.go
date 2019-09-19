package swap

import (
	"encoding/json"
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
