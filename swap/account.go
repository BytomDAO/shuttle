package swap

type Account struct {
	AccountID    string `json:"id"`
	AccountAlias string `json:"alias"`
}

type AccountsResponse struct {
	Status string    `json:"status"`
	Data   []Account `json:"data"`
}

// func ListAccounts() []Account {
// 	data := []byte(`{}`)
// 	body := request(listAccountsURL, data)

// 	accountsResp := new(AccountsResponse)
// 	if err := json.Unmarshal(body, accountsResp); err != nil {
// 		fmt.Println(err)
// 	}
// 	return accountsResp.Data
// }

type Address struct {
	AccountAlias   string `json:"account_alias"`
	AccountID      string `json:"account_id"`
	Address        string `json:"address"`
	ControlProgram string `json:"control_program"`
	Change         bool   `json:"change"`
	KeyIndex       uint64 `json:"key_index"`
}

type AddressesResponse struct {
	Status string    `json:"status"`
	Data   []Address `json:"data"`
}

// func ListAddresses(accountAlias string) []Address {
// 	data := []byte(`{"account_alias": "` + accountAlias + `"}`)
// 	body := request(listAddressesURL, data)

// 	addresses := new(AddressesResponse)
// 	if err := json.Unmarshal(body, addresses); err != nil {
// 		fmt.Println(err)
// 	}
// 	return addresses.Data
// }

type Balance struct {
	AccountID string `json:"account_id"`
	Amount    uint64 `json:"amount"`
}

type BalancesResponse struct {
	Status string    `json:"status"`
	Data   []Balance `json:"data"`
}

// func ListBalances(accountAlias string) []Balance {
// 	data := []byte(`{"account_alias": "` + accountAlias + `"}`)
// 	body := request(listBalancesURL, data)

// 	balancesResp := new(BalancesResponse)
// 	if err := json.Unmarshal(body, balancesResp); err != nil {
// 		fmt.Println(err)
// 	}
// 	return balancesResp.Data
// }

type PubkeyInfo struct {
	Pubkey string   `json:"pubkey"`
	Path   []string `json:"derivation_path"`
}

type KeyInfo struct {
	XPubkey     string       `json:"root_xpub"`
	PubkeyInfos []PubkeyInfo `json:"pubkey_infos"`
}

type PubkeysResponse struct {
	Status string  `json:"status"`
	Data   KeyInfo `json:"data"`
}

// func ListPubkeys(accountAlias string) KeyInfo {
// 	data := []byte(`{"account_alias": "` + accountAlias + `"}`)
// 	body := request(listPubkeysURL, data)

// 	pubkeysResp := new(PubkeysResponse)
// 	if err := json.Unmarshal(body, pubkeysResp); err != nil {
// 		fmt.Println(err)
// 	}
// 	return pubkeysResp.Data
// }
