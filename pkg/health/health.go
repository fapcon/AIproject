package health

import (
	"net/http"
	"sync"

	"github.com/go-chi/chi/v5"
)

//nolint:gochecknoglobals // local package variables
var (
	status = http.StatusServiceUnavailable
	mu     sync.RWMutex
)

func Status() int {
	mu.RLock()
	defer mu.RUnlock()
	return status
}

func SetStatus(s int) {
	mu.Lock()
	defer mu.Unlock()
	status = s
}

func Routes() chi.Router {
	r := chi.NewRouter()
	r.Get("/", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(Status())
	})
	return r
}
