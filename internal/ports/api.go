package ports

import (
	"context"

	"google.golang.org/protobuf/types/known/emptypb"
	"studentgit.kata.academy/quant/torque/internal/app"
	protoapi "studentgit.kata.academy/quant/torque/internal/ports/types/api"
	"studentgit.kata.academy/quant/torque/pkg/logster"
)

type APIHandler struct {
	logger logster.Logger
	api    *app.API

	protoapi.UnimplementedAPIServiceServer
}

func NewAPIHandler(logger logster.Logger, api *app.API) *APIHandler {
	return &APIHandler{
		logger: logger,
		api:    api,

		UnimplementedAPIServiceServer: protoapi.UnimplementedAPIServiceServer{},
	}
}

func (h *APIHandler) SubscribeInstruments(
	ctx context.Context, request *protoapi.Instruments,
) (*emptypb.Empty, error) {
	h.api.SubscribeInstrumentsRequest(ctx, request.ToDomain())
	return new(emptypb.Empty), nil
}

func (h *APIHandler) UnsubscribeInstruments(
	ctx context.Context, request *protoapi.Instruments,
) (*emptypb.Empty, error) {
	return new(emptypb.Empty), h.api.UnsubscribeInstruments(ctx, request.ToDomain())
}

func (h *APIHandler) ModifyNameResolvers(
	ctx context.Context, request *protoapi.Instruments,
) (*emptypb.Empty, error) {
	return new(emptypb.Empty), h.api.AddInstruments(ctx, request.ToDomain())
}

func (h *APIHandler) ModifyInstrumentDetails(
	ctx context.Context, request *protoapi.InstrumentDetails,
) (*emptypb.Empty, error) {
	return new(emptypb.Empty), h.api.AddInstrumentDetails(ctx, request.ToDomain())
}
