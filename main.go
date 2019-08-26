package main

import (
	"fmt"

	"github.com/btm-swap-tool/swap"
)

func main() {
	// balances := swap.ListBalances("a1")
	// fmt.Println("balances:", balances)

	// accounts := swap.ListAccounts()
	// fmt.Println("accounts:", accounts)

	// addresses := swap.ListAddresses("a1")
	// fmt.Println("addresses:", addresses)

	// pubkeyInfo := swap.ListPubkeys("a1")
	// fmt.Println(pubkeyInfo)

	accountIDSent := "10CJPO1HG0A02" // account a1 build swap contract
	assetRequested := "ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff"
	amountRequested := uint64(1000000000)
	seller := "00145dd7b82556226d563b6e7d573fe61d23bd461c1f"                            // control program which want to receive assetRequested
	cancelKey := "3e5d7d52d334964eef173021ef6a04dc0807ac8c41700fe718f5a80c2109f79e"     // cancelKey can cancel swap contract
	assetIDLocked := "bae7e17bb8f5d0cfbfd87a92f3204da082d388d4c9b10e8dcd36b3d0a18ceb3a" // assetIDLocked represents locked asset ID
	amountLocked := uint64(20000000000)                                                 // amountLocked represents locked asset amount
	txFee := uint64(10000000)                                                           // txFee represents transaction fee

	contractInfo := swap.CompileLockContract(assetRequested, seller, cancelKey, amountRequested)
	fmt.Println("contract info:", contractInfo)

	tx := swap.BuildLockTransaction(accountIDSent, assetIDLocked, contractInfo.Program, amountLocked, txFee)
	fmt.Println("tx:", string(tx))

	password := "12345"
	signedTransaction := swap.SignTransaction(password, string(tx))
	fmt.Println("signedTransaction:", signedTransaction)

	txID := swap.SubmitTransaction(signedTransaction)
	fmt.Println("txID:", txID)

	contractUTXOID, err := swap.GetContractUTXOID(txID, contractInfo.Program)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("contractUTXOID:", contractUTXOID)

	// tx = buildUnlockContractTransaction(contractUTXOID, seller)
	// fmt.Println(tx)
}
