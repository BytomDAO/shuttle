package swap

import (
	"fmt"
	"testing"
)

func TestGetUTXOID(t *testing.T) {
	server := &Server{
		IP:   "52.82.73.202",
		Port: "3060",
	}

	txID := "0d2b40feb0e64e910194ed19eac9627683064b848c196da674bef3a94dc3eba8"
	controlProgram := "001418b791936982ba3cc33112284aa65f575736d913"
	utxoID, err := getUTXOID(server, txID, controlProgram)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("utxoID:", utxoID)
}
