package main

import (
	"bufio"
	"encoding/hex"
	"fmt"
	"os"
	"strconv"

	"github.com/equity/compiler"
	equ "github.com/equity/equity/util"
	"github.com/spf13/cobra"

	"github.com/shuttle/swap"
)

func init() {
	// deploy contract arguments
	deployTradeoffCmd.PersistentFlags().Uint64Var(&txFee, "txFee", 40000000, "contract transaction fee")
	deployTradeoffCmd.PersistentFlags().StringVar(&assetRequested, "assetRequested", "", "tradeoff contract paramenter with requested assetID")
	deployTradeoffCmd.PersistentFlags().Uint64Var(&amountRequested, "amountRequested", 0, "tradeoff contract paramenter with requested amount")
	deployTradeoffCmd.PersistentFlags().StringVar(&seller, "seller", "", "tradeoff contract paramenter with seller control-program")
	deployTradeoffCmd.PersistentFlags().StringVar(&cancelKey, "cancelKey", "", "tradeoff contract paramenter with seller pubkey for cancelling the contract")
	deployTradeoffCmd.PersistentFlags().StringVar(&assetLocked, "assetLocked", "", "tradeoff contract locked value with assetID")
	deployTradeoffCmd.PersistentFlags().Uint64Var(&amountLocked, "amountLocked", 0, "tradeoff contract locked value with amount")
	deployTradeoffCmd.PersistentFlags().StringVar(&ip, "ip", "127.0.0.1", "network address")
	deployTradeoffCmd.PersistentFlags().StringVar(&port, "port", "9888", "network port")

	// deploy HTLC contract arguments
	deployHTLCCmd.PersistentFlags().Uint64Var(&txFee, "txFee", 40000000, "contract transaction fee")
	deployHTLCCmd.PersistentFlags().StringVar(&senderPublicKey, "sender", "", "HTLC contract paramenter with sender PublicKey")
	deployHTLCCmd.PersistentFlags().StringVar(&recipientPublicKey, "recipient", "", "HTLC contract paramenter with recipientPublicKey")
	deployHTLCCmd.PersistentFlags().Uint64Var(&blockHeight, "blockHeight", 0, "HTLC contract locked value with blockHeight")
	deployHTLCCmd.PersistentFlags().StringVar(&hash, "hash", "", "HTLC contract locked value with hash")
	deployHTLCCmd.PersistentFlags().StringVar(&assetLocked, "assetLocked", "", "HTLC contract locked value with assetID")
	deployHTLCCmd.PersistentFlags().Uint64Var(&amountLocked, "amountLocked", 0, "HTLC contract locked value with amount")
	deployHTLCCmd.PersistentFlags().StringVar(&ip, "ip", "127.0.0.1", "network address")
	deployHTLCCmd.PersistentFlags().StringVar(&port, "port", "9888", "network port")

	// call contract arguments
	callTradeoffCmd.PersistentFlags().Uint64Var(&txFee, "txFee", 40000000, "contract transaction fee")
	callTradeoffCmd.PersistentFlags().StringVar(&ip, "ip", "127.0.0.1", "network address")
	callTradeoffCmd.PersistentFlags().StringVar(&port, "port", "9888", "network port")

	// call HTLC contract arguments
	callHTLCCmd.PersistentFlags().Uint64Var(&txFee, "txFee", 40000000, "contract transaction fee")
	callHTLCCmd.PersistentFlags().StringVar(&ip, "ip", "127.0.0.1", "network address")
	callHTLCCmd.PersistentFlags().StringVar(&port, "port", "9888", "network port")

	// cancel tradeoff contract arguments
	cancelTradeoffCmd.PersistentFlags().Uint64Var(&txFee, "txFee", 40000000, "contract transaction fee")
	cancelTradeoffCmd.PersistentFlags().StringVar(&ip, "ip", "127.0.0.1", "network address")
	cancelTradeoffCmd.PersistentFlags().StringVar(&port, "port", "9888", "network port")

	// cancel HTLC contract arguments
	cancelHTLCCmd.PersistentFlags().Uint64Var(&txFee, "txFee", 40000000, "contract transaction fee")
	cancelHTLCCmd.PersistentFlags().StringVar(&ip, "ip", "127.0.0.1", "network address")
	cancelHTLCCmd.PersistentFlags().StringVar(&port, "port", "9888", "network port")

	// compile contract locally
	equityCmd.PersistentFlags().BoolVar(&bin, strBin, false, "Binary of the contracts in hex.")
	equityCmd.PersistentFlags().BoolVar(&shift, strShift, false, "Function shift of the contracts.")
	equityCmd.PersistentFlags().BoolVar(&instance, strInstance, false, "Object of the Instantiated contracts.")
	equityCmd.PersistentFlags().BoolVar(&ast, strAst, false, "AST of the contracts.")
	equityCmd.PersistentFlags().BoolVar(&version, strVersion, false, "Version of equity compiler.")

	// build deploy contract tx
	buildTxCmd.PersistentFlags().StringVar(&ip, "ip", "127.0.0.1", "network address")
	buildTxCmd.PersistentFlags().StringVar(&port, "port", "3000", "network port")
}

