package swap

// var (
// 	localURL = "http://127.0.0.1:9888/"

// 	buildTransactionURL     = localURL + "build-transaction"
// 	getTransactionURL       = localURL + "get-transaction"
// 	signTransactionURL      = localURL + "sign-transaction"
// 	decodeRawTransactionURL = localURL + "decode-raw-transaction"
// 	submitTransactionURL    = localURL + "submit-transaction"
// 	compileURL              = localURL + "compile"
// 	decodeProgramURL        = localURL + "decode-program"
// 	signMessageURl          = localURL + "sign-message"
// 	listAccountsURL         = localURL + "list-accounts"
// 	listAddressesURL        = localURL + "list-addresses"
// 	listBalancesURL         = localURL + "list-balances"
// 	listPubkeysURL          = localURL + "list-pubkeys"
// 	listUnspentOutputsURL   = localURL + "list-unspent-outputs"
// )

var (
	// localURL = "http://127.0.0.1:9888/"

	buildTransactionURL     = "build-transaction"
	getTransactionURL       = "get-transaction"
	signTransactionURL      = "sign-transaction"
	decodeRawTransactionURL = "decode-raw-transaction"
	submitTransactionURL    = "submit-transaction"
	compileURL              = "compile"
	decodeProgramURL        = "decode-program"
	signMessageURl          = "sign-message"
	listAccountsURL         = "list-accounts"
	listAddressesURL        = "list-addresses"
	listBalancesURL         = "list-balances"
	listPubkeysURL          = "list-pubkeys"
	listUnspentOutputsURL   = "list-unspent-outputs"
)

type AccountInfo struct {
	AccountID string
	Password  string
	Receiver  string
	TxFee     uint64
}

type AssetAmount struct {
	Asset  string
	Amount uint64
}

type ContractArgs struct {
	AssetAmount
	Seller    string
	CancelKey string
}
