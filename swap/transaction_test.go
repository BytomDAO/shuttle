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
	outputID := "fa6e8ae89b3acdcfe8d8256c9adce856d87a658c0fe9c711136eca190b66c763"
	lockedAsset := "bae7e17bb8f5d0cfbfd87a92f3204da082d388d4c9b10e8dcd36b3d0a18ceb3a"
	lockedAmount := uint64(100)
	contractProgram := "20eec15ce68d46569f92ecebd7769101b22e34109892cc7ddfd54dc772f850c2431600145b0a81adc5c2d68a9967082a09c96e82d62aa058016420ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff741a547a6413000000007b7b51547ac1631a000000547a547aae7cac00c0"

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
	contractUTXOID := "dd5ebcbd1c8a9feaa82aad3b6d9b4c28784c4bd1d94acacce6156b47269dc429"
	contractAsset := "bae7e17bb8f5d0cfbfd87a92f3204da082d388d4c9b10e8dcd36b3d0a18ceb3a"
	spendWalletAmount := fee
	receiver := "sm1qe5gne93c8wx75ualxkju5yyec20j54ynjxd8zj" // account a4
	contractAmount := uint64(100)

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

func TestBuildCallTradeoffTx(t *testing.T) {
	guid := "e18b91ba-91a5-4837-9d41-ce2b76cea81c" // acount a1
	fee := uint64(40000000)
	confirmations := uint64(1)
	contractUTXOID := "dd5ebcbd1c8a9feaa82aad3b6d9b4c28784c4bd1d94acacce6156b47269dc429"
	assetRequested := "ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff"
	amountRequested := uint64(100)
	contractAsset := "bae7e17bb8f5d0cfbfd87a92f3204da082d388d4c9b10e8dcd36b3d0a18ceb3a"
	spendWalletAmount := fee
	seller := "00145b0a81adc5c2d68a9967082a09c96e82d62aa058" // seller program
	contractAmount := uint64(100)

	buildTxResp, err := buildCallTradeoffTx(server, guid, contractUTXOID, contractAsset, seller, assetRequested, fee, spendWalletAmount, confirmations, contractAmount, amountRequested)
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
	rawTx := "07010002016201603b2c6b69759cd0a2245b2f1a5681cf782e485f077a86ed13f82cae677a671d66bae7e17bb8f5d0cfbfd87a92f3204da082d388d4c9b10e8dcd36b3d0a18ceb3aaeb7a5f4a9e91502011600145b0a81adc5c2d68a9967082a09c96e82d62aa05801000161015fd284f4407b0b0a3ec00ec97b6145103c9943679db65001cbc7fdb61065c173ceffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8099c4d5990100011600145b0a81adc5c2d68a9967082a09c96e82d62aa058220120eec15ce68d46569f92ecebd7769101b22e34109892cc7ddfd54dc772f850c2430301b001bae7e17bb8f5d0cfbfd87a92f3204da082d388d4c9b10e8dcd36b3d0a18ceb3a904e018b01202cf24dba5fb0a30e26e83b2ac5b9e29e1b161e5c1fa7425e73043362938b98240164206ea28f3f1389efd6a731de070fb38ab69dc93dae6c73b6524bac901b662f601d20eec15ce68d46569f92ecebd7769101b22e34109892cc7ddfd54dc772f850c2437422547a6418000000557aa8547a88537a7bae7cac63220000007bcd9f69537a7cae7cac00c000013fbae7e17bb8f5d0cfbfd87a92f3204da082d388d4c9b10e8dcd36b3d0a18ceb3a9ee9a4f4a9e915011600145b0a81adc5c2d68a9967082a09c96e82d62aa05800013effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff80e5bac29901011600145b0a81adc5c2d68a9967082a09c96e82d62aa05800"
	memo := ""
	spendUTXOSig := "bdc314665cec20027fc25fab5b91c083e6ee614ccc4d50ef6da25884e398eeaf176fbbd0319dc1c3dca23e015040559ea313e992bd31dc30ea5f3d083957550e"
	spendUTXOPublicKey := "eec15ce68d46569f92ecebd7769101b22e34109892cc7ddfd54dc772f850c243"
	spendWalletSig := "e2b94f6604baceb9e8a6da7f87be542b9f24384010aed97c11cc613db43efd38b274a01f391d59996d1ea2fa7e2c7bcea4a7472910f29e4a5c223f080482a60b"

	spendUTXOSignatures := append([]string{}, spendUTXOSig, spendUTXOPublicKey)
	txID, err := submitPayment(server, guid, rawTx, memo, spendWalletSig, spendUTXOSignatures)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("submit tx result:", txID)
}

