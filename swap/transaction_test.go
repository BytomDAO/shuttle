package swap

import (
	"encoding/json"
	"fmt"
	"testing"
)

var server = &Server{
	IP:   "127.0.0.1",
	Port: "3000",
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
	guid := "dbcb0423-ce18-47f0-8f03-dbc1260a3d76" // acount a1
	fee := uint64(40000000)
	confirmations := uint64(1)
	outputID := "e5b5a81e560536aae67ad21098109be17c56027f0b7dce9ab0be92bd38858c54"
	lockedAsset := "bae7e17bb8f5d0cfbfd87a92f3204da082d388d4c9b10e8dcd36b3d0a18ceb3a"
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

func TestSubmitPayment(t *testing.T) {
	guid := "dbcb0423-ce18-47f0-8f03-dbc1260a3d76" // acount a1
	rawTx := "07010002015f015d5af8651f89ce87c4ab542ce8db3d945cee238d6a0b90b43a7bbb34e5afa54c2effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff80d7cc5a0001160014b655549012a6b2ec01f6207638f65fcb2d6b7ef9630240506c41bedec62e2751f4fc17103b2c55fd52c929cc2c591f387f5b8c759e06574a13415a4b11dfb050ac512a39292420f6c121047a11056b789ce5c26e0d820b20be5144b1d13a0afa7d79d439551ade6419ca5af8cdbe7823378d7d9cabad53ae015f015d95a8062e70ef687c08b25ba497c564b8a3f05c49e9ab88c66d090784f1b91efeffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff80c9fe3d0001160014538138c216e836097ffcae2eb29cfa6eb9ebb88d63024066133e67617b54acc96c8ea8dcec40cb29a71a5f21e687f06c49c1ce1b67a2f2bfe0aea9d712ee496b58d966e177729da08541bf6fe2923144d9ad7c70d1aa0320e4a9372b3513b8462e56a95fb6a17a53d19c40f1b206f357a391cc7386403f0602013cffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff80e1eb170116001458ff9e04104430313e2a4a624be748453c956d9200013cffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff808bd66d0116001470f1ffbfd3c90a2e2d93429d1438e4f4b7025d3600"
	memo := ""

	txID, err := submitPayment(server, guid, rawTx, memo)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("submit tx result:", txID)
}

func TestSignMsg(t *testing.T) {
	signData := "12abef"
	xprv := "0a6c7936304a753592b2c70af998ab35ab39200f2bcc2655cfffef505412f8ecadf574f93a6a0a2cc6573bfa96c6deabd2ac06beb0bd59a4a77513d5b6e51319"
	sig, err := signMsg(signData, xprv)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("sig:", sig)
}
