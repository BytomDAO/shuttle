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

	accountIDLocked := "10CJPO1HG0A02"                                                   // accountIDLocked represents account which create locked contract
	accountPasswordLocked := "12345"                                                     // accountPasswordLocked represents account password which create locked contract
	assetRequested := "ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff" // assetRequested represents asset ID which can unlock contract
	amountRequested := uint64(1000000000)                                                // amountRequested represents asset amount which can unlock contract
	seller := "00145dd7b82556226d563b6e7d573fe61d23bd461c1f"                             // control program which want to receive assetRequested
	cancelKey := "3e5d7d52d334964eef173021ef6a04dc0807ac8c41700fe718f5a80c2109f79e"      // cancelKey can cancel swap contract
	assetIDLocked := "bae7e17bb8f5d0cfbfd87a92f3204da082d388d4c9b10e8dcd36b3d0a18ceb3a"  // assetIDLocked represents locked asset ID
	amountLocked := uint64(20000000000)                                                  // amountLocked represents locked asset amount
	txFee := uint64(50000000)                                                            // txFee represents transaction fee

	accountIDUnlocked := "10CKAD3000A02"                                 // accountIDUnlocked represents account ID which create unlocked contract
	buyerContolProgram := "00140fdee108543d305308097019ceb5aec3da60ec66" // buyerContolProgram represents buyer control program
	// accountPasswordUnlocked := "12345"                                   // accountPasswordUnlocked represents account password which create locked contract

	// compile locked contract
	contractInfo := swap.CompileLockContract(assetRequested, seller, cancelKey, amountRequested)
	fmt.Println("--> contract info:", contractInfo)

	// build locked contract
	txLocked := swap.BuildLockTransaction(accountIDLocked, assetIDLocked, contractInfo.Program, amountLocked, txFee)
	fmt.Println("--> txLocked:", string(txLocked))

	// sign locked contract transaction
	signedTransaction := swap.SignTransaction(accountPasswordLocked, string(txLocked))
	fmt.Println("--> signedTransaction:", signedTransaction)

	// submit signed transaction
	txID := swap.SubmitTransaction(signedTransaction)
	fmt.Println("--> txID:", txID)

	// get contract output ID
	contractUTXOID, err := swap.GetContractUTXOID(txID, contractInfo.Program)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("--> contractUTXOID:", contractUTXOID)

	//
	// contractUTXOID = "38c0b998f8437aebe2be45e3c1893551ce303ff940fd2fcbb696f003233bdf5b"
	txUnlocked := swap.BuildUnlockContractTransaction(accountIDUnlocked, contractUTXOID, seller, assetIDLocked, assetRequested, buyerContolProgram, amountRequested, amountLocked, txFee)
	fmt.Println("--> txUnlocked:", string(txUnlocked))
}
