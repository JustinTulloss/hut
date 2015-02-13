package hut

import (
	"github.com/cactus/go-statsd-client/statsd"
)

type Statter statsd.Statter

func (s *Service) NewStatsd() (client Statter) {
	address := s.Env.GetUDPServiceAddress("STATSD", 8125)
	if address == "" {
		client, _ = statsd.NewNoop()
		return client
	}
	// TODO: put the right prefix in here
	client, err := statsd.New(address, s.name)
	if err != nil {
		s.Log.Error().Printf("Could not connect to statsd: %s\n", err)
		client, _ = statsd.NewNoop()
	}
	return
}
