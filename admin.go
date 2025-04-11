package main

import (
	"fmt"
	"net/http"
	"sync/atomic"
)


type apiConfig struct {
	fileserverHits atomic.Int32
}

func (cfg *apiConfig) metricsHandler(w http.ResponseWriter, req *http.Request) {

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(200)

	// w.Write([]byte(fmt.Sprintf("Hits: %d",cfg.fileserverHits.Load())))
	page := []byte(fmt.Sprintf(
	`<html>
		<body>
			<h1>Welcome, Chirpy Admin</h1>
			<p>Chirpy has been visited %d times!</p>
		</body>
	</html>`,
	cfg.fileserverHits.Load()))
	w.Write(page)
}

func (cfg *apiConfig) resetMetricsHandler(w http.ResponseWriter, req *http.Request) {

	cfg.fileserverHits = atomic.Int32{}
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(200)

	w.Write([]byte(fmt.Sprintf("Hits: %d",cfg.fileserverHits.Load())))
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, req * http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, req)
	})

}