func TestSubmitUnlockedHTLCPayment(t *testing.T) {
	guid := "e18b91ba-91a5-4837-9d41-ce2b76cea81c" // acount a1
	rawTx := "0701000201d30101d0013b2c6b69759cd0a2245b2f1a5681cf782e485f077a86ed13f82cae677a671d66bae7e17bb8f5d0cfbfd87a92f3204da082d388d4c9b10e8dcd36b3d0a18ceb3a0100018b01202cf24dba5fb0a30e26e83b2ac5b9e29e1b161e5c1fa7425e73043362938b98240164206ea28f3f1389efd6a731de070fb38ab69dc93dae6c73b6524bac901b662f601d20eec15ce68d46569f92ecebd7769101b22e34109892cc7ddfd54dc772f850c2437422547a6418000000557aa8547a88537a7bae7cac63220000007bcd9f69537a7cae7cac00c001000161015fabf3111e8449df088eda1072c1bd4322157b62f588d8817d36f548f444092591ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8099c4d5990100011600145b0a81adc5c2d68a9967082a09c96e82d62aa058220120eec15ce68d46569f92ecebd7769101b22e34109892cc7ddfd54dc772f850c243020139bae7e17bb8f5d0cfbfd87a92f3204da082d388d4c9b10e8dcd36b3d0a18ceb3a0101160014cd113c96383b8dea73bf35a5ca1099c29f2a549300013effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff80e5bac29901011600145b0a81adc5c2d68a9967082a09c96e82d62aa05800"
	memo := ""
	preimage := "68656c6c6f" // hex("hello")
	spendUTXOSig := "749ebba8f0e59d9815b50bc9440667d7d5ff09baffa48c057b7d17e547e8f88e571cc690a791194c610ac5e4c91494031ac1142ffd7a13f48b548d6ab3a89a06"
	spendWalletSig := "c7bf8f1f18ac9bf94c051c80d3d6f9ce52e2edda3f4d36c37d795d747d649a6776d99d28e88eba7f1a1e521b25216f17b0f032d78b8c0c1c8a53e2c42d6e220c"

	spendUTXOSignatures := append([]string{}, preimage, spendUTXOSig, "")
	txID, err := submitPayment(server, guid, rawTx, memo, spendWalletSig, spendUTXOSignatures)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("submit tx result:", txID)
}

func TestSubmitCancelPayment(t *testing.T) {
	guid := "e18b91ba-91a5-4837-9d41-ce2b76cea81c" // acount a1
	rawTx := "0701000201d40101d1011edb4b73e76241f96de1bebf96dfa2c65e15a065cec78a23aa2fe1e1f3478a4ebae7e17bb8f5d0cfbfd87a92f3204da082d388d4c9b10e8dcd36b3d0a18ceb3a904e00018b01202cf24dba5fb0a30e26e83b2ac5b9e29e1b161e5c1fa7425e73043362938b98240164206ea28f3f1389efd6a731de070fb38ab69dc93dae6c73b6524bac901b662f601d20eec15ce68d46569f92ecebd7769101b22e34109892cc7ddfd54dc772f850c2437422547a6418000000557aa8547a88537a7bae7cac63220000007bcd9f69537a7cae7cac00c001000161015f06ba9c21d8cb432dc89282815d1254b29f1182cb164d7ad64f6a8ace8f328297ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8099c4d5990100011600145b0a81adc5c2d68a9967082a09c96e82d62aa058220120eec15ce68d46569f92ecebd7769101b22e34109892cc7ddfd54dc772f850c243030139bae7e17bb8f5d0cfbfd87a92f3204da082d388d4c9b10e8dcd36b3d0a18ceb3a6301160014cd113c96383b8dea73bf35a5ca1099c29f2a549300013abae7e17bb8f5d0cfbfd87a92f3204da082d388d4c9b10e8dcd36b3d0a18ceb3aad4d011600145b0a81adc5c2d68a9967082a09c96e82d62aa05800013effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff80e5bac29901011600145b0a81adc5c2d68a9967082a09c96e82d62aa05800"
	memo := ""
	spendUTXOSig := "84fa1ef1b5ec3ae07c616999323671d03392d994f317545c3c130d2a2eb2fe1110fb92f00e47e61f76f0b168c871eb4d66c77169795cd0fd05517e29f0f79505"
	spendWalletSig := "54d73d6d4b5ee8c5d675e347767e50784b8eafe3b80dae3bf054e70a4fb29bde964d37e514ff67af904b687feeaeaea21c48b4169fdc88d133d9aaf0c2c8070b"

	spendUTXOSignatures := append([]string{}, spendUTXOSig, "01")
	txID, err := submitPayment(server, guid, rawTx, memo, spendWalletSig, spendUTXOSignatures)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("submit tx result:", txID)
}

