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
	// compile contract locally
	equityCmd.PersistentFlags().BoolVar(&bin, strBin, false, "Binary of the contracts in hex.")
	equityCmd.PersistentFlags().BoolVar(&shift, strShift, false, "Function shift of the contracts.")
	equityCmd.PersistentFlags().BoolVar(&instance, strInstance, false, "Object of the Instantiated contracts.")
	equityCmd.PersistentFlags().BoolVar(&ast, strAst, false, "AST of the contracts.")
	equityCmd.PersistentFlags().BoolVar(&version, strVersion, false, "Version of equity compiler.")

	// build deploy htlc contract tx
	deployHTLCCmd.PersistentFlags().StringVar(&ip, "ip", "127.0.0.1", "network address")
	deployHTLCCmd.PersistentFlags().StringVar(&port, "port", "3000", "network port")

	// call HTLC contract arguments
	callHTLCCmd.PersistentFlags().StringVar(&ip, "ip", "127.0.0.1", "network address")
	callHTLCCmd.PersistentFlags().StringVar(&port, "port", "9888", "network port")

	// cancel HTLC contract arguments
	cancelHTLCCmd.PersistentFlags().StringVar(&ip, "ip", "127.0.0.1", "network address")
	cancelHTLCCmd.PersistentFlags().StringVar(&port, "port", "9888", "network port")

	// build deploy tradeoff contract tx
	deployTradeoffCmd.PersistentFlags().StringVar(&ip, "ip", "127.0.0.1", "network address")
	deployTradeoffCmd.PersistentFlags().StringVar(&port, "port", "3000", "network port")

	// call HTLC contract arguments
	callTradeoffCmd.PersistentFlags().StringVar(&ip, "ip", "127.0.0.1", "network address")
	callTradeoffCmd.PersistentFlags().StringVar(&port, "port", "9888", "network port")

	// cancel tradeoff contract arguments
	cancelTradeoffCmd.PersistentFlags().StringVar(&ip, "ip", "127.0.0.1", "network address")
	cancelTradeoffCmd.PersistentFlags().StringVar(&port, "port", "9888", "network port")

	// submit tx
	submitPaymentCmd.PersistentFlags().StringVar(&ip, "ip", "127.0.0.1", "network address")
	submitPaymentCmd.PersistentFlags().StringVar(&port, "port", "3000", "network port")
}

