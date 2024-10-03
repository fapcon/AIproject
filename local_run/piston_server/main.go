package main

import (
	"log"
	"net"
	"os"
	"time"

	"github.com/oklog/ulid/v2"
	"google.golang.org/grpc"
	pb "studentgit.kata.academy/quant/torque/internal/adapters/grpc_clients/types/piston"
	"studentgit.kata.academy/quant/torque/pkg/logster"
)

type server struct {
	pb.UnimplementedPistonServer
	logger     logster.Logger
	orderSent  bool
	id         string
	exchangeID string
}

func (s *server) ConnectMDGateway(stream pb.Piston_ConnectMDGatewayServer) error {
	go func() {
		for {
			<-time.After(5 * time.Second)
			log.Println("Sending SubscribeMessage")
			if err := stream.Send(&pb.MarketdataMessage{
				MsgType: pb.MarketdataMsgType_SUBSCRIBE,
				Message: &pb.MarketdataMessage_SubscribeMessage{
					SubscribeMessage: &pb.SubscribeMessage{
						Instruments: []string{"BTC/USDT"},
					},
				},
			}); err != nil {
				log.Fatalf("Failed to send message: %v", err)
			}
			return
		}
	}()

	for {
		in, err := stream.Recv()
		if err != nil {
			return err
		}

		_ = in
		s.logger.Printf("Received Asks: %v", in.GetBookMessage().GetBook().GetSnapshot().GetAsks()[:2])
		s.logger.Printf("Received Bids: %v", in.GetBookMessage().GetBook().GetSnapshot().GetBids()[:2])
	}
}

func (s *server) ConnectTradingGateway(stream pb.Piston_ConnectTradingGatewayServer) error {
	go func() {
		for {
			<-time.After(5 * time.Second)

			//s.logger.Println("Sending TradingMessageType_RequestInstrumentsDetails")
			//if err := stream.Send(&pb.TradingMessage{
			//	Type: pb.TradingMessageType_RequestInstrumentsDetails,
			//	Data: &pb.TradingMessage_InstrumentsDetailsRequest{},
			//}); err != nil {
			//	s.logger.Infof("Failed to send TradingMessageType_RequestInstrumentsDetails: %v", err)
			//}
			//
			//<-time.After(5 * time.Second)

			//log.Println("Sending TradingHistoryType_OrderMessage OrderHistoryRequest")
			//if err := stream.Send(&pb.TradingMessage{
			//	Type: pb.TradingMessageType_RequestTradeHistory,
			//	Data: &pb.TradingMessage_TradeHistoryReq{TradeHistoryReq: &pb.TradeHistoryRequest{
			//		Id:   1,
			//		From: time.Now().AddDate(0, -1, 0).UnixNano(),
			//		To:   time.Now().UnixNano(),
			//	}},
			//}); err != nil {
			//	log.Fatalf("Failed to send TradingMessage_Order OrderHistoryRequest: %v", err)
			//}
			//
			//<-time.After(5 * time.Second)

			//log.Println("Sending TradingMessageType_RequestOpenOrders")
			//if err := stream.Send(&pb.TradingMessage{
			//	Type: pb.TradingMessageType_RequestOpenOrders,
			//	Data: &pb.TradingMessage_OpenOrdersRequest{OpenOrdersRequest: &pb.OpenOrdersRequest{
			//		RequestID: 100,
			//	}},
			//}); err != nil {
			//	log.Fatalf("Failed to send RequestOpenOrders: %v", err)
			//}

			//<-time.After(2 * time.Second)
			//log.Println("Sending TradingMessageType_BalancesRequest")
			//if err := stream.Send(&pb.TradingMessage{
			//	Type: pb.TradingMessageType_RequestBalance,
			//}); err != nil {
			//	log.Fatalf("Failed to send RequestBalances: %v", err)
			//}
			//<-time.After(5 * time.Second)
			//var id string
			if !s.orderSent {

				s.id = ulid.Make().String()
				s.logger.Println("Sending TradingMessageType_OrderMessage NewOrder")
				if err := stream.Send(&pb.TradingMessage{
					Type: pb.TradingMessageType_AddOrderRequest,
					Data: &pb.TradingMessage_AddOrder{AddOrder: &pb.AddOrder{
						OrderRequestID: s.id,
						PistonID:       s.id,
						Side:           pb.OrderSide_BID,
						Instrument:     "BTC/USDT",
						Type:           pb.OrderType_LIMIT,
						Size:           0.00001,
						Price:          10000.3,
						Created:        time.Now().UnixNano(),
					}},
				}); err != nil {
					s.logger.Errorf("Failed to send TradingMessage_Order New: %v", err)
				}
				s.orderSent = true
			}

			<-time.After(5 * time.Second)
			//s.logger.Println("Sending TradingMessageType_RequestOpenOrders")
			//if err := stream.Send(&pb.TradingMessage{
			//	Type: pb.TradingMessageType_RequestOpenOrders,
			//	Data: &pb.TradingMessage_OpenOrdersRequest{OpenOrdersRequest: &pb.OpenOrdersRequest{
			//		RequestID: 100,
			//	}},
			//}); err != nil {
			//	s.logger.Fatalf("Failed to send RequestOpenOrders: %v", err)
			//}
			//<-time.After(5 * time.Second)

			//s.logger.Infof("Sending TradingMessageType_OrderMessage Cancel exchangeID %s", s.exchangeID)
			//if err := stream.Send(&pb.TradingMessage{
			//	Type: pb.TradingMessageType_CancelOrderRequest,
			//	Data: &pb.TradingMessage_CancelOrder{CancelOrder: &pb.CancelOrder{
			//		OrderRequestID: s.id,
			//		PistonID:    s.id,
			//		ExchangeID:     s.exchangeID,
			//		Instrument:     "BTC/USDT",
			//	}},
			//}); err != nil {
			//	s.logger.Errorf("Failed to send TradingMessage_Order Cancel: %v", err)
			//}

			//<-time.After(time.Second * 2)

			s.logger.Println("Sending TradingMessageType_MoveOrderRequest")
			if err := stream.Send(&pb.TradingMessage{
				Type: pb.TradingMessageType_MoveOrderRequest,
				Data: &pb.TradingMessage_MoveOrder{MoveOrder: &pb.MoveOrder{
					OrderRequestID:        s.id,
					PistonID:              s.id,
					ExchangeID:            s.exchangeID,
					ExchangeClientOrderID: "",
					Exchange:              "",
					Instrument:            "BTC/USDT",
					Side:                  pb.OrderSide_BID,
					Type:                  pb.OrderType_LIMIT,
					Size:                  0.00001,
					Price:                 10000,
					Created:               time.Now().UnixNano(),
				}},
			}); err != nil {
				s.logger.Fatalf("Failed to send TradingMessage_Order Cancel: %v", err)
			}
			<-time.After(2 * time.Second)

			//s.logger.Println("Sending TradingMessageType_MoveOrderRequest")
			//if err := stream.Send(&pb.TradingMessage{
			//	Type: pb.TradingMessageType_MoveOrderRequest,
			//	Data: &pb.TradingMessage_MoveOrder{MoveOrder: &pb.MoveOrder{
			//		OrderRequestID:        id,
			//		PistonID:           id,
			//		ExchangeID:            "",
			//		ClientOrderID: "",
			//		Exchange:              "OKX",
			//		Instrument:            "BTC/USDT",
			//		Side:                  pb.OrderSide_BID,
			//		Type:                  pb.OrderType_LIMIT,
			//		Size:                  0.0005,
			//		Price:                 10000,
			//		Created:               time.Now().UnixNano(),
			//	}},
			//}); err != nil {
			//	s.logger.Fatalf("Failed to send TradingMessage_Order Cancel: %v", err)
			//}
			//
			//<-time.After(time.Second * 10)

			// Send Market Order
			//s.logger.Println("Sending TradingMessageType_OrderMessage NewOrder")
			//if err := stream.Send(&pb.TradingMessage{
			//	Type: pb.TradingMessageType_AddOrderRequest,
			//	Data: &pb.TradingMessage_AddOrder{AddOrder: &pb.AddOrder{
			//		OrderRequestID: "proximity_client_order_id",
			//		PistonID:    "proximityID",
			//		Side:           pb.OrderSide_ASK,
			//		Instrument:     "BTC/USDT",
			//		Type:           pb.OrderType_MARKET,
			//		Size:           0.003,
			//	}},
			//}); err != nil {
			//	s.logger.Errorf("Failed to send TradingMessage_Order New: %v", err)
			//}

		}
	}()

	for {
		in, err := stream.Recv()
		if err != nil {
			return err
		}
		s.logger.Printf("Received message: %v", in)

		switch in.Type {
		case pb.TradingMessageType_OrderAddedResponse:
			orderAdded := in.GetOrderAdded()
			s.exchangeID = orderAdded.ExchangeID
		case pb.TradingMessageType_OrderMovedResponse:
			orderMoved := in.GetOrderMoved()
			s.exchangeID = orderMoved.ExchangeID
		}

	}
}

func main() {

	var cfg = logster.Config{
		Project:           "gateways",
		Format:            "text",
		Level:             "debug",
		Env:               "local",
		DisableStackTrace: true,
		System:            "",
		Inst:              "",
	}

	var logger = logster.New(os.Stdout, cfg)

	lis, err := net.Listen("tcp", ":5800")
	if err != nil {
		logger.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterPistonServer(s, &server{logger: logger, orderSent: false})
	logger.Println("Server is running on port 5800")
	if err = s.Serve(lis); err != nil {
		logger.Fatalf("failed to serve: %v", err)
	}
}
