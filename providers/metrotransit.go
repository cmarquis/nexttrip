package providers

import (
	"net/http"
)

// HTTPClient is the client used for contacting metrotransit this can be swapped
// for testing
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// MetroTransitProvider implements Provider for MetroTransit
type MetroTransit struct {
	APIClient  HTTPClient
	UseSandbox bool
}

func (p *MetroTransit) GetNextTrip(route, stop, direction string) (uint8, error) {
	return 0, nil
}
