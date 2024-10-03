package observation

import "studentgit.kata.academy/quant/torque/pkg/metrics"

type Name = string

const (
	Namespace = "arb"
	Subsystem = "gateways"
)

const (
	WebsocketReconnect              metrics.Name = "websocket_reconnect"
	WebsocketEventRoundTripDuration metrics.Name = "event_round_trip_duration"
	WebsocketEventDecodeDuration    metrics.Name = "event_decode_duration"

	MakerDataUpdate        Name = "maker_data_update"
	OrderBook              Name = "order_book"
	RestOrderRequest       Name = "rest_order_request"
	RestOrderCancelRequest Name = "rest_order_cancel_request"
	RestOpenOrdersRequest  Name = "rest_open_orders_request"
	RestAccountInfoRequest Name = "rest_account_info_request"
)
