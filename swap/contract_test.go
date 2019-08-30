package swap

import (
	"fmt"
	"testing"
)

func TestContract(t *testing.T) {
	accountInfo := AccountInfo{
		AccountID: "10CJPO1HG0A02",
		Password:  "12345",
		TxFee:     uint64(40000000),
	}

	requestedContractValue := AssetAmount{
		Asset:  "ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff",
		Amount: uint64(1000000000),
	}

	LockedContractValue := AssetAmount{
		Asset:  "bae7e17bb8f5d0cfbfd87a92f3204da082d388d4c9b10e8dcd36b3d0a18ceb3a",
		Amount: uint64(20000000000),
	}

	contractArgs := ContractArgs{
		AssetAmount: requestedContractValue,
		Seller:      "00145dd7b82556226d563b6e7d573fe61d23bd461c1f",
		CancelKey:   "3e5d7d52d334964eef173021ef6a04dc0807ac8c41700fe718f5a80c2109f79e",
	}

	contractUTXOID, err := DeployContract(accountInfo, contractArgs, LockedContractValue)
	if err != nil {
		panic(err)
	}

	accountInfo = AccountInfo{
		AccountID: "10CKAD3000A02",
		Password:  "12345",
		Receiver:  "00140fdee108543d305308097019ceb5aec3da60ec66",
		TxFee:     uint64(40000000),
	}
	txID, err := CallContract(accountInfo, contractUTXOID, contractArgs, LockedContractValue)
	if err != nil {
		panic(err)
	}
	fmt.Println("--> txID:", txID)
}
