package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/semidesnatada/chirpy/internal/auth"
	"github.com/semidesnatada/chirpy/internal/database"
)


func (cfg *apiConfig) usersCreationHandler(w http.ResponseWriter, req *http.Request) {
	
	type requestValues struct {
		Email string `json:"email"`
		Password string `json:"password"`
	}
	type responseValues struct {
		Id uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Email string `json:"email"`
		IsChirpyRed bool `json:"is_chirpy_red"`
	}

	decoder := json.NewDecoder(req.Body)
	reqParams := requestValues{}
	decodeErr := decoder.Decode(&reqParams)

	if decodeErr != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't decode the request", decodeErr)
		return
	}

	HashedPassword, hErr := auth.HashPassword(reqParams.Password)
	if hErr != nil {
		respondWithError(w, http.StatusBadRequest, "this password is no bueno", hErr)
		return
	}

	u, uErr := cfg.DB.CreateUser(req.Context(), database.CreateUserParams{
		Email: reqParams.Email,
		HashedPassword: HashedPassword,
	})

	if uErr != nil {
		respondWithError(w, http.StatusBadRequest, "couldn't create a user with this email and password", uErr)
		return
	}

	respondWithJSON(w, http.StatusCreated, responseValues{
		Id: u.ID,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
		Email: u.Email,
		IsChirpyRed: u.IsChirpyRed,
	})
}

func (cfg *apiConfig) usersUpdateHandler(w http.ResponseWriter, req *http.Request) {

	accessToken, err := auth.GetBearerToken(req.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "not able to extract authorization token", err)
		return
	}

	userId, authErr := auth.ValidateJWT(accessToken, cfg.JWTSecret)
	if authErr != nil {
		respondWithError(w, http.StatusUnauthorized, "not authorized", authErr)
		return
	}

	type requestValues struct {
		Email string `json:"email"`
		Password string `json:"password"`
	}
	type responseValues struct {
		Id uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Email string `json:"email"`
		IsChirpyRed bool `json:"is_chirpy_red"`
	}

	decoder := json.NewDecoder(req.Body)
	reqParams := requestValues{}
	decodeErr := decoder.Decode(&reqParams)

	if decodeErr != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't decode the request", decodeErr)
		return
	}

	HashedPassword, hErr := auth.HashPassword(reqParams.Password)
	if hErr != nil {
		respondWithError(w, http.StatusBadRequest, "this password is no bueno", hErr)
		return
	}

	u, uErr := cfg.DB.UpdateUser(req.Context(), database.UpdateUserParams{
		ID: userId,
		Email: reqParams.Email,
		HashedPassword: HashedPassword,
	})

	if uErr != nil {
		respondWithError(w, http.StatusBadRequest, "couldn't update the user", uErr)
		return
	}

	respondWithJSON(w, http.StatusOK, responseValues{
		Id: u.ID,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
		Email: u.Email,
		IsChirpyRed: u.IsChirpyRed,
	})
}