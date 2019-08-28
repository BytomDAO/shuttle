package main

import (
	"errors"
	"flag"
	"fmt"

	"github.com/btm-swap-tool/swap"
)

var (
	accountIDLocked       = "10CJPO1HG0A02"                                                    // accountIDLocked represents account which create locked contract
	accountPasswordLocked = "12345"                                                            // accountPasswordLocked represents account password which create locked contract
	assetRequested        = "ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff" // assetRequested represents asset ID which can unlock contract
	amountRequested       = uint64(1000000000)                                                 // amountRequested represents asset amount which can unlock contract
	seller                = "00145dd7b82556226d563b6e7d573fe61d23bd461c1f"                     // control program which want to receive assetRequested
	cancelKey             = "3e5d7d52d334964eef173021ef6a04dc0807ac8c41700fe718f5a80c2109f79e" // cancelKey can cancel swap contract
	assetIDLocked         = "bae7e17bb8f5d0cfbfd87a92f3204da082d388d4c9b10e8dcd36b3d0a18ceb3a" // assetIDLocked represents locked asset ID
	amountLocked          = uint64(20000000000)                                                // amountLocked represents locked asset amount
	txFee                 = uint64(50000000)                                                   // txFee represents transaction fee

	accountIDUnlocked       = "10CKAD3000A02"                                // accountIDUnlocked represents account ID which create unlocked contract
	buyerContolProgram      = "00140fdee108543d305308097019ceb5aec3da60ec66" // buyerContolProgram represents buyer control program
	accountPasswordUnlocked = "12345"                                        // accountPasswordUnlocked represents account password which create locked contract
)

var (
	errInvalidAssetRequested = errors.New("assetRequested is invalid")
	errInvalidSeller         = errors.New("seller is invalid")
	errInvalidCancelKey      = errors.New("cancelKey is invalid")
	errInvalidAccountID      = errors.New("accountID is invalid")
	errInvalidAssetIDLocked  = errors.New("assetIDLocked is invalid")
	errInvalidTxFee          = errors.New("txFee is invalid")
	errInvalidAmountLocked   = errors.New("amountLocked is invalid")
)

func main() {
	var accountID string
	flag.StringVar(&accountID, "accountID", "", "account ID, which is an identifier of account.")

	var password string
	flag.StringVar(&password, "password", "", "account password, which can decrypt account key.")

	var assetRequested string
	flag.StringVar(&assetRequested, "assetRequested", "", "asset ID, which is asset ID sent to locked contract by buyer.")

	var amountRequested uint64
	flag.Uint64Var(&amountRequested, "amountRequested", uint64(0), "asset amount, which is asset amount sent to locked contract by buyer.")

	var seller string
	flag.StringVar(&seller, "seller", "", "seller is control program of contract creator.")

	var cancelKey string
	flag.StringVar(&cancelKey, "cancelKey", "", "public key, which can cancel contract.")

	var assetIDLocked string
	flag.StringVar(&assetIDLocked, "assetIDLocked", "", "asset ID, which is asset ID locked in contract by creator.")

	var amountLocked uint64
	flag.Uint64Var(&amountLocked, "amountLocked", uint64(0), "asset amount, which is asset amount locked in contract by creator.")

	var txFee uint64
	flag.Uint64Var(&txFee, "txFee", uint64(50000000), "tx fee, which is transaction fee.")

	flag.Parse()

	// fmt.Println("arg is :", flag.Arg(0))

	if flag.Arg(0) == "deploy" {
		if err := checkDeployParameter(assetRequested, seller, cancelKey, accountIDLocked, assetIDLocked, txFee, amountLocked); err != nil {
			panic(err)
		}
		contractUTXOID := swap.DeployContract(assetRequested, seller, cancelKey, accountIDLocked, assetIDLocked, accountPasswordLocked, amountRequested, amountLocked, txFee)
		fmt.Println("--> contractUTXOID:", contractUTXOID)
	}

	// contractUTXOID := swap.DeployContract(assetRequested, seller, cancelKey, accountIDLocked, assetIDLocked, accountPasswordLocked, amountRequested, amountLocked, txFee)
	// txID := swap.CallContract(accountIDUnlocked, contractUTXOID, seller, assetIDLocked, assetRequested, buyerContolProgram, accountPasswordUnlocked, amountRequested, amountLocked, txFee)
	// fmt.Println("--> txID:", txID)
}

func checkDeployParameter(assetRequested, seller, cancelKey, accountIDLocked, assetIDLocked string, txFee, amountLocked uint64) error {
	if len(assetRequested) != 64 {
		return errInvalidAssetRequested
	}

	if len(seller) == 0 {
		return errInvalidSeller
	}

	if len(cancelKey) == 0 {
		return errInvalidCancelKey
	}

	if len(accountIDLocked) == 0 {
		return errInvalidAccountID
	}

	if len(assetIDLocked) != 64 {
		return errInvalidAssetIDLocked
	}

	if txFee == uint64(0) {
		return errInvalidTxFee
	}

	if amountLocked == uint64(0) {
		return errInvalidAmountLocked
	}

	return nil
}
