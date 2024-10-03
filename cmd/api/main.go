package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"net"
	"net/http"

	"studentgit.kata.academy/quant/torque/config"

	pkggrpc "studentgit.kata.academy/quant/torque/pkg/grpc"

	"os"

	"studentgit.kata.academy/quant/torque/internal/adapters/kafka"
	"studentgit.kata.academy/quant/torque/internal/ports"
	"studentgit.kata.academy/quant/torque/internal/ports/types/api"

	"studentgit.kata.academy/quant/torque/internal/app"
	"studentgit.kata.academy/quant/torque/pkg/health"
	"studentgit.kata.academy/quant/torque/pkg/postgres"
	"studentgit.kata.academy/quant/torque/pkg/sig"

	"github.com/oklog/ulid/v2"
	"golang.org/x/sync/errgroup"
	"studentgit.kata.academy/quant/torque/internal/storage"
	"studentgit.kata.academy/quant/torque/pkg/clocks"
	pkghttp "studentgit.kata.academy/quant/torque/pkg/http"
	"studentgit.kata.academy/quant/torque/pkg/logster"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
)

func main() {
	var apiConfig config.Config
	var apiConfigFile string

	flag.StringVar(&apiConfigFile, "apiConfig", "config/api/api_local.yaml", "Path to the apiConfig file")
	flag.Parse()
	err := config.LoadConfig(apiConfigFile, &apiConfig)
	if err != nil {
		panic(err)
	}

	logger := logster.New(os.Stdout, apiConfig.Log)
	defer func() { _ = logger.Sync() }()

	g, ctx := errgroup.WithContext(context.Background())

	defer logster.LogShutdownDuration(ctx, logger)()
	g.Go(func() error {
		return sig.ListenSignal(ctx, logger)
	})

	logger.Infof("service starting with apiConfig %+v", apiConfig)

	// metrics.
	registry := prometheus.NewRegistry()
	registerer := prometheus.WrapRegistererWith(prometheus.Labels{"go_project": apiConfig.Log.Project}, registry)
	registerer.MustRegister(
		collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}),
		collectors.NewGoCollector(),
	)
	// postgres
	postgresDB, err := postgres.New(logger, apiConfig.Postgres)
	logster.InitResource(logger, postgresDB, err)

	clock := clocks.RealClock{}
	ulidGenerator := func() string { return ulid.Make().String() }

	// kafka
	consumers := kafka.NewRealConsumersBuilder(logger, kafka.ConsumerConfigs{ //nolint:exhaustruct // we don't consume
		// group internal events
		InternalEvents: apiConfig.Kafka.Consumers.Internal,
	})
	producers := kafka.NewRealProducersBuilder(logger, kafka.ProducerConfigs{
		InternalEvents: apiConfig.Kafka.Producers.Internal,
	})

	serviceStorage, err := storage.NewBuilder(
		logger, postgresDB, consumers, producers, clock, ulidGenerator, registerer, g, apiConfig.Log.Project,
	).
		WithSubscribeInstrumentsRequestProducer().
		Build(ctx)

	logster.InitResource(logger, serviceStorage, err)
	defer logster.CloseResource(logger, serviceStorage)

	g.Go(func() error {
		return logster.LogIfError(logger, serviceStorage.Run(ctx), "Storage error")
	})

	techHandler := pkghttp.NewHandler("/", pkghttp.DefaultTechOptions(logger, registry))

	g.Go(func() error {
		return logster.LogIfError(
			logger, pkghttp.RunServer(ctx, apiConfig.PrivateAddr, logger, techHandler),
			"Tech server error",
		)
	})

	apiService := app.NewAPI(logger, serviceStorage)

	g.Go(func() error {
		return logster.LogIfError(logger, RunGrpcServer(ctx, logger, apiService, &apiConfig), "Grpc server error")
	})

	health.SetStatus(http.StatusOK)

	logger.Infof("Waiting for signal")
	err = g.Wait()
	logger.Infof("error: %v", err)
	if err != nil && !errors.Is(err, sig.ErrSignalReceived) {
		logger.WithError(err).Errorf("Exit reason")
	}
}

func RunGrpcServer(ctx context.Context, logger logster.Logger, apiService *app.API, cfg *config.Config) error {
	logger.Infof("Starting grpc server")
	grpcServer := pkggrpc.NewGRPCServer(logger)
	grpcHandler := ports.NewAPIHandler(logger, apiService)
	api.RegisterAPIServiceServer(grpcServer, grpcHandler)

	listen, err := net.Listen("tcp", cfg.GRPCApiAddr)
	if err != nil {
		return fmt.Errorf("GRPC server can't listen requests")
	}

	errListen := make(chan error, 1)
	go func() {
		errListen <- grpcServer.Serve(listen)
	}()

	select {
	case <-ctx.Done():
		grpcServer.GracefulStop()
		return nil
	case err = <-errListen:
		return fmt.Errorf("can't run grpc server: %w", err)
	}
}
