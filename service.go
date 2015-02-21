package hut

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime"
	"strings"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

type Service struct {
	name string
	// Use this to route HTTP requests to this service.
	Router *mux.Router
	// Represents the current environment the service is running in. The
	// environment can be filled out in a number of ways, but services can
	// always find the variables here.
	Env Env
	// A log helper with different levels of logging.
	Log Log
	// Record stats about your service here. Defaults to having a prefix
	// that is the same as the service's binary name.
	Stats Statter
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
	var service *Service
	if r == nil {
		r = mux.NewRouter()
		r.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			service.NotFound(w)
		})
	}
	fullPath := strings.Split(os.Args[0], "/")
	name := fullPath[len(fullPath)-1]
	// Default to OsEnv until we have more ways of doing configuration
	service = &Service{
		name:   name,
		Router: r,
		Env:    NewOsEnv(),
		Log:    NewStdLog(),
	}
	service.Stats = service.NewStatsd()
	r.HandleFunc("/health", service.defaultHandler)
	return service
}

// Starts this service running on a port specified in the env.
// Note that it's not necessary to call start if this service is going to be
// served by a different service.
func (s *Service) Start() {
	http.Handle("/", handlers.CombinedLoggingHandler(os.Stdout, s.Router))

	port := s.Env.GetString("port")
	// Actually start the server
	port = fmt.Sprintf(":%s", port)
	log.Printf("Server started on %s\n", port)
	log.Fatal(http.ListenAndServe(port, nil))
}
