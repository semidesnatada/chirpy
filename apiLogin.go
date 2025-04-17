package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/semidesnatada/chirpy/internal/auth"
)

func (cfg *apiConfig) loginHandler(w http.ResponseWriter, req *http.Request) {
	
	type requestValues struct {
		Email string `json:"email"`
		Password string `json:"password"`
	}
	type responseValues struct {
		Id uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Email string `json:"email"`
	}

	decoder := json.NewDecoder(req.Body)
	reqParams := requestValues{}
	decodeErr := decoder.Decode(&reqParams)

	if decodeErr != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't decode the request", decodeErr)
		return
	}

	// hashedPass, hErr := cfg.DB.GetHashedPassword(req.Context(), reqParams.Email)
	// if hErr != nil {
	// 	respondWithError(w, http.StatusUnauthorized, "Incorrect email or password", hErr)
	// 	return
	// }

	user, uErr := cfg.DB.GetUser(req.Context(), reqParams.Email)
	if uErr != nil {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password", uErr)
		return
	}

	correctPass := auth.CheckPasswordHash(user.HashedPassword, reqParams.Password)
	if correctPass != nil {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password", correctPass)
		return
	}

	respondWithJSON(w, http.StatusOK, 
		responseValues{
			Id: user.ID,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
			Email: user.Email,
		},
	)

}