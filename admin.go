package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)


func (cfg *apiConfig) metricsHandler(w http.ResponseWriter, req *http.Request) {

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)

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

	godotenv.Load()
	PLATFORM := os.Getenv("PLATFORM")
	if PLATFORM != "dev" {
		// fmt.Println("DB_URL must be set")
		// os.Exit(1)
		respondWithError(w, http.StatusForbidden, "cannot access this endpoint in production", nil)
		return
	}

	cfg.DB.DeleteAllUsers(req.Context())

	cfg.fileserverHits.Store(0)
	
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	w.Write([]byte(fmt.Sprintf("Hits: %d",cfg.fileserverHits.Load())))
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, req * http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, req)
	})

}