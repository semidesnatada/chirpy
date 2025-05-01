package main

import (
	"net/http"

	_ "github.com/lib/pq"
)

func main() {

	apiCfg := createState()

	// user, uErr := apiCfg.DB.CreateUser(context.Background(), "chicken@beans.com")
	// if uErr == nil {
	// 	fmt.Println(user)
	// } else {
	// 	fmt.Println(uErr.Error())
	// }

	serveHandler := http.NewServeMux()

	//app namespace
	serveHandler.Handle("/app/", apiCfg.middlewareMetricsInc(http.StripPrefix("/app/", http.FileServer(http.Dir("./serve/")))))
	
	//api namespace
	serveHandler.HandleFunc("GET /api/healthz", healthzHandler)
	// serveHandler.HandleFunc("POST /api/validate_chirp", chirpValidationHandler)
	serveHandler.HandleFunc("POST /api/chirps", apiCfg.chirpsCreationHandler)
	serveHandler.HandleFunc("POST /api/users", apiCfg.usersCreationHandler)
	serveHandler.HandleFunc("GET /api/chirps", apiCfg.chirpsGetHandler)
	serveHandler.HandleFunc("GET /api/chirps/{chirpID}", apiCfg.chirpsGetSingleHandler)
	serveHandler.HandleFunc("POST /api/login", apiCfg.loginHandler)
	serveHandler.HandleFunc("POST /api/refresh", apiCfg.refreshHandler)
	serveHandler.HandleFunc("POST /api/revoke", apiCfg.revokeHandler)

	//admin namespace
	serveHandler.HandleFunc("GET /admin/metrics", apiCfg.metricsHandler)
	serveHandler.HandleFunc("POST /admin/reset", apiCfg.resetMetricsHandler)


	toast := http.Server{
		Handler: serveHandler,
		Addr: ":8080",
	}

	toast.ListenAndServe()

}