func TestSubmitTradeoffPayment(t *testing.T) {
	guid := "e18b91ba-91a5-4837-9d41-ce2b76cea81c" // acount a1
	rawTx := "07010002016201601edb4b73e76241f96de1bebf96dfa2c65e15a065cec78a23aa2fe1e1f3478a4ebae7e17bb8f5d0cfbfd87a92f3204da082d388d4c9b10e8dcd36b3d0a18ceb3a9ee9a4f4a9e91501011600145b0a81adc5c2d68a9967082a09c96e82d62aa05801000161015f3ec2cd05e9c3d0e3834386d3fc1d041b82892600996e485c7539be53e9168faaffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8099c4d5990100011600145b0a81adc5c2d68a9967082a09c96e82d62aa058220120eec15ce68d46569f92ecebd7769101b22e34109892cc7ddfd54dc772f850c24303019c01bae7e17bb8f5d0cfbfd87a92f3204da082d388d4c9b10e8dcd36b3d0a18ceb3a64017920eec15ce68d46569f92ecebd7769101b22e34109892cc7ddfd54dc772f850c2431600145b0a81adc5c2d68a9967082a09c96e82d62aa058016420ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff741a547a6413000000007b7b51547ac1631a000000547a547aae7cac00c000013fbae7e17bb8f5d0cfbfd87a92f3204da082d388d4c9b10e8dcd36b3d0a18ceb3abae8a4f4a9e915011600145b0a81adc5c2d68a9967082a09c96e82d62aa05800013effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff80e5bac29901011600145b0a81adc5c2d68a9967082a09c96e82d62aa05800"
	memo := ""
	spendUTXOSig := "c88cdd08f099c7a4a3aa4075e6439285e93c782dfe03f3160f79a171feb2c2f4a4acf716aba3cd4b0a33308f63c75ff23965912ef2cdf5599574eaa426c70d05"
	spendUTXOPublicKey := "eec15ce68d46569f92ecebd7769101b22e34109892cc7ddfd54dc772f850c243"
	spendWalletSig := "633c41fc2de572dcb4284e2d66e6412c23696feba4e8f381abab3ce8b3508eace660a1b3d80b8421a49ca49edadf2a262b054379cb258ffb9d7f37c2b87d1f09"

	spendUTXOSignatures := append([]string{}, spendUTXOSig, spendUTXOPublicKey)
	txID, err := submitPayment(server, guid, rawTx, memo, spendWalletSig, spendUTXOSignatures)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("submit tx result:", txID)
}

func TestSubmitCallTradeoffPayment(t *testing.T) {
	guid := "e18b91ba-91a5-4837-9d41-ce2b76cea81c" // acount a1
	rawTx := "0701000201c00101bd012b30c07d12cc4e20268976694f4213fd0aa0d2406bab92b9770185d62415dcf9bae7e17bb8f5d0cfbfd87a92f3204da082d388d4c9b10e8dcd36b3d0a18ceb3a6400017920eec15ce68d46569f92ecebd7769101b22e34109892cc7ddfd54dc772f850c2431600145b0a81adc5c2d68a9967082a09c96e82d62aa058016420ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff741a547a6413000000007b7b51547ac1631a000000547a547aae7cac00c001000161015fc9a1ab37e062e33d6dbc1aa1a829ea031abbfc3a5892a3fd8f56c9b06a82fa07ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8099c4d5990100011600145b0a81adc5c2d68a9967082a09c96e82d62aa058220120eec15ce68d46569f92ecebd7769101b22e34109892cc7ddfd54dc772f850c243020139bae7e17bb8f5d0cfbfd87a92f3204da082d388d4c9b10e8dcd36b3d0a18ceb3a6401160014cd113c96383b8dea73bf35a5ca1099c29f2a549300013effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff80e5bac29901011600145b0a81adc5c2d68a9967082a09c96e82d62aa05800"
	memo := ""
	spendWalletSig := "71a2577526e82b92c1bda051df12057405f8df90937fca7e0837fd990b1d8586d74d22e80f499a1e25533557b6ad12327a8ebb4814183f9acbd8dade73352a0b"

	spendUTXOSignatures := append([]string{}, "")
	txID, err := submitPayment(server, guid, rawTx, memo, spendWalletSig, spendUTXOSignatures)
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
