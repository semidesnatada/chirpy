package main

import (
	"errors"
	"net/http"

	"github.com/semidesnatada/chirpy/internal/auth"
)

func (cfg *apiConfig) refreshHandler(w http.ResponseWriter, req *http.Request) {

	type responseValues struct {
		Token string `json:"token"`
	}

	refreshToken, err := auth.GetBearerToken(req.Header)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "can't process this request properly", err)
		return
	}

	tokenExists, dBerr := cfg.DB.CheckRefreshTokenExistsAndIsValid(req.Context(), refreshToken)

	if dBerr != nil {
		respondWithError(w, http.StatusInternalServerError, "problem with checking this token", dBerr)
		return
	}

	if !tokenExists {
		respondWithError(w, http.StatusUnauthorized, "couldn't find this token in the DB", errors.New("not authorised"))
		return
	}

	userID, uErr := cfg.DB.GetUserFromAccessToken(req.Context(), refreshToken)
	if uErr != nil {
		respondWithError(w, http.StatusInternalServerError, "problem with checking which user this is", uErr)
		return
	}
	
	//need to get JWT auth token too to include in response
	newAccessToken, aErr := auth.MakeJWT(userID, cfg.JWTSecret)
	if aErr != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't create JWT auth", aErr)
	}

	respondWithJSON(w, http.StatusOK, responseValues{
		Token: newAccessToken,
	})

}