var (
	ip   = "127.0.0.1"
	port = "3000"
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

var deployHTLCCmd = &cobra.Command{
	Use:   "deployhtlc <guid> <lockedAsset> <lockedAmount> <contractProgram> [URL flags(ip and port)]",
	Short: "deploy HTLC contract",
	Args:  cobra.ExactArgs(4),
	Run: func(cmd *cobra.Command, args []string) {
		guid := args[0]
		if len(guid) == 0 {
			fmt.Println("The part field of guid is invalid:", guid)
			os.Exit(0)
		}

		lockedAsset := args[1]
		if _, err := hex.DecodeString(lockedAsset); err != nil || len(lockedAsset) != 64 {
			fmt.Println("The part field of lockedAsset is invalid:", lockedAsset)
			os.Exit(0)
		}

		lockedAmount, err := strconv.ParseUint(args[2], 10, 64)
		if err != nil {
			fmt.Println("parse locked amount err:", err)
			os.Exit(0)
		}

		contractProgram := args[3]
		if _, err := hex.DecodeString(contractProgram); err != nil || len(contractProgram) == 0 {
			fmt.Println("The part field of contractProgram is invalid:", contractProgram)
			os.Exit(0)
		}

		server := &swap.Server{
			IP:   ip,
			Port: port,
		}

		res, err := swap.BuildTx(server, guid, lockedAsset, contractProgram, lockedAmount)
		if err != nil {
			fmt.Println("build tx err:", err)
			os.Exit(0)
		}

		fmt.Println("build tx result:", res)
	},
}

var callHTLCCmd = &cobra.Command{
	Use:   "callhtlc <guid> <contractUTXOID> <contractAsset> <contractAmount> <receiver> [URL flags(ip and port)]",
	Short: "call HTLC contract",
	Args:  cobra.ExactArgs(5),
	Run: func(cmd *cobra.Command, args []string) {
		guid := args[0]
		if len(guid) == 0 {
			fmt.Println("The part field of guid is invalid:", guid)
			os.Exit(0)
		}

		contractUTXOID := args[1]
		if _, err := hex.DecodeString(contractUTXOID); err != nil || len(contractUTXOID) != 64 {
			fmt.Println("The part field of contractUTXOID is invalid:", contractUTXOID)
			os.Exit(0)
		}

		contractAsset := args[2]
		if _, err := hex.DecodeString(contractAsset); err != nil || len(contractAsset) != 64 {
			fmt.Println("The part field of contractAsset is invalid:", contractAsset)
			os.Exit(0)
		}

		contractAmount, err := strconv.ParseUint(args[3], 10, 64)
		if err != nil {
			fmt.Println("parse contract amount err:", err)
			os.Exit(0)
		}

		receiver := args[4]
		if len(receiver) == 0 {
			fmt.Println("The part field of receiver is invalid:", receiver)
			os.Exit(0)
		}

		server := &swap.Server{
			IP:   ip,
			Port: port,
		}

		res, err := swap.BuildUnlockedTx(server, guid, contractUTXOID, contractAsset, receiver, contractAmount)
		if err != nil {
			fmt.Println("build call htlc tx err:", err)
			os.Exit(0)
		}

		fmt.Println("build call htlc tx result:", res)
	},
}

var cancelHTLCCmd = &cobra.Command{
	Use:   "cancelhtlc <guid> <contractUTXOID> <contractAsset> <contractAmount> <receiver> [URL flags(ip and port)]",
	Short: "cancel HTLC contract",
	Args:  cobra.ExactArgs(5),
	Run: func(cmd *cobra.Command, args []string) {
		guid := args[0]
		if len(guid) == 0 {
			fmt.Println("The part field of guid is invalid:", guid)
			os.Exit(0)
		}

		contractUTXOID := args[1]
		if _, err := hex.DecodeString(contractUTXOID); err != nil || len(contractUTXOID) != 64 {
			fmt.Println("The part field of contractUTXOID is invalid:", contractUTXOID)
			os.Exit(0)
		}

		contractAsset := args[2]
		if _, err := hex.DecodeString(contractAsset); err != nil || len(contractAsset) != 64 {
			fmt.Println("The part field of contractAsset is invalid:", contractAsset)
			os.Exit(0)
		}

		contractAmount, err := strconv.ParseUint(args[3], 10, 64)
		if err != nil {
			fmt.Println("parse contract amount err:", err)
			os.Exit(0)
		}

		receiver := args[4]
		if len(receiver) == 0 {
			fmt.Println("The part field of receiver is invalid:", receiver)
			os.Exit(0)
		}

		server := &swap.Server{
			IP:   ip,
			Port: port,
		}

		res, err := swap.BuildUnlockedTx(server, guid, contractUTXOID, contractAsset, receiver, contractAmount)
		if err != nil {
			fmt.Println("build call htlc tx err:", err)
			os.Exit(0)
		}

		fmt.Println("build call htlc tx result:", res)
	},
}

var deployTradeoffCmd = &cobra.Command{
	Use:   "deploytradeoff <guid> <lockedAsset> <lockedAmount> <contractProgram> [URL flags(ip and port)]",
	Short: "deploy tradeoff contract",
	Args:  cobra.ExactArgs(4),
	Run: func(cmd *cobra.Command, args []string) {
		guid := args[0]
		if len(guid) == 0 {
			fmt.Println("The part field of guid is invalid:", guid)
			os.Exit(0)
		}

		lockedAsset := args[1]
		if _, err := hex.DecodeString(lockedAsset); err != nil || len(lockedAsset) != 64 {
			fmt.Println("The part field of lockedAsset is invalid:", lockedAsset)
			os.Exit(0)
		}

		lockedAmount, err := strconv.ParseUint(args[2], 10, 64)
		if err != nil {
			fmt.Println("parse locked amount err:", err)
			os.Exit(0)
		}

		contractProgram := args[3]
		if _, err := hex.DecodeString(contractProgram); err != nil || len(contractProgram) == 0 {
			fmt.Println("The part field of contractProgram is invalid:", contractProgram)
			os.Exit(0)
		}

		server := &swap.Server{
			IP:   ip,
			Port: port,
		}

		res, err := swap.BuildTx(server, guid, lockedAsset, contractProgram, lockedAmount)
		if err != nil {
			fmt.Println("build tx err:", err)
			os.Exit(0)
		}

		fmt.Println("build tx result:", res)
	},
}

var callTradeoffCmd = &cobra.Command{
	Use:   "calltradeoff <guid> <contractUTXOID> <assetRequested> <amountRequested> <seller> [URL flags(ip and port)]",
	Short: "call tradeoff contract",
	Args:  cobra.ExactArgs(5),
	Run: func(cmd *cobra.Command, args []string) {
		guid := args[0]
		if len(guid) == 0 {
			fmt.Println("The part field of guid is invalid:", guid)
			os.Exit(0)
		}

		contractUTXOID := args[1]
		if _, err := hex.DecodeString(contractUTXOID); err != nil || len(contractUTXOID) != 64 {
			fmt.Println("The part field of contractUTXOID is invalid:", contractUTXOID)
			os.Exit(0)
		}

		assetRequested := args[2]
		if _, err := hex.DecodeString(assetRequested); err != nil || len(assetRequested) != 64 {
			fmt.Println("The part field of assetRequested is invalid:", assetRequested)
			os.Exit(0)
		}

		amountRequested, err := strconv.ParseUint(args[3], 10, 64)
		if err != nil {
			fmt.Println("parse amountRequested err:", err)
			os.Exit(0)
		}

		seller := args[4]
		if _, err := hex.DecodeString(seller); err != nil || len(seller) == 0 {
			fmt.Println("The part field of seller is invalid:", seller)
			os.Exit(0)
		}

		server := &swap.Server{
			IP:   ip,
			Port: port,
		}

		res, err := swap.BuildCallTradeoffTx(server, guid, contractUTXOID, seller, assetRequested, amountRequested)
		if err != nil {
			fmt.Println("build call tradeoff tx err:", err)
			os.Exit(0)
		}

		fmt.Println("build call tradeoff tx result:", res)
	},
}

var cancelTradeoffCmd = &cobra.Command{
	Use:   "canceltradeoff <guid> <contractUTXOID> <lockedAsset> <lockedAmount> <receiver> [URL flags(ip and port)]",
	Short: "cancel tradeoff contract",
	Args:  cobra.ExactArgs(5),
	Run: func(cmd *cobra.Command, args []string) {
		guid := args[0]
		if len(guid) == 0 {
			fmt.Println("The part field of guid is invalid:", guid)
			os.Exit(0)
		}

		contractUTXOID := args[1]
		if _, err := hex.DecodeString(contractUTXOID); err != nil || len(contractUTXOID) != 64 {
			fmt.Println("The part field of contractUTXOID is invalid:", contractUTXOID)
			os.Exit(0)
		}

		lockedAsset := args[2]
		if _, err := hex.DecodeString(lockedAsset); err != nil || len(lockedAsset) != 64 {
			fmt.Println("The part field of lockedAsset is invalid:", lockedAsset)
			os.Exit(0)
		}

		lockedAmount, err := strconv.ParseUint(args[3], 10, 64)
		if err != nil {
			fmt.Println("parse contract amount err:", err)
			os.Exit(0)
		}

		receiver := args[4]
		if len(receiver) == 0 {
			fmt.Println("The part field of receiver is invalid:", receiver)
			os.Exit(0)
		}

		server := &swap.Server{
			IP:   ip,
			Port: port,
		}

		res, err := swap.BuildUnlockedTx(server, guid, contractUTXOID, lockedAsset, receiver, lockedAmount)
		if err != nil {
			fmt.Println("build call htlc tx err:", err)
			os.Exit(0)
		}

		fmt.Println("build call htlc tx result:", res)
	},
}

var signMessageCmd = &cobra.Command{
	Use:   "sign <xprv> <message>",
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

		res, err := swap.SignMessage(message, xprv)
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

var submitPaymentCmd = &cobra.Command{
	Use:   "submit <action> <guid> <rawTx> [spend parameters] [URL flags(ip and port)]",
	Short: "submit a payment",
	Args:  cobra.MinimumNArgs(3),
	Run: func(cmd *cobra.Command, args []string) {
		action := args[0]
		guid := args[1]
		if len(guid) == 0 {
			fmt.Println("The part field of guid is invalid:", guid)
			os.Exit(0)
		}

		rawTx := args[2]
		if _, err := hex.DecodeString(rawTx); err != nil {
			fmt.Println("The part field of rawTx is invalid:", rawTx)
			os.Exit(0)
		}

		// witness_arguments
		wa := []string{}
		was := [][]string{}

		switch action {
		case "deployhtlc":
			for i := 3; i < len(args); i++ {
				wa = append(wa, args[i])
				was = append(was, wa)
				wa = []string{}
			}
		case "callhtlc":
			// callhtlc need 6 arguments at least
			if len(args) < 6 {
				fmt.Println("callhtlc need 6 arguments at least, len(args) is:", len(args))
				os.Exit(0)
			}

			wa = append(wa, args[3], args[4], "")
			was = append(was, wa)
			wa = []string{}

			for i := 5; i < len(args); i++ {
				wa = append(wa, args[i])
				was = append(was, wa)
				wa = []string{}
			}
		case "cancelhtlc":
			// callhtlc need 5 arguments at least
			if len(args) < 5 {
				fmt.Println("callhtlc need 5 arguments at least, len(args) is:", len(args))
				os.Exit(0)
			}

			wa = append(wa, args[3], "01")
			was = append(was, wa)
			wa = []string{}

			for i := 4; i < len(args); i++ {
				wa = append(wa, args[i])
				was = append(was, wa)
				wa = []string{}
			}
		case "deploytradeoff":
			for i := 3; i < len(args); i++ {
				wa = append(wa, args[i])
				was = append(was, wa)
				wa = []string{}
			}
		case "calltradeoff":
			// callhtlc need 4 arguments at least
			if len(args) < 4 {
				fmt.Println("calltradeoff need 4 arguments at least, len(args) is:", len(args))
				os.Exit(0)
			}

			wa = append(wa, "")
			was = append(was, wa)
			wa = []string{}

			for i := 3; i < len(args); i++ {
				wa = append(wa, args[i])
				was = append(was, wa)
				wa = []string{}
			}
		case "canceltradeoff":
			// calltradeoff need 5 arguments at least
			if len(args) < 5 {
				fmt.Println("calltradeoff need 5 arguments at least, len(args) is:", len(args))
				os.Exit(0)
			}

			wa = append(wa, args[3], "01")
			was = append(was, wa)
			wa = []string{}

			for i := 4; i < len(args); i++ {
				wa = append(wa, args[i])
				was = append(was, wa)
				wa = []string{}
			}
		default:
			fmt.Println("action is invalid:", action)
			os.Exit(0)
		}

		server := &swap.Server{
			IP:   ip,
			Port: port,
		}
		res, err := swap.SubmitPayment(server, guid, rawTx, "", was)
		if err != nil {
			fmt.Println("submit tx err:", err)
			os.Exit(0)
		}

		fmt.Printf("submit %s tx result: %s\n", action, res)
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
