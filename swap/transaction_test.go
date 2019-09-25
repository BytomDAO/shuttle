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
	guid := "e18b91ba-91a5-4837-9d41-ce2b76cea81c" // acount a1
	fee := uint64(40000000)
	confirmations := uint64(1)
	outputID := "21f174b985c667fcbf79d07a7b2e58a91a37d13d28f354978acaa70c822e0b97"
	lockedAsset := "bae7e17bb8f5d0cfbfd87a92f3204da082d388d4c9b10e8dcd36b3d0a18ceb3a"
	lockedAmount := uint64(1)
	contractProgram := "202cf24dba5fb0a30e26e83b2ac5b9e29e1b161e5c1fa7425e73043362938b98240164206ea28f3f1389efd6a731de070fb38ab69dc93dae6c73b6524bac901b662f601d20eec15ce68d46569f92ecebd7769101b22e34109892cc7ddfd54dc772f850c2437422547a6418000000557aa8547a88537a7bae7cac63220000007bcd9f69537a7cae7cac00c0"

	buildTxResp, err := buildTx(server, guid, outputID, lockedAsset, contractProgram, fee, confirmations, lockedAmount)
	if err != nil {
		fmt.Println(err)
	}
	res, err := json.Marshal(buildTxResp)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(res))
}

func TestBuildUnlockedTx(t *testing.T) {
	guid := "e18b91ba-91a5-4837-9d41-ce2b76cea81c" // acount a1
	fee := uint64(40000000)
	confirmations := uint64(1)
	contractUTXOID := "aa6cecf0c8768c05182ce5389d5a70c9fb9ca6e5697dc098568ffcc735093235"
	contractAsset := "bae7e17bb8f5d0cfbfd87a92f3204da082d388d4c9b10e8dcd36b3d0a18ceb3a"
	spendWalletAmount := fee
	receiver := "sm1qe5gne93c8wx75ualxkju5yyec20j54ynjxd8zj" // account a4
	contractAmount := uint64(1)

	buildTxResp, err := buildUnlockedTx(server, guid, contractUTXOID, contractAsset, receiver, fee, spendWalletAmount, confirmations, contractAmount)
	if err != nil {
		fmt.Println(err)
	}
	res, err := json.Marshal(buildTxResp)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("build unlocked response:", string(res))
}

func TestSubmitPayment(t *testing.T) {
	guid := "e18b91ba-91a5-4837-9d41-ce2b76cea81c" // acount a1
	rawTx := "0701000201620160ab5d3e83a2055d7d02381106a1b4fd44af5c4fac2fbaa2dde5e40a6e5a49932fbae7e17bb8f5d0cfbfd87a92f3204da082d388d4c9b10e8dcd36b3d0a18ceb3ab0b7a5f4a9e91502011600145b0a81adc5c2d68a9967082a09c96e82d62aa05801000161015f8057da31069f7630b2985d4a328b6d3c2d353652dbd2a3f44d54ba6085524ddcffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8099c4d5990100011600145b0a81adc5c2d68a9967082a09c96e82d62aa058220120eec15ce68d46569f92ecebd7769101b22e34109892cc7ddfd54dc772f850c2430301af01bae7e17bb8f5d0cfbfd87a92f3204da082d388d4c9b10e8dcd36b3d0a18ceb3a01018b01202cf24dba5fb0a30e26e83b2ac5b9e29e1b161e5c1fa7425e73043362938b98240164206ea28f3f1389efd6a731de070fb38ab69dc93dae6c73b6524bac901b662f601d20eec15ce68d46569f92ecebd7769101b22e34109892cc7ddfd54dc772f850c2437422547a6418000000557aa8547a88537a7bae7cac63220000007bcd9f69537a7cae7cac00c000013fbae7e17bb8f5d0cfbfd87a92f3204da082d388d4c9b10e8dcd36b3d0a18ceb3aafb7a5f4a9e915011600145b0a81adc5c2d68a9967082a09c96e82d62aa05800013effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff80e5bac29901011600145b0a81adc5c2d68a9967082a09c96e82d62aa05800"
	memo := ""
	spendUTXOSig := "ca5058e664c716dd7f086affc2026080de97237bb2aa73963aae6328baf4142b866fb0288039578a1eb85f50b5aaaa7c4f0a4985a7f6bd3b3c03b3e1cb6dee0e"
	spendUTXOPublicKey := "eec15ce68d46569f92ecebd7769101b22e34109892cc7ddfd54dc772f850c243"
	spendWalletSig := "a95af3c493e147f216d433d75e015d137cf6bcf10ec33dc5c8a152c10c7a05bb65dd91815e93f4ecbaaa7a4a390bab71e70beb1a2257b11b052047a61a93de01"

	txID, err := submitPayment(server, guid, rawTx, memo, spendUTXOSig, spendUTXOPublicKey, spendWalletSig)
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
