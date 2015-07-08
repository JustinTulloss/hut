package hut

import (
	"github.com/cactus/go-statsd-client/statsd"
)

// Statter an interface for collecting statistics. It mirrors statsd directly.
type Statter statsd.Statter

// Returns a new Statter object. Defaults to looking at the service environment
// for a statsd service, linked through docker link environment variables
// (`STATSD_PORT_8125_TCP_ADDR` and `STATSD_PORT_8125_TCP_PORT`)
// If those variables aren't set, returns a noop Statter that will still work
// but won't log your stats anywhere.
func (s *Service) NewStatsd() (client Statter) {
	address := s.Env.GetUDPServiceAddress("STATSD", 8125)
	if address == "" {
		client, _ = statsd.NewNoop()
		return client
	}
	client, err := statsd.New(address, s.name)
	if err != nil {
		s.Log.Error("Could not connect to statsd", "err", err)
		client, _ = statsd.NewNoop()
	}
	return
}
