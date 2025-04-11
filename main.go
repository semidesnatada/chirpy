package main

import (
	"context"
	"fmt"
	"net/http"

	_ "github.com/lib/pq"
	"github.com/semidesnatada/chirpy/internal/database"
)

func main() {

	s := database.CreateState()

	user, uErr := s.DB.CreateUser(context.Background(), "chicken@beans.com")
	if uErr == nil {
		fmt.Println(user)
	} else {
		fmt.Println(uErr.Error())
	}

	apiCfg := apiConfig{}
	serveHandler := http.NewServeMux()

	//app namespace
	serveHandler.Handle("/app/", apiCfg.middlewareMetricsInc(http.StripPrefix("/app/", http.FileServer(http.Dir("./serve/")))))
	
	//api namespace
	serveHandler.HandleFunc("GET /api/healthz", healthzHandler)
	serveHandler.HandleFunc("POST /api/validate_chirp", chirpValidationHandler)

	//admin namespace
	serveHandler.HandleFunc("GET /admin/metrics", apiCfg.metricsHandler)
	serveHandler.HandleFunc("POST /admin/reset", apiCfg.resetMetricsHandler)


	toast := http.Server{
		Handler: serveHandler,
		Addr: ":8080",
	}

	toast.ListenAndServe()

}
