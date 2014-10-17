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

func NewService() *Service {
	r := mux.NewRouter()
	// Default to OsEnv until we have more ways of doing configuration
	service := &Service{Router: r, Env: NewOsEnv()}
	r.HandleFunc("/", service.defaultHandler)
	return service
}

func (s *Service) Start() {
	http.Handle("/", handlers.CombinedLoggingHandler(os.Stdout, s.Router))

	// Actually start the server
	port := fmt.Sprintf(":%s", os.Getenv("PORT"))
	log.Printf("Server started on %s\n", port)
	log.Fatal(http.ListenAndServe(port, nil))
}