var (
	txFee = uint64(0)
	ip    = "127.0.0.1"
	port  = "3000"

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

var (
	strBin      string = "bin"
	strShift    string = "shift"
	strInstance string = "instance"
	strAst      string = "ast"
	strVersion  string = "version"
)

var (
	bin      = false
	shift    = false
	instance = false
	ast      = false
	version  = false
)

var deployTradeoffCmd = &cobra.Command{
	Use:   "deployTradeoff <accountID> <password> [contract flags(paramenters and locked value)] [txFee flag] [URL flags(ip and port)]",
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

		server := &swap.Server{
			IP:   ip,
			Port: port,
		}

		contractUTXOID, err := swap.DeployTradeoffContract(server, accountInfo, contractArgs, contractValue)
		if err != nil {
			fmt.Println(err)
			os.Exit(0)
		}
		fmt.Println("--> contractUTXOID:", contractUTXOID)
	},
}

var buildTxCmd = &cobra.Command{
	Use:   "build <guid> <outputID> <lockedAsset> <contractProgram> <lockedAmount> [URL flags(ip and port)]",
	Short: "build contract",
	Args:  cobra.ExactArgs(5),
	Run: func(cmd *cobra.Command, args []string) {
		guid := args[0]
		if len(guid) == 0 {
			fmt.Println("The part field of guid is invalid:", guid)
			os.Exit(0)
		}

		outputID := args[1]
		if _, err := hex.DecodeString(outputID); err != nil || len(outputID) != 64 {
			fmt.Println("The part field of outputID is invalid:", outputID)
			os.Exit(0)
		}

		lockedAsset := args[2]
		if _, err := hex.DecodeString(lockedAsset); err != nil || len(lockedAsset) != 64 {
			fmt.Println("The part field of lockedAsset is invalid:", lockedAsset)
			os.Exit(0)
		}

		contractProgram := args[3]
		if _, err := hex.DecodeString(contractProgram); err != nil || len(contractProgram) == 0 {
			fmt.Println("The part field of contractProgram is invalid:", contractProgram)
			os.Exit(0)
		}

		lockedAmount, err := strconv.ParseUint(args[4], 10, 64)
		if err != nil {
			fmt.Println("parse locked amount err:", err)
			os.Exit(0)
		}

		server := &swap.Server{
			IP:   ip,
			Port: port,
		}

		res, err := swap.BuildTx(server, guid, outputID, lockedAsset, contractProgram, lockedAmount)
		if err != nil {
			fmt.Println("build tx err:", err)
			os.Exit(0)
		}

		fmt.Println("build tx result:", res)
	},
}

var signMessageCmd = &cobra.Command{
	Use:   "sign [xprv] [message]",
	Short: "sign message",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		xprv := args[0]
		if _, err := hex.DecodeString(xprv); err != nil || len(xprv) != 128 {
			fmt.Println("The part field of xprv is invalid:", xprv)
			os.Exit(0)
		}

		message := args[1]
		if _, err := hex.DecodeString(message); err != nil || len(message) != 64 {
			fmt.Println("The part field of message is invalid:", message)
			os.Exit(0)
		}

		res, err := swap.SignMsg(message, xprv)
		if err != nil {
			fmt.Println("sign message err:", err)
			os.Exit(0)
		}

		fmt.Printf("\nsign result:\n"+
			"xprv: %s\n"+
			"message: %s\n"+
			"signature: %s\n",
			xprv, message, res)
	},
}

var callTradeoffCmd = &cobra.Command{
	Use:   "callTradeoff <accountID> <password> <buyer-program> <contractUTXOID> [txFee flag] [URL flags(ip and port)]",
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

		server := &swap.Server{
			IP:   ip,
			Port: port,
		}

		txID, err := swap.CallTradeoffContract(server, accountInfo, contractUTXOID)
		if err != nil {
			fmt.Println(err)
			os.Exit(0)
		}
		fmt.Println("--> txID:", txID)
	},
}

var cancelTradeoffCmd = &cobra.Command{
	Use:   "cancelTradeoff <accountID> <password> <redeem-program> <contractUTXOID> [txFee flag] [URL flags(ip and port)]",
	Short: "cancel tradeoff contract for asset swapping",
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

		server := &swap.Server{
			IP:   ip,
			Port: port,
		}

		txID, err := swap.CancelTradeoffContract(server, accountInfo, contractUTXOID)
		if err != nil {
			fmt.Println(err)
			os.Exit(0)
		}
		fmt.Println("--> txID:", txID)
	},
}

