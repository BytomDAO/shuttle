package swap

var (
	buildTransactionURL     = "build-transaction"
	getTransactionURL       = "api/v1/btm/merchant/get-transaction" // get-transaction blockcenter url
	signTransactionURL      = "sign-transaction"
	decodeRawTransactionURL = "decode-raw-transaction"
	submitTransactionURL    = "api/v1/btm/merchant/submit-payment"
	compileURL              = "compile"
	decodeProgramURL        = "decode-program"
	signMessageURl          = "sign-message"
	listAccountsURL         = "list-accounts"
	listAddressesURL        = "list-addresses"
	listBalancesURL         = "list-balances"
	listPubkeysURL          = "list-pubkeys"
	listUnspentOutputsURL   = "list-unspent-outputs"
)

const (
	BTMAssetID = "ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff"
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
