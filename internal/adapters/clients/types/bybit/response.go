package bybit

//go:generate easyjson -all response.go

type CommonResponse struct {
	Code int    `json:"retCode"`
	Msg  string `json:"retMsg"`
	Info struct {
	} `json:"retExtInfo"`
	Time int64 `json:"time"`
}
