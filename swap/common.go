package swap

var (
	buildTransactionURL  = "api/v1/btm/merchant/build-transaction"
	getTransactionURL    = "api/v1/btm/merchant/get-transaction" // get-transaction blockcenter url
	submitTransactionURL = "api/v1/btm/merchant/submit-payment"
)

const (
	BTMAssetID    = "ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff"
	fee           = uint64(40000000)
	confirmations = uint64(1)
)

type AssetAmount struct {
	Asset  string
	Amount uint64
}
