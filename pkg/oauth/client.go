package oauth

import (
	"context"
	"net/http"

	"github.com/kyma-project/eventing-publisher-proxy/pkg/env"
	"github.com/kyma-project/eventing-publisher-proxy/pkg/tracing/propagation/tracecontextb3"
	"go.opencensus.io/plugin/ochttp"
	"golang.org/x/oauth2"
)

// NewClient returns a new HTTP client which have nested transports for handling oauth2 security,
// HTTP connection pooling, and tracing.
func NewClient(ctx context.Context, cfg *env.EventMeshConfig) *http.Client {
	// configure auth client
	config := Config(cfg)
	client := config.Client(ctx)

	// configure connection transport
	var base = http.DefaultTransport.(*http.Transport).Clone()
	cfg.ConfigureTransport(base)
	client.Transport.(*oauth2.Transport).Base = base

	// configure tracing transport
	client.Transport = &ochttp.Transport{
		Base:        client.Transport,
		Propagation: tracecontextb3.TraceContextEgress(),
	}

	return client
}
