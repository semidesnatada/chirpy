package main

import (
	"net/http"

	"github.com/semidesnatada/chirpy/internal/auth"
)

func (cfg *apiConfig) revokeHandler(w http.ResponseWriter, req *http.Request) {

	refreshToken, err := auth.GetBearerToken(req.Header)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "can't process this request properly", err)
		return
	}

	revokeErr := cfg.DB.RevokeRefreshToken(req.Context(), refreshToken)
	if revokeErr != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't revoke this toke", revokeErr)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
	// respondWithJSON(w, , struct{}{})
}