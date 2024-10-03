package bybit

//go:generate easyjson -all wallet.go
type GetBalancesResponse struct {
	CommonResponse
	Result struct {
		List []struct {
			Coin []struct {
				WalletBalance string `json:"walletBalance"`
				Coin          string `json:"coin"`
			} `json:"coin"`
		} `json:"list"`
	} `json:"result"`
}

type AccountType string

// bybit account types.
const (
	UnifiedAccount  AccountType = "UNIFIED"
	ContractAccount AccountType = "CONTRACT"
	SpotAccount     AccountType = "SPOT"
)
