package oauth

import (
	"github.com/kyma-project/eventing-publisher-proxy/pkg/env"
	"golang.org/x/oauth2/clientcredentials"
)

// Config returns a new oauth2 client credentials config instance.
func Config(cfg *env.EventMeshConfig) clientcredentials.Config {
	return clientcredentials.Config{
		ClientID:     cfg.ClientID,
		ClientSecret: cfg.ClientSecret,
		TokenURL:     cfg.TokenEndpoint,
	}
}
