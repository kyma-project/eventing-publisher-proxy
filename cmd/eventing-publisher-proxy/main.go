package main //nolint:cyclop // it is only starting required instances.

import (
	log "log"

	"github.com/kelseyhightower/envconfig"
	"github.com/kyma-project/eventing-publisher-proxy/pkg/commander"
	"github.com/kyma-project/eventing-publisher-proxy/pkg/commander/eventmesh"
	"github.com/kyma-project/eventing-publisher-proxy/pkg/commander/nats"
	"github.com/kyma-project/eventing-publisher-proxy/pkg/metrics"
	"github.com/kyma-project/eventing-publisher-proxy/pkg/metrics/latency"
	"github.com/kyma-project/eventing-publisher-proxy/pkg/options"
	"github.com/prometheus/client_golang/prometheus"

	emlogger "github.com/kyma-project/eventing-manager/pkg/logger"
)

const (
	backendEventMesh = "beb"
	backendNATS      = "nats"
)

type Config struct {
	// Backend used for Eventing. It could be "nats" or "beb".
	Backend string `envconfig:"BACKEND" required:"true"`

	// AppLogFormat defines the log format.
	AppLogFormat string `default:"json" envconfig:"APP_LOG_FORMAT"`

	// AppLogLevel defines the log level.
	AppLogLevel string `default:"info" envconfig:"APP_LOG_LEVEL"`
}

func main() {
	opts := options.New()
	if err := opts.Parse(); err != nil {
		log.Fatalf("Failed to parse options, error: %v", err)
	}

	// parse the config for main:
	cfg := new(Config)
	if err := envconfig.Process("", cfg); err != nil {
		log.Fatalf("Failed to read configuration, error: %v", err)
	}

	// init the logger
	logger, err := emlogger.New(cfg.AppLogFormat, cfg.AppLogLevel)
	if err != nil {
		log.Fatalf("Failed to initialize logger, error: %v", err)
	}
	defer func() {
		if err := logger.WithContext().Sync(); err != nil {
			log.Printf("Failed to flush logger, error: %v", err)
		}
	}()
	setupLogger := logger.WithContext().With("backend", cfg.Backend)

	// metrics collector
	metricsCollector := metrics.NewCollector(latency.NewBucketsProvider())
	prometheus.MustRegister(metricsCollector)
	metricsCollector.SetHealthStatus(true)

	// Instantiate configured commander.
	var c commander.Commander
	switch cfg.Backend {
	case backendEventMesh:
		c = eventmesh.NewCommander(opts, metricsCollector, logger)
	case backendNATS:
		c = nats.NewCommander(opts, metricsCollector, logger)
	default:
		setupLogger.Fatalf("Invalid publisher backend: %v", cfg.Backend)
	}

	// Init the commander.
	if err := c.Init(); err != nil {
		setupLogger.Fatalw("Commander initialization failed", "error", err)
	}

	// Start the metrics server.
	metricsServer := metrics.NewServer(logger)
	defer metricsServer.Stop()
	if err := metricsServer.Start(opts.MetricsAddress); err != nil {
		setupLogger.Infow("Failed to start metrics server", "error", err)
	}

	setupLogger.Infof("Starting publisher to: %v", cfg.Backend)

	// Start the commander.
	if err := c.Start(); err != nil {
		setupLogger.Fatalw("Failed to start publisher", "error", err)
	}

	setupLogger.Info("Shutdown the Event Publisher")
}
