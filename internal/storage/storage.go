package storage

import (
	"context"
	"fmt"
	"time"

	"studentgit.kata.academy/quant/torque/internal/adapters/kafka"

	grpcclients "studentgit.kata.academy/quant/torque/internal/adapters/grpc_clients"

	websocketclients "studentgit.kata.academy/quant/torque/internal/adapters/websocket_clients"

	"studentgit.kata.academy/quant/torque/internal/adapters/cache"
	"studentgit.kata.academy/quant/torque/internal/adapters/clients"
	"studentgit.kata.academy/quant/torque/internal/adapters/postgres"
	"studentgit.kata.academy/quant/torque/internal/adapters/prometheus"
	"studentgit.kata.academy/quant/torque/pkg/clocks"

	prom "github.com/prometheus/client_golang/prometheus"
	"golang.org/x/sync/errgroup"
	pkgpostgres "studentgit.kata.academy/quant/torque/pkg/postgres"

	"studentgit.kata.academy/quant/torque/pkg/logster"
)

type Builder struct {
	builds []func(*Storage) error
	inits  []func(*Storage, context.Context) error
}

func NewBuilder(
	logger logster.Logger,
	db *pkgpostgres.DB,
	consumers kafka.ConsumersBuilder,
	producers kafka.ProducersBuilder,
	clock clocks.Clock,
	ulidGenerator func() string,
	registerer prom.Registerer,
	eg *errgroup.Group,
	subsystem string,
) *Builder {
	return &Builder{
		builds: []func(*Storage) error{
			func(s *Storage) error {
				collector, err := prometheus.NewCollector(logger, registerer, subsystem)
				if err != nil {
					return err
				}
				postgresDB, err := postgres.NewDatabase(db)
				if err != nil {
					return err
				}
				*s = Storage{
					logger:           logger,
					clock:            clock,
					ulidGenerator:    ulidGenerator,
					postgres:         postgresDB,
					prometheus:       collector,
					kafka:            kafka.NewKafka(logger, consumers, producers),
					cache:            cache.NewInMemoryStorage(),
					clients:          clients.NewClients(logger),
					grpcClients:      grpcclients.NewGrpcClients(logger),
					websocketClients: websocketclients.NewWebsocketClients(logger),
					capabilities:     make(map[Capability]struct{}),
					eg:               eg,
					runs:             nil,
				}
				return nil
			},
		},
		inits: []func(*Storage, context.Context) error{
			func(s *Storage, ctx context.Context) error {
				return s.websocketClients.Init(ctx)
			},
			func(s *Storage, ctx context.Context) error {
				return s.kafka.Init(ctx)
			},
		},
	}
}

func (b *Builder) Build(ctx context.Context) (*Storage, error) {
	s := new(Storage)
	for _, f := range b.builds {
		err := f(s)
		if err != nil {
			return nil, err
		}
	}

	for _, f := range b.inits {
		err := f(s, ctx)
		if err != nil {
			return nil, err
		}
	}
	return s, nil
}

func (b *Builder) buildError(f func(s *Storage) error) {
	b.builds = append(b.builds, f)
}

func (b *Builder) build(f func(s *Storage)) {
	b.buildError(func(s *Storage) error {
		f(s)
		return nil
	})
}

func (s *Storage) Close() error {
	s.logger.Infof("storage closed")
	return s.kafka.Close()
}

func (b *Builder) buildCap(c Capability) {
	b.buildError(func(s *Storage) error {
		if _, ok := s.capabilities[c]; ok {
			return fmt.Errorf("capability %s already activated", c)
		}
		s.capabilities[c] = struct{}{}
		return nil
	})
}

func (b *Builder) buildKafka(f func(k *kafka.Kafka) error) {
	b.buildError(func(s *Storage) error {
		return f(s.kafka)
	})
}

func (b *Builder) init(f func(s *Storage, ctx context.Context) error) {
	b.inits = append(b.inits, func(s *Storage, ctx context.Context) error {
		return f(s, ctx)
	})
}

func (b *Builder) runCron(interval time.Duration, f func(s *Storage, ctx context.Context)) {
	b.runCtx(func(s *Storage, ctx context.Context) {
		ticks, stop := s.clock.NewTicker(interval)
		defer stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticks:
				f(s, ctx)
			}
		}
	})
}

func (b *Builder) runCtx(f func(s *Storage, ctx context.Context)) {
	b.build(func(s *Storage) {
		s.runs = append(s.runs, func(ctx context.Context) error {
			f(s, ctx)
			return nil
		})
	})
}

type Storage struct {
	logger           logster.Logger
	clock            clocks.Clock
	ulidGenerator    func() string
	postgres         *postgres.Database
	kafka            *kafka.Kafka
	prometheus       *prometheus.Collector
	cache            *cache.InMemoryStorage
	clients          *clients.Clients
	websocketClients *websocketclients.Clients
	grpcClients      *grpcclients.Clients
	capabilities     map[Capability]struct{}
	runs             []func(context.Context) error
	eg               *errgroup.Group
}

func (s *Storage) Run(ctx context.Context) error {
	for _, f := range s.runs {
		f := f
		s.eg.Go(func() error {
			return f(ctx)
		})
	}
	s.eg.Go(func() error {
		return s.websocketClients.Run(ctx)
	})
	time.Sleep(time.Second) // Need to fully initialize gateway
	s.eg.Go(func() error {
		return s.grpcClients.Run(ctx)
	})
	s.eg.Go(func() error {
		return s.kafka.Run(ctx)
	})
	return nil
}

func (s *Storage) Cache() *cache.InMemoryStorage {
	return s.cache
}

//nolint:unused // Will need later
func (s *Storage) has(c Capability) bool {
	_, ok := s.capabilities[c]
	return ok
}

func (s *Storage) mustHave(capabilities ...Capability) {
	for _, c := range capabilities {
		if _, ok := s.capabilities[c]; !ok {
			panic(fmt.Sprintf("Storage doesn't have capability %s", c))
		}
	}
}
