package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/btm-swap-tool/swap"
)

func init() {
	// deploy contract arguments
	deployCmd.PersistentFlags().Uint64Var(&txFee, "txFee", 40000000, "contract transaction fee")

	deployCmd.PersistentFlags().StringVar(&assetRequested, "assetRequested", "", "tradeoff contract paramenter with requested assetID")
	deployCmd.PersistentFlags().Uint64Var(&amountRequested, "amountRequested", 0, "tradeoff contract paramenter with requested amount")
	deployCmd.PersistentFlags().StringVar(&seller, "seller", "", "tradeoff contract paramenter with seller control-program")
	deployCmd.PersistentFlags().StringVar(&cancelKey, "cancelKey", "", "tradeoff contract paramenter with seller pubkey for cancelling the contract")

	deployCmd.PersistentFlags().StringVar(&assetLocked, "assetLocked", "", "tradeoff contract locked value with assetID")
	deployCmd.PersistentFlags().Uint64Var(&amountLocked, "amountLocked", 0, "tradeoff contract locked value with amount")

	// deploy HTLC contract arguments
	deployHTLCCmd.PersistentFlags().Uint64Var(&txFee, "txFee", 40000000, "contract transaction fee")

	deployHTLCCmd.PersistentFlags().StringVar(&senderPublicKey, "sender", "", "HTLC contract paramenter with sender PublicKey")
	deployHTLCCmd.PersistentFlags().StringVar(&recipientPublicKey, "recipient", "", "HTLC contract paramenter with recipientPublicKey")
	deployHTLCCmd.PersistentFlags().Uint64Var(&blockHeight, "blockHeight", 0, "HTLC contract locked value with blockHeight")
	deployHTLCCmd.PersistentFlags().StringVar(&hash, "hash", "", "HTLC contract locked value with hash")

	deployHTLCCmd.PersistentFlags().StringVar(&assetLocked, "assetLocked", "", "HTLC contract locked value with assetID")
	deployHTLCCmd.PersistentFlags().Uint64Var(&amountLocked, "amountLocked", 0, "HTLC contract locked value with amount")

	// call contract arguments
	callCmd.PersistentFlags().Uint64Var(&txFee, "txFee", 40000000, "contract transaction fee")

	callCmd.PersistentFlags().StringVar(&assetRequested, "assetRequested", "", "tradeoff contract paramenter with requested assetID")
	callCmd.PersistentFlags().Uint64Var(&amountRequested, "amountRequested", 0, "tradeoff contract paramenter with requested amount")
	callCmd.PersistentFlags().StringVar(&seller, "seller", "", "tradeoff contract paramenter with seller control-program")

	callCmd.PersistentFlags().StringVar(&assetLocked, "assetLocked", "", "tradeoff contract locked value with assetID")
	callCmd.PersistentFlags().Uint64Var(&amountLocked, "amountLocked", 0, "tradeoff contract locked value with amount")

	// call HTLC contract arguments
	callHTLCCmd.PersistentFlags().StringVar(&preimage, "preimage", "", "HTLC contract locked value with preimage")
}

var (
	txFee = uint64(0)

	// contract paramenters
	assetRequested  = ""
	amountRequested = uint64(0)
	seller          = ""
	cancelKey       = ""

	// contract locked value
	assetLocked  = ""
	amountLocked = uint64(0)

	// unlock contract paramenters
	contractUTXOID = ""
	buyer          = ""
)

var (
	senderPublicKey    = ""
	recipientPublicKey = ""
	blockHeight        = uint64(0)
	hash               = ""
	preimage           = ""
)

var deployCmd = &cobra.Command{
	Use:   "deploy <accountID> <password> [contract flags(paramenters and locked value)] [txFee flag]",
	Short: "deploy tradeoff contract",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		accountInfo := swap.AccountInfo{
			AccountID: args[0],
			Password:  args[1],
			TxFee:     txFee,
		}
		if len(accountInfo.AccountID) == 0 || len(accountInfo.Password) == 0 {
			fmt.Println("The part field of the structure AccountInfo is empty:", accountInfo)
			os.Exit(0)
		}

		contractArgs := swap.ContractArgs{
			AssetAmount: swap.AssetAmount{
				Asset:  assetRequested,
				Amount: amountRequested,
			},
			Seller:    seller,
			CancelKey: cancelKey,
		}
		if len(contractArgs.Asset) == 0 || contractArgs.Amount == uint64(0) || len(contractArgs.Seller) == 0 || len(contractArgs.CancelKey) == 0 {
			fmt.Println("The part field of the structure ContractArgs is empty:", contractArgs)
			os.Exit(0)
		}

		contractValue := swap.AssetAmount{
			Asset:  assetLocked,
			Amount: amountLocked,
		}
		if len(contractValue.Asset) == 0 || contractValue.Amount == uint64(0) {
			fmt.Println("The part field of the structure ContractValue AssetAmount is empty:", contractValue)
			os.Exit(0)
		}

		contractUTXOID, err := swap.DeployContract(accountInfo, contractArgs, contractValue)
		if err != nil {
			fmt.Println(err)
			os.Exit(0)
		}
		fmt.Println("--> contractUTXOID:", contractUTXOID)
	},
}

