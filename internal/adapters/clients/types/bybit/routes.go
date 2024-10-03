package bybit

type Endpoint string

const (
	GetOrderBook    Endpoint = "/v5/market/orderbook"
	CreateOrder     Endpoint = "/v5/order/create"
	AmendOrder      Endpoint = "/v5/order/amend"
	CancelOrder     Endpoint = "/v5/order/cancel"
	GetOpenedOrders Endpoint = "/v5/order/realtime"
	CancelAllOrders Endpoint = "/v5/order/cancel-all"
	GetBalances     Endpoint = "/v5/account/wallet-balance"
	GetOrderHistory Endpoint = "/v5/order/history"
)
