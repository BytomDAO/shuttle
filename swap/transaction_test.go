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

func TestSubmitPayment(t *testing.T) {
	guid := "90294d7d-0db1-4f6f-8485-990af6197e88" // acount a1
	rawTx := "070100010161015f1773bf774fc34758c4bf86cc9fd021d8d19021414b286317d8bcf3ede61de55affffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffa0afba9ca31700011600144e5b5b946da8618f4ed30f3178de4813a7a62516630240d5d284f92d342355723f0452ec0cd6e2156bcc979cc1e09eeaff7fadf6158505d9491212431e5504ab7f201232798e18b64cc4b73909d4220eb32b9ba1f3ed0920c9e310c4d7b9bd973e2836e30d1fd6cf083252dc8a2c35b4281e0d393300693102013effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffa0f0da9ba217011600145f50069d54981661527191a4b0d1e9c79a7d8e1300013cffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff808bd66d01160014dc3e2f3dfacd20a5de691c1055655d57043aeaa500"
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

8c5bcf0f5e8ed94c6db496a26c0d77355967e5f147167ba8ab37529a9c76686f
59967b2d7ac15693d56c166dbbd017b5d6cb3013a0b766ac9479d705ae820109