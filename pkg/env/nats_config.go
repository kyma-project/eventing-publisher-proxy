package env

import (
	"fmt"
	"time"
)

// compile time check.
var _ fmt.Stringer = &NATSConfig{}

const JetStreamSubjectPrefix = "kyma"

// NATSConfig represents the environment config for the Event Publisher to NATS.
type NATSConfig struct {
	Port                  int           `default:"8080"       envconfig:"INGRESS_PORT"`
	URL                   string        `envconfig:"NATS_URL" required:"true"`
	RetryOnFailedConnect  bool          `default:"true"       envconfig:"RETRY_ON_FAILED_CONNECT"`
	MaxReconnects         int           `default:"-1"         envconfig:"MAX_RECONNECTS"` // Negative means keep try reconnecting.
	ReconnectWait         time.Duration `default:"5s"         envconfig:"RECONNECT_WAIT"`
	RequestTimeout        time.Duration `default:"5s"         envconfig:"REQUEST_TIMEOUT"`
	ApplicationCRDEnabled bool          `default:"true"       envconfig:"APPLICATION_CRD_ENABLED"`

	// Legacy Namespace is used as the event source for legacy events
	LegacyNamespace string `default:"kyma" envconfig:"LEGACY_NAMESPACE"`
	// EventTypePrefix is the prefix of each event as per the eventing specification
	// It follows the eventType format: <eventTypePrefix>.<appName>.<event-name>.<version>
	EventTypePrefix string `default:"kyma" envconfig:"EVENT_TYPE_PREFIX"`

	// JetStream-specific configs
	JSStreamName string `default:"kyma" envconfig:"JS_STREAM_NAME"`
}

// ToConfig converts to a default EventMeshConfig.
func (c *NATSConfig) ToConfig() *EventMeshConfig {
	cfg := &EventMeshConfig{
		EventMeshNamespace: c.LegacyNamespace,
		EventTypePrefix:    c.EventTypePrefix,
	}
	return cfg
}

// String implements the fmt.Stringer interface.
func (c *NATSConfig) String() string {
	return fmt.Sprintf("%#v", c)
}