var deployHTLCCmd = &cobra.Command{
	Use:   "deployHTLC <accountID> <password> [contract flags(paramenters and locked value)] [txFee flag] [URL flags(ip and port)]",
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

		server := &swap.Server{
			IP:   ip,
			Port: port,
		}

		contractUTXOID, err := swap.DeployHTLCContract(server, account, contractValue, contractArgs)
		if err != nil {
			fmt.Println(err)
			os.Exit(0)
		}
		fmt.Println("--> contractUTXOID:", contractUTXOID)
	},
}

var callHTLCCmd = &cobra.Command{
	Use:   "callHTLC <accountID> <password> <buyer-program> <preimage> <contractUTXOID> [txFee flag] [URL flags(ip and port)]",
	Short: "call HTLC contract for asset swapping",
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

		contractUTXOID := args[4]
		if len(contractUTXOID) == 0 {
			fmt.Println("contract utxoID is empty:", contractUTXOID)
			os.Exit(0)
		}

		preimage := args[3]
		server := &swap.Server{
			IP:   ip,
			Port: port,
		}
		txID, err := swap.CallHTLCContract(server, account, contractUTXOID, preimage)
		if err != nil {
			fmt.Println(err)
			os.Exit(0)
		}
		fmt.Println("--> txID:", txID)
	},
}

var cancelHTLCCmd = &cobra.Command{
	Use:   "cancelHTLC <accountID> <password> <redeem-program> <contractUTXOID> [txFee flag] [URL flags(ip and port)]",
	Short: "cancel HTLC contract for asset swapping",
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

		server := &swap.Server{
			IP:   ip,
			Port: port,
		}

		txID, err := swap.CancelHTLCContract(server, accountInfo, contractUTXOID)
		if err != nil {
			fmt.Println(err)
			os.Exit(0)
		}
		fmt.Println("--> txID:", txID)
	},
}

var equityCmd = &cobra.Command{
	Use:     "equity <input_file>",
	Short:   "equity commandline compiler",
	Example: "equity contract_name [contract_args...] --bin --instance",
	Args:    cobra.RangeArgs(0, 100),
	Run: func(cmd *cobra.Command, args []string) {
		if version {
			version := compiler.VersionWithCommit(compiler.GitCommit)
			fmt.Println("equity, the equity compiler commandline interface")
			fmt.Println("Version:", version)
			os.Exit(0)
		}

		if len(args) < 1 {
			cmd.Usage()
			os.Exit(0)
		}

		if err := handleCompiled(args); err != nil {
			os.Exit(-1)
		}
	},
}

func handleCompiled(args []string) error {
	contractFile, err := os.Open(args[0])
	if err != nil {
		fmt.Printf("An error [%v] occurred on opening the file, please check whether the file exists or can be accessed.\n", err)
		return err
	}
	defer contractFile.Close()

	reader := bufio.NewReader(contractFile)
	contracts, err := compiler.Compile(reader)
	if err != nil {
		fmt.Println("Compile contract failed:", err)
		return err
	}

	// Print the result for all contracts
	for i, contract := range contracts {
		fmt.Printf("======= %v =======\n", contract.Name)
		if bin {
			fmt.Println("Binary:")
			fmt.Printf("%v\n\n", hex.EncodeToString(contract.Body))
		}

		if shift {
			fmt.Println("Clause shift:")
			clauseMap, err := equ.Shift(contract)
			if err != nil {
				fmt.Println("Statistics contract clause shift error:", err)
				return err
			}

			for clause, shift := range clauseMap {
				fmt.Printf("    %s:  %v\n", clause, shift)
			}
			fmt.Printf("\nNOTE: \n    If the contract contains only one clause, Users don't need clause selector when unlock contract." +
				"\n    Furthermore, there is no signification for ending clause shift except for display.\n\n")
		}

		if instance {
			if i != len(contracts)-1 {
				continue
			}

			fmt.Println("Instantiated program:")
			if len(args)-1 < len(contract.Params) {
				fmt.Printf("Error: The number of input arguments %d is less than the number of contract parameters %d\n", len(args)-1, len(contract.Params))
				usage := fmt.Sprintf("Usage:\n  equity %s", args[0])
				for _, param := range contract.Params {
					usage = usage + " <" + param.Name + ">"
				}
				fmt.Printf("%s\n\n", usage)
				return err
			}

			contractArgs, err := equ.ConvertArguments(contract, args[1:len(contract.Params)+1])
			if err != nil {
				fmt.Println("Convert arguments into contract parameters error:", err)
				return err
			}

			instantProg, err := equ.InstantiateContract(contract, contractArgs)
			if err != nil {
				fmt.Println("Instantiate contract error:", err)
				return err
			}
			fmt.Printf("%v\n\n", hex.EncodeToString(instantProg))
		}

		if ast {
			fmt.Println("Ast:")
			rawData, err := equ.JSONMarshal(contract, true)
			if err != nil {
				fmt.Println("Marshal the struct of contract to json error:", err)
				return err
			}
			fmt.Println(string(rawData))
		}
	}

	return nil
}
