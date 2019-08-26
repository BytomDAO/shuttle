package main

import (
	"fmt"

	"github.com/btm-swap-tool/swap"
)

func main() {
	balances := swap.ListBalances("a1")
	fmt.Println("balances:", balances)

	accounts := swap.ListAccounts()
	fmt.Println("accounts:", accounts)

	addresses := swap.ListAddresses("a1")
	fmt.Println("addresses:", addresses)

	pubkeyInfo := swap.ListPubkeys("a1")
	fmt.Println(pubkeyInfo)

	assetRequested := "ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff"
	seller := "00145dd7b82556226d563b6e7d573fe61d23bd461c1f"
	cancelKey := "3e5d7d52d334964eef173021ef6a04dc0807ac8c41700fe718f5a80c2109f79e"
	amountRequested := uint64(1000000000)
	contractInfo := swap.CompileLockContract(assetRequested, seller, cancelKey, amountRequested)
	fmt.Println("contract info:", contractInfo)

	assetID := "bae7e17bb8f5d0cfbfd87a92f3204da082d388d4c9b10e8dcd36b3d0a18ceb3a"
	amount := uint64(20000000000)
	controlProgram := "203e5d7d52d334964eef173021ef6a04dc0807ac8c41700fe718f5a80c2109f79e1600145dd7b82556226d563b6e7d573fe61d23bd461c1f0400ca9a3b20ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff741a547a6413000000007b7b51547ac1631a000000547a547aae7cac00c0"
	tx := swap.BuildTransaction(assetID, controlProgram, amount)
	fmt.Println("tx:", tx)

	password := "12345"
	signedTransaction := swap.SignTransaction(password, string(tx))
	fmt.Println("signedTransaction:", signedTransaction)

	txID := swap.SubmitTransaction(signedTransaction)
	fmt.Println("txID:", txID)

	// txID := "e15674a9c694f56a0d172152df8dc28bbbe89f9828feb8398d3ebfb4e2f104ae"
	contractUTXOID, err := swap.GetContractUTXOID(txID, controlProgram)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("contractUTXOID:", contractUTXOID)

	// tx = buildUnlockContractTransaction(contractUTXOID, seller)
	// fmt.Println(tx)
}
