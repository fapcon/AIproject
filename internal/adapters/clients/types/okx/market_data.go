package okx

//easyjson:json
type OrderBookResponse struct {
	/*{
		"code": "0",
		"msg": "",
		"data": [
			{
				"asks": [
					[
							"41006.8",
							"0.60038921",
							"0",
							"1"
					]
				],
				"bids": [
					[
							"41006.3",
							"0.30178218",
							"0",
							"2"
					]
				],
				"ts": "1629966436396"
			}
	  	]
	  }*/
	Code string `json:"code"`
	Msg  string `json:"msg"`
	Data []struct {
		Asks [][]string `json:"asks"`
		Bids [][]string `json:"bids"`
		TS   string     `json:"ts"`
	} `json:"data"`
}
