package nats

import (
	"fmt"

	"github.com/nats-io/nats.go"
)

var ErrNATSConnectionNotConnected = fmt.Errorf("NATS connection not connected")

type Opt = nats.Option

//nolint:gochecknoglobals // cloning functions as variables.
var (
	WithRetryOnFailedConnect = nats.RetryOnFailedConnect
	WithMaxReconnects        = nats.MaxReconnects
	WithReconnectWait        = nats.ReconnectWait
	WithName                 = nats.Name
)

// Connect returns a NATS connection that is ready for use, or an error if connection to the NATS server failed.
// It uses the nats.Connect function which is thread-safe.
func Connect(url string, opts ...Opt) (*nats.Conn, error) {
	connection, err := nats.Connect(url,
		opts...,
	)
	if err != nil {
		return nil, err
	}

	if status := connection.Status(); status != nats.CONNECTED {
		return nil, fmt.Errorf("%w with status:%v", ErrNATSConnectionNotConnected, status)
	}

	return connection, err
}
