package nats_test

import (
	"testing"
	"time"

	natsgo "github.com/nats-io/nats.go"

	"github.com/stretchr/testify/assert"

	eppnats "github.com/kyma-project/eventing-publisher-proxy/pkg/nats"
	epptestingutils "github.com/kyma-project/eventing-publisher-proxy/testing"
)

func TestConnect(t *testing.T) {
	testCases := []struct {
		name                      string
		givenRetryOnFailedConnect bool
		givenMaxReconnect         int
		givenReconnectWait        time.Duration
	}{
		{
			name:                      "do not retry failed connections",
			givenRetryOnFailedConnect: false,
			givenMaxReconnect:         0,
			givenReconnectWait:        time.Millisecond,
		},
		{
			name:                      "keep retrying failed connections",
			givenRetryOnFailedConnect: true,
			givenMaxReconnect:         -1,
			givenReconnectWait:        time.Millisecond,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			// given
			natsServer := epptestingutils.StartNATSServer()
			assert.NotNil(t, natsServer)
			defer natsServer.Shutdown()

			clientURL := natsServer.ClientURL()
			assert.NotEmpty(t, clientURL)

			// when
			connection, err := eppnats.Connect(clientURL,
				eppnats.WithRetryOnFailedConnect(tc.givenRetryOnFailedConnect),
				eppnats.WithMaxReconnects(tc.givenMaxReconnect),
				eppnats.WithReconnectWait(tc.givenReconnectWait),
			)
			assert.Nil(t, err)
			assert.NotNil(t, connection)
			defer func() { connection.Close() }()

			// then
			assert.Equal(t, connection.Status(), natsgo.CONNECTED)
			assert.Equal(t, clientURL, connection.Opts.Servers[0])
			assert.Equal(t, tc.givenRetryOnFailedConnect, connection.Opts.RetryOnFailedConnect)
			assert.Equal(t, tc.givenMaxReconnect, connection.Opts.MaxReconnect)
			assert.Equal(t, tc.givenReconnectWait, connection.Opts.ReconnectWait)
		})
	}
}
