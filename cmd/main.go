package main

import (
	"errors"
	"flag"
	"fmt"

	"github.com/btm-swap-tool/swap"
)

var (
	errInvalidAssetRequested     = errors.New("assetRequested is invalid")
	errInvalidSeller             = errors.New("seller is invalid")
	errInvalidCancelKey          = errors.New("cancelKey is invalid")
	errInvalidAccountID          = errors.New("accountID is invalid")
	errInvalidAssetIDLocked      = errors.New("assetIDLocked is invalid")
	errInvalidTxFee              = errors.New("txFee is invalid")
	errInvalidAmountLocked       = errors.New("amountLocked is invalid")
	errInvalidContractUTXOID     = errors.New("contractUTXOID is invalid")
	errInvalidBuyerContolProgram = errors.New("buyerContolProgram is invalid")
	errInvalidAmountRequested    = errors.New("amountRequested is invalid")
)

/*
deploy:
go run main.go -assetRequested=ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff -seller=00145dd7b82556226d563b6e7d573fe61d23bd461c1f -cancelKey=3e5d7d52d334964eef173021ef6a04dc0807ac8c41700fe718f5a80c2109f79e -accountID=10CJPO1HG0A02 -assetIDLocked=bae7e17bb8f5d0cfbfd87a92f3204da082d388d4c9b10e8dcd36b3d0a18ceb3a -password=12345 -amountRequested=1000000000 -amountLocked=20000000000 -txFee=50000000  deploy

call:
go run main.go -accountID=10CKAD3000A02 -contractUTXOID=9544dac2b66314b242d7f5ce1e1d6a4d8dda9e38f6fb5008e99356aa91efb565 -seller=00145dd7b82556226d563b6e7d573fe61d23bd461c1f -assetIDLocked=bae7e17bb8f5d0cfbfd87a92f3204da082d388d4c9b10e8dcd36b3d0a18ceb3a -assetRequested=ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff -buyerContolProgram=00140fdee108543d305308097019ceb5aec3da60ec66 -password=12345 -amountRequested=1000000000 -amountLocked=20000000000 -txFee=50000000 call
*/

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

	var contractUTXOID string
	flag.StringVar(&contractUTXOID, "contractUTXOID", "", "contract UTXO ID, which is locked contract output ID.")

	var buyerContolProgram string
	flag.StringVar(&buyerContolProgram, "buyerContolProgram", "", "which is control program of buyer.")

	flag.Parse()
	if flag.Arg(0) == "deploy" {
		if err := checkDeployParameter(assetRequested, seller, cancelKey, accountID, assetIDLocked, txFee, amountLocked); err != nil {
			fmt.Println(err)
			return
		}
		contractUTXOID = swap.DeployContract(assetRequested, seller, cancelKey, accountID, assetIDLocked, password, amountRequested, amountLocked, txFee)
		fmt.Println("--> contractUTXOID:", contractUTXOID)
	} else if flag.Arg(0) == "call" {
		if err := checkCallParameter(accountID, contractUTXOID, seller, assetIDLocked, assetRequested, buyerContolProgram, amountRequested, amountLocked, txFee); err != nil {
			fmt.Println(err)
			return
		}
		txID := swap.CallContract(accountID, contractUTXOID, seller, assetIDLocked, assetRequested, buyerContolProgram, password, amountRequested, amountLocked, txFee)
		fmt.Println("--> txID:", txID)
	} else {
		fmt.Println("You should specify deploy or call.")
		return
	}
}

func checkDeployParameter(assetRequested, seller, cancelKey, accountID, assetIDLocked string, txFee, amountLocked uint64) error {
	if len(assetRequested) != 64 {
		return errInvalidAssetRequested
	}

	if len(seller) == 0 {
		return errInvalidSeller
	}

	if len(cancelKey) == 0 {
		return errInvalidCancelKey
	}

	if len(accountID) == 0 {
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

func checkCallParameter(accountID, contractUTXOID, seller, assetIDLocked, assetRequested, buyerContolProgram string, amountRequested, amountLocked, txFee uint64) error {
	if len(accountID) == 0 {
		return errInvalidAccountID
	}

	if len(contractUTXOID) != 64 {
		return errInvalidContractUTXOID
	}

	if len(seller) == 0 {
		return errInvalidSeller
	}

	if len(assetIDLocked) != 64 {
		return errInvalidAssetIDLocked
	}

	if len(assetRequested) != 64 {
		return errInvalidAssetRequested
	}

	if len(buyerContolProgram) == 0 {
		return errInvalidBuyerContolProgram
	}

	if amountRequested == uint64(0) {
		return errInvalidAmountRequested
	}

	if amountLocked == uint64(0) {
		return errInvalidAmountLocked
	}

	if txFee == uint64(0) {
		return errInvalidTxFee
	}

	return nil
}
