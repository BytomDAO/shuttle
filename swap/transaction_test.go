package swap

import (
	"encoding/json"
	"fmt"
	"testing"
)

var server = &Server{
	IP:   "52.82.73.202",
	Port: "3060",
}

func TestGetUTXOID(t *testing.T) {
	txID := "0d2b40feb0e64e910194ed19eac9627683064b848c196da674bef3a94dc3eba8"
	controlProgram := "001418b791936982ba3cc33112284aa65f575736d913"
	utxoID, err := getUTXOID(server, txID, controlProgram)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("utxoID:", utxoID)
}

func TestBuildTx(t *testing.T) {
	guid := "90294d7d-0db1-4f6f-8485-990af6197e88" // acount a1
	fee := uint64(40000000)
	confirmations := uint64(1)
	outputID := "e5b5a81e560536aae67ad21098109be17c56027f0b7dce9ab0be92bd38858c54"
	lockedAsset := "31f326e629719627b3b07a428e28572dfc62bedf8d9d3ee4a44911b93e5a128b"
	controlAddress := "sm1qrzmerymfs2arese3zg5y4fjl2atndkgn3j8dea"
	contractProgram := "001418b791936982ba3cc33112284aa65f575736d913" // redeem to a1
	lockedAmount := uint64(1)

	buildTxResp, err := buildTx(server, guid, outputID, lockedAsset, controlAddress, contractProgram, fee, confirmations, lockedAmount)
	if err != nil {
		fmt.Println(err)
	}
	res, err := json.Marshal(buildTxResp)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(res))
}