var callCmd = &cobra.Command{
	Use:   "call <accountID> <password> <buyer-program> <contractUTXOID> [txFee flag]",
	Short: "call tradeoff contract for asset swapping",
	Args:  cobra.ExactArgs(4),
	Run: func(cmd *cobra.Command, args []string) {
		accountInfo := swap.AccountInfo{
			AccountID: args[0],
			Password:  args[1],
			Receiver:  args[2],
			TxFee:     txFee,
		}
		if len(accountInfo.AccountID) == 0 || len(accountInfo.Password) == 0 || len(accountInfo.Receiver) == 0 {
			fmt.Println("The part field of the structure AccountInfo is empty:", accountInfo)
			os.Exit(0)
		}

		contractUTXOID := args[3]
		if len(contractUTXOID) == 0 {
			fmt.Println("contract utxoID is empty:", contractUTXOID)
			os.Exit(0)
		}

		program, contractValue, err := swap.ListUnspentOutputs(contractUTXOID)
		if err != nil {
			fmt.Println("list unspent outputs err:", err)
			os.Exit(0)
		}

		if len(contractValue.Asset) == 0 || contractValue.Amount == uint64(0) {
			fmt.Println("The part field of the structure ContractValue AssetAmount is empty:", contractValue)
			os.Exit(0)
		}

		contractArgs, err := swap.DecodeProgram(program)
		if err != nil {
			fmt.Println("decode program err:", err)
			os.Exit(0)
		}

		if len(contractArgs.Asset) == 0 || contractArgs.Amount == uint64(0) || len(contractArgs.Seller) == 0 {
			fmt.Println("The part field of the structure ContractArgs is empty:", contractArgs)
			os.Exit(0)
		}

		txID, err := swap.CallContract(accountInfo, contractUTXOID, *contractArgs, *contractValue)
		if err != nil {
			fmt.Println(err)
			os.Exit(0)
		}
		fmt.Println("--> txID:", txID)
	},
}

var deployHTLCCmd = &cobra.Command{
	Use:   "deployHTLC <accountID> <password> [contract flags(paramenters and locked value)] [txFee flag]",
	Short: "deploy HTLC contract",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		account := swap.AccountInfo{
			AccountID: args[0],
			Password:  args[1],
			TxFee:     txFee,
		}
		if len(account.AccountID) == 0 || len(account.Password) == 0 {
			fmt.Println("The part field of the structure AccountInfo is empty:", account)
			os.Exit(0)
		}

		contractArgs := swap.HTLCContractArgs{
			SenderPublicKey:    senderPublicKey,
			RecipientPublicKey: recipientPublicKey,
			BlockHeight:        blockHeight,
			Hash:               hash,
		}
		if len(contractArgs.SenderPublicKey) == 0 || len(contractArgs.RecipientPublicKey) == 0 || contractArgs.BlockHeight == uint64(0) || len(contractArgs.Hash) == 0 {
			fmt.Println("The part field of the structure ContractArgs is empty:", contractArgs)
			os.Exit(0)
		}

		contractValue := swap.AssetAmount{
			Asset:  assetLocked,
			Amount: amountLocked,
		}
		if len(contractValue.Asset) == 0 || contractValue.Amount == uint64(0) {
			fmt.Println("The part field of the structure ContractValue AssetAmount is empty:", contractValue)
			os.Exit(0)
		}

		contractUTXOID, err := swap.DeployHTLCContract(account, contractValue, contractArgs)
		if err != nil {
			fmt.Println(err)
			os.Exit(0)
		}
		fmt.Println("--> contractUTXOID:", contractUTXOID)
	},
}

var callHTLCCmd = &cobra.Command{
	Use:   "callHTLC <accountID> <password> <buyer-program> <contractUTXOID> <preimage> [txFee flag]",
	Short: "callHTLC HTLC contract for asset swapping",
	Args:  cobra.ExactArgs(5),
	Run: func(cmd *cobra.Command, args []string) {
		account := swap.AccountInfo{
			AccountID: args[0],
			Password:  args[1],
			Receiver:  args[2],
			TxFee:     txFee,
		}
		if len(account.AccountID) == 0 || len(account.Password) == 0 || len(account.Receiver) == 0 {
			fmt.Println("The part field of the structure Account is empty:", account)
			os.Exit(0)
		}

		contractUTXOID := args[3]
		if len(contractUTXOID) == 0 {
			fmt.Println("contract utxoID is empty:", contractUTXOID)
			os.Exit(0)
		}

		program, contractValue, err := swap.ListUnspentOutputs(contractUTXOID)
		if err != nil {
			fmt.Println("list unspent outputs err:", err)
			os.Exit(0)
		}

		if len(contractValue.Asset) == 0 || contractValue.Amount == uint64(0) {
			fmt.Println("The part field of the structure ContractValue AssetAmount is empty:", contractValue)
			os.Exit(0)
		}

		contractArgs, err := swap.DecodeHTLCProgram(program)
		if err != nil {
			fmt.Println("decode program err:", err)
			os.Exit(0)
		}

		if len(contractArgs.Hash) == 0 || len(contractArgs.RecipientPublicKey) == 0 || len(contractArgs.SenderPublicKey) == 0 || contractArgs.BlockHeight == uint64(0) {
			fmt.Println("The part field of the structure ContractArgs is empty:", contractArgs)
			os.Exit(0)
		}

		preimage := args[4]
		txID, err := swap.CallHTLCContract(account, contractUTXOID, preimage, *contractArgs, *contractValue)
		if err != nil {
			fmt.Println(err)
			os.Exit(0)
		}
		fmt.Println("--> txID:", txID)
	},
}
