package swap

var (
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
