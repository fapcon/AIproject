package postgres

import (
	"context"
	"fmt"
	"net"
	"net/url"
	"time"

	"github.com/jackc/pgx/v4"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4/pgxpool"
	"go.uber.org/multierr"
	"studentgit.kata.academy/quant/torque/pkg/logster"
)

const timeout = 2 * time.Second

type logAdapter interface {
	Log(ctx context.Context, level pgx.LogLevel, msg string, data map[string]interface{})
}

type loggerAdapter struct {
	logger logster.Logger
}

func (l *loggerAdapter) Log(_ context.Context, level pgx.LogLevel, msg string, data map[string]interface{}) {
	switch level {
	case pgx.LogLevelTrace:
		l.logger.WithFields(data).Debugf(msg)
	case pgx.LogLevelDebug:
		l.logger.WithFields(data).Debugf(msg)
	case pgx.LogLevelInfo:
		l.logger.WithFields(data).Debugf(msg)
	case pgx.LogLevelWarn:
		l.logger.WithFields(data).Warnf(msg)
	case pgx.LogLevelError:
		l.logger.WithFields(data).Errorf(msg)
	case pgx.LogLevelNone:
		l.logger.WithFields(data).Errorf(msg)
	}
}

type DB struct {
	pool   *pgxpool.Pool
	logger logster.Logger
}

type Config struct {
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	Database string `yaml:"database"`
	Params   string `yaml:"params"`
	Region   string `yaml:"region"`

	MinConnections    int           `yaml:"minConnections"`
	MaxConnections    int           `yaml:"maxConnections"`
	HealthCheckPeriod time.Duration `yaml:"healthcheckPeriod"`
	MaxConnIdleTime   time.Duration `yaml:"maxConnIdletime"`
	MaxConnLifetime   time.Duration `yaml:"maxConnLifetime"`
}

func (c Config) DSN() string {
	return fmt.Sprintf("postgres://%s:%s@%s/%s%s", c.User, c.Password,
		net.JoinHostPort(c.Host, c.Port), c.Database, c.Params)
}

func (c Config) DSNNEW() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s",
		c.Host, c.Port, c.User, c.Password, c.Database)
}

func (c Config) String() string {
	return fmt.Sprintf(
		"{DSN:%s MaxConnections:%d MinConnections:%d MaxConnLifetime:%s MaxConnIdleTime:%s HealthCheckPeriod:%s}",
		maskDSN(c.DSN()), c.MaxConnections, c.MinConnections, c.MaxConnLifetime, c.MaxConnIdleTime, c.HealthCheckPeriod,
	)
}

func (c Config) GoString() string {
	return c.String()
}

func New(logger logster.Logger, cfg Config) (*DB, error) {
	poolConfig, err := pgxpool.ParseConfig(cfg.DSNNEW())
	if err != nil {
		return nil, fmt.Errorf("unable to parse postgres config: %w", err)
	}

	poolConfig.MaxConns = int32(cfg.MaxConnections)
	poolConfig.MinConns = int32(cfg.MinConnections)
	poolConfig.MaxConnLifetime = cfg.MaxConnLifetime
	poolConfig.MaxConnIdleTime = cfg.MaxConnIdleTime
	poolConfig.HealthCheckPeriod = cfg.HealthCheckPeriod
	if l, ok := logger.(logAdapter); ok {
		poolConfig.ConnConfig.Logger = l
	} else {
		poolConfig.ConnConfig.Logger = &loggerAdapter{logger: logger.WithPrefix(logster.LibPrefix)}
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	pool, err := pgxpool.ConnectConfig(ctx, poolConfig)
	if err != nil {
		return nil, fmt.Errorf("can't open connection to postgres: %w", err)
	}

	db := &DB{pool: pool, logger: logger}

	go db.checkLatency(cfg.HealthCheckPeriod)

	return db, nil
}

func (d *DB) Close() error {
	d.pool.Close()
	return nil
}

func (d *DB) Stats() *pgxpool.Stat {
	return d.Pool().Stat()
}

func (d *DB) Pool() *pgxpool.Pool {
	return d.pool
}

func (d *DB) Logger() logster.Logger {
	return d.logger
}

func (d *DB) validateQuery(query string) error {
	conn, err := d.pool.Acquire(context.Background())
	if err != nil {
		return fmt.Errorf("can't acquire connection: %w", err)
	}
	defer conn.Release()

	const warnThreshold = time.Second
	stmtName := uuid.NewString()
	now := time.Now()

	if _, err = conn.Conn().Prepare(context.Background(), stmtName, query); err != nil {
		return fmt.Errorf("can't create statement from query %q: %w", query, err)
	}
	if err = conn.Conn().Deallocate(context.Background(), stmtName); err != nil {
		return fmt.Errorf("can't close statement from query %q: %w", query, err)
	}

	if time.Since(now) > warnThreshold {
		d.logger.WithField("query", query).Warnf("Prepare time > %s", warnThreshold)
	}

	return nil
}

func (d *DB) ValidateQueries(queries []string) error {
	var err error
	for _, query := range queries {
		if errX := d.validateQuery(query); errX != nil {
			err = multierr.Append(err, errX)
		}
	}
	if err != nil {
		return err
	}
	d.logger.Infof("All queries are validated")
	return nil
}

func (d *DB) checkLatency(period time.Duration) {
	ticker := time.NewTicker(period)
	logger := d.logger.WithField("service", "check_latency_task")
	for range ticker.C {
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		var f int
		err := d.Pool().QueryRow(ctx, "SELECT 1").Scan(&f)
		if err != nil {
			logger.WithError(err).Errorf("can't check postgres latency")
		}
		cancel()
	}
}

func maskDSN(dsn string) string {
	parsed, err := url.Parse(dsn)
	if err != nil {
		return dsn
	}
	parsed.User = url.UserPassword("****", "****")

	// cause url.String() return escaped * symbol, it shouldn't return err
	unescaped, err := url.QueryUnescape(parsed.String())
	if err != nil {
		return dsn
	}

	return unescaped
}
