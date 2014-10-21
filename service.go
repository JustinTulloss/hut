package hut

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

type Service struct {
	Router *mux.Router
	Env    Env
}

// This serves as a health check and a default route just to make sure we're
// up and running where we expect. The implementation is expected to fill out
// the rest of the routes somewhere other than "/"
func (s *Service) defaultHandler(w http.ResponseWriter, r *http.Request) {
	memStats := &runtime.MemStats{}
	runtime.ReadMemStats(memStats)
	rep := map[string]interface{}{
		"status":     "running",
		"memStats":   memStats,
		"cpus":       runtime.NumCPU(),
		"goroutines": runtime.NumGoroutine(),
		"cgocalls":   runtime.NumCgoCall(),
	}
	s.Reply(rep, w)
}

// Optionally you can pass in a router to the service in order to instantiate
// it somewhere down the routing chaing. By default, the service assumes it's
// the root.
func NewService(r *mux.Router) *Service {
	if r == nil {
		r = mux.NewRouter()
	}
	// Default to OsEnv until we have more ways of doing configuration
	service := &Service{Router: r, Env: NewOsEnv()}
	r.HandleFunc("/health", service.defaultHandler)
	return service
}

// Starts this service running on a port specified in the env.
// Note that it's not necessary to call start if this service is going to be
// served by a different service.
func (s *Service) Start() {
	http.Handle("/", handlers.CombinedLoggingHandler(os.Stdout, s.Router))

	port, err := s.Env.GetString("port")
	if err != nil {
		log.Fatal("Could not get port to start on: ", err)
	}
	// Actually start the server
	port = fmt.Sprintf(":%s", port)
	log.Printf("Server started on %s\n", port)
	log.Fatal(http.ListenAndServe(port, nil))
}
