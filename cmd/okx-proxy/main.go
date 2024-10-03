package main

import (
	"context"
	"errors"
	"flag"
	"net/http"
	"os"

	"github.com/oklog/ulid/v2"
	"golang.org/x/sync/errgroup"
	"studentgit.kata.academy/quant/torque/config"
	"studentgit.kata.academy/quant/torque/internal/adapters/kafka"
	"studentgit.kata.academy/quant/torque/internal/app"
	"studentgit.kata.academy/quant/torque/internal/storage"
	"studentgit.kata.academy/quant/torque/pkg/clocks"
	"studentgit.kata.academy/quant/torque/pkg/health"
	pkghttp "studentgit.kata.academy/quant/torque/pkg/http"
	"studentgit.kata.academy/quant/torque/pkg/logster"
	"studentgit.kata.academy/quant/torque/pkg/postgres"
	"studentgit.kata.academy/quant/torque/pkg/sig"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
)

func main() {
	var okxConfig config.Config
	var configFile string

	flag.StringVar(&configFile, "config", "config/okx-proxy/okx_local.yaml", "Path to the config file")
	flag.Parse()
	err := config.LoadConfig(configFile, &okxConfig)
	if err != nil {
		panic(err)
	}

	logger := logster.New(os.Stdout, okxConfig.Log)
	defer func() { _ = logger.Sync() }()

	g, ctx := errgroup.WithContext(context.Background())

	defer logster.LogShutdownDuration(ctx, logger)()
	g.Go(func() error {
		return sig.ListenSignal(ctx, logger)
	})

	logger.Infof("service starting with config %+v", okxConfig)

	// metrics.
	registry := prometheus.NewRegistry()
	registerer := prometheus.WrapRegistererWith(prometheus.Labels{"go_project": okxConfig.Log.Project}, registry)
	registerer.MustRegister(
		collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}),
		collectors.NewGoCollector(),
	)
	// postgres
	postgresDB, err := postgres.New(logger, okxConfig.Postgres)
	logster.InitResource(logger, postgresDB, err)

	clock := clocks.RealClock{}
	ulidGenerator := func() string { return ulid.Make().String() }

	// kafka
	consumers := kafka.NewRealConsumersBuilder(logger, kafka.ConsumerConfigs{
		InternalEvents:      okxConfig.Kafka.Consumers.Internal,
		GroupInternalEvents: okxConfig.Kafka.Consumers.GroupInternal,
	})
	producers := kafka.NewRealProducersBuilder(logger, kafka.ProducerConfigs{
		InternalEvents: okxConfig.Kafka.Producers.Internal,
	})

	serviceStorage, err := storage.NewBuilder(
		logger, postgresDB, consumers, producers, clock, ulidGenerator, registerer, g, okxConfig.Log.Project,
	).
		WithNewOrderRequestsConsumer().
		Build(ctx)

	logster.InitResource(logger, serviceStorage, err)
	defer logster.CloseResource(logger, serviceStorage)

	g.Go(func() error {
		return logster.LogIfError(logger, serviceStorage.Run(ctx), "Storage error")
	})

	okxService := app.NewOKXProxy(logger, serviceStorage)
	g.Go(func() error {
		return logster.LogIfError(logger, okxService.Run(ctx), "OKX service error")
	})

	techHandler := pkghttp.NewHandler("/", pkghttp.DefaultTechOptions(logger, registry))

	g.Go(func() error {
		return logster.LogIfError(
			logger, pkghttp.RunServer(ctx, okxConfig.PrivateAddr, logger, techHandler),
			"Tech server error",
		)
	})

	health.SetStatus(http.StatusOK)
	logger.Infof("Waiting for signal")
	err = g.Wait()
	logger.Infof("error: %v", err)
	if err != nil && !errors.Is(err, sig.ErrSignalReceived) {
		logger.WithError(err).Errorf("Exit reason")
	}
}