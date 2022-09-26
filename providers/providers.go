package providers

import (
	"net/http"
)

// Provider is a abstracted transit provider
type Provider interface {
	GetNextTrip(route, stop, direction string) (int64, error)
}

// Providers provides the ability to get a commmunication interface to the
// appropriate provider, based on the passed config
type Providers interface {
	GetProvider(providerName string) Provider
}

// DefaultProviders is the default set of providers to use. This can be exchanged for testing.
type DefaultProviders struct{ Sandboxed bool }

// GetProvider will return the appropriate provider communicator based on the config
func (providers *DefaultProviders) GetProvider(providerName string) Provider {
	switch providerName {
	case "metrotransit":
		return &MetroTransitProvider{
			UseSandbox: providers.Sandboxed,
			APIClient:  &http.Client{},
		}
	default:
		return nil
	}
}
