package swap

import (
	"fmt"
	"testing"
)

func TestDeployHTLCContract(t *testing.T) {
	account := HTLCAccount{
		AccountID: "10CJPO1HG0A02",
		Password:  "12345",
		TxFee:     uint64(100000000),
	}
	contractArgs := HTLCContractArgs{
		SenderPublicKey:    "3e5d7d52d334964eef173021ef6a04dc0807ac8c41700fe718f5a80c2109f79e",
		RecipientPublicKey: "198787c8380ed1ba6fec1f81bb68c17c16432c4bc646effe0a5fae4f1b528f16",
		BlockHeight:        uint64(950),
		Hash:               "2cf24dba5fb0a30e26e83b2ac5b9e29e1b161e5c1fa7425e73043362938b9824",
	}
	contractValue := AssetAmount{
		Asset:  "bae7e17bb8f5d0cfbfd87a92f3204da082d388d4c9b10e8dcd36b3d0a18ceb3a",
		Amount: uint64(20000000000),
	}
	contractUTXOID, err := DeployHTLCContract(account, contractValue, contractArgs)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("contractUTXOID:", contractUTXOID)
}
