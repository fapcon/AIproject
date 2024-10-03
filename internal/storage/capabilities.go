package storage

type Capability string

const (
	CapabilityOKXAPIClient                         Capability = "okx_api_client"
	CapabilityOKXSpotTradingWebsocketClient        Capability = "okx_spot_trading_websocket_client"
	CapabilityOKXMarketDataWebsocketClient         Capability = "okx_market_data_websocket_client"
	CapabilityOKXUserStreamWebsocketClient         Capability = "okx_user_stream_websocket_client"
	CapabilityNameResolver                         Capability = "name_resolver"
	CapabilityInstrumentDetails                    Capability = "instrument_details"
	CapabilityOrderBookRead                        Capability = "order_book_read"
	CapabilityOrderBookWrite                       Capability = "order_book_write"
	CapabilityPistonClient                         Capability = "piston_client"
	CapabilityPistonIDCache                        Capability = "piston_id_cache"
	CapabilityOrderMoveCache                       Capability = "order_move_cache"
	CapabilitySubscribeInstrumentsRequestsProducer Capability = "subscribe_instruments_requests_producer"
	CapabilityBybitAPIClient                       Capability = "bybit_api_client"
	CapabilityBybitMarketDataWebsocketClient       Capability = "bybit_market_data_websocket_client"
	CapabilityBybitUserStreamWebsocketClient       Capability = "bybit_user_stream_websocket_client"
	CapabilityNewOrderRequestsConsumer             Capability = "new_order_request_consumer"
)
