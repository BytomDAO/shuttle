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
	outputID := "fa898bb1daa5c7bf799809c01823edaafec0af920f7ab6de26dd17aa7e6c29a4"
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
	contractUTXOID := "11317c0bf0c39b7b1d3082a292e6ecbd010d59cdf59007d44ea7bc7b9c36c337"
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
	rawTx := "0701000201620160b631510ab58859eb6834068a3ec9dc6104efe8c13fba707908279be029b8a29abae7e17bb8f5d0cfbfd87a92f3204da082d388d4c9b10e8dcd36b3d0a18ceb3aafb7a5f4a9e91501011600145b0a81adc5c2d68a9967082a09c96e82d62aa05801000161015f39276de0fc73758efc6753169284862e2ca67ab6218ae6500ef23f9ca87f8f2cffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8099c4d5990100011600145b0a81adc5c2d68a9967082a09c96e82d62aa058220120eec15ce68d46569f92ecebd7769101b22e34109892cc7ddfd54dc772f850c2430301af01bae7e17bb8f5d0cfbfd87a92f3204da082d388d4c9b10e8dcd36b3d0a18ceb3a01018b01202cf24dba5fb0a30e26e83b2ac5b9e29e1b161e5c1fa7425e73043362938b98240164206ea28f3f1389efd6a731de070fb38ab69dc93dae6c73b6524bac901b662f601d20eec15ce68d46569f92ecebd7769101b22e34109892cc7ddfd54dc772f850c2437422547a6418000000557aa8547a88537a7bae7cac63220000007bcd9f69537a7cae7cac00c000013effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff80e5bac29901011600145b0a81adc5c2d68a9967082a09c96e82d62aa05800013fbae7e17bb8f5d0cfbfd87a92f3204da082d388d4c9b10e8dcd36b3d0a18ceb3aaeb7a5f4a9e915011600145b0a81adc5c2d68a9967082a09c96e82d62aa05800"
	memo := ""
	spendUTXOSig := "749ebba8f0e59d9815b50bc9440667d7d5ff09baffa48c057b7d17e547e8f88e571cc690a791194c610ac5e4c91494031ac1142ffd7a13f48b548d6ab3a89a06"
	spendUTXOPublicKey := "6ea28f3f1389efd6a731de070fb38ab69dc93dae6c73b6524bac901b662f601d"
	spendWalletSig := "c7bf8f1f18ac9bf94c051c80d3d6f9ce52e2edda3f4d36c37d795d747d649a6776d99d28e88eba7f1a1e521b25216f17b0f032d78b8c0c1c8a53e2c42d6e220c"

	spendUTXOSignatures := append([]string{}, spendUTXOSig, spendUTXOPublicKey)
	txID, err := submitPayment(server, guid, rawTx, memo, spendWalletSig, spendUTXOSignatures)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("submit tx result:", txID)
}

func TestSubmitUnlockedPayment(t *testing.T) {
	guid := "e18b91ba-91a5-4837-9d41-ce2b76cea81c" // acount a1
	rawTx := "0701000201d30101d0013b2c6b69759cd0a2245b2f1a5681cf782e485f077a86ed13f82cae677a671d66bae7e17bb8f5d0cfbfd87a92f3204da082d388d4c9b10e8dcd36b3d0a18ceb3a0100018b01202cf24dba5fb0a30e26e83b2ac5b9e29e1b161e5c1fa7425e73043362938b98240164206ea28f3f1389efd6a731de070fb38ab69dc93dae6c73b6524bac901b662f601d20eec15ce68d46569f92ecebd7769101b22e34109892cc7ddfd54dc772f850c2437422547a6418000000557aa8547a88537a7bae7cac63220000007bcd9f69537a7cae7cac00c001000161015fabf3111e8449df088eda1072c1bd4322157b62f588d8817d36f548f444092591ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8099c4d5990100011600145b0a81adc5c2d68a9967082a09c96e82d62aa058220120eec15ce68d46569f92ecebd7769101b22e34109892cc7ddfd54dc772f850c243020139bae7e17bb8f5d0cfbfd87a92f3204da082d388d4c9b10e8dcd36b3d0a18ceb3a0101160014cd113c96383b8dea73bf35a5ca1099c29f2a549300013effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff80e5bac29901011600145b0a81adc5c2d68a9967082a09c96e82d62aa05800"
	memo := ""
	preimage := "68656c6c6f" // hex("hello")
	spendUTXOSig := "749ebba8f0e59d9815b50bc9440667d7d5ff09baffa48c057b7d17e547e8f88e571cc690a791194c610ac5e4c91494031ac1142ffd7a13f48b548d6ab3a89a06"
	spendWalletSig := "c7bf8f1f18ac9bf94c051c80d3d6f9ce52e2edda3f4d36c37d795d747d649a6776d99d28e88eba7f1a1e521b25216f17b0f032d78b8c0c1c8a53e2c42d6e220c"

	spendUTXOSignatures := append([]string{}, preimage, spendUTXOSig, "")
	txID, err := submitPayment(server, guid, rawTx, memo, spendWalletSig, spendUTXOSignatures)
	// txID, err := submitUnlockedPayment(server, guid, rawTx, memo, spendWalletSig, spendUTXOSignatures)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("submit tx result:", txID)
}

func TestSignMsg(t *testing.T) {
	signData := "6b2c56ad857c98b602ea2211641c94b5e7979214b418777b359d9cc3856f7f12"
	xprv := "682d87647c76edafb0c0bdb8b9a87e84f79627c86a4d7620c89a9ef7551ecf47013095e747f609c86703ee7c0281b2182dfaca66d60ea58814d7929e6b6968a5"
	sig, err := signMsg(signData, xprv)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("sig:", sig)
}
