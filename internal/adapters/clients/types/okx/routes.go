package okx

type Endpoint string

const (
	OrderBook      Endpoint = "/api/v5/market/books"
	Balances       Endpoint = "/api/v5/account/balance"
	OpenOrders     Endpoint = "/api/v5/trade/orders-pending"
	OrderHistory   Endpoint = "/api/v5/trade/orders-history"
	FillsHistory   Endpoint = "/api/v5/trade/fills-history"
	Order          Endpoint = "/api/v5/trade/order"
	CancelOrder    Endpoint = "/api/v5/trade/cancel-order"
	AmendOrder     Endpoint = "/api/v5/trade/amend-order"
	CancelAllAfter Endpoint = "/api/v5/trade/cancel-all-after"
)
