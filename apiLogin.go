package main

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/semidesnatada/chirpy/internal/auth"
	"github.com/semidesnatada/chirpy/internal/database"
)

func (cfg *apiConfig) loginHandler(w http.ResponseWriter, req *http.Request) {
	
	type requestValues struct {
		Email string `json:"email"`
		Password string `json:"password"`
		ExpiresInSeconds int `json:"expires_in_seconds"`
	}
	type responseValues struct {
		Id uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Email string `json:"email"`
		Token string `json:"token"`
		RefreshToken string `json:"refresh_token"`
		IsChirpyRed bool `json:"is_chirpy_red"`
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

	correctPass := auth.CheckPasswordHash(reqParams.Password, user.HashedPassword)
	if correctPass != nil {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password", correctPass)
		return
	}

	//need to get JWT auth token too to include in response
	token, aErr := auth.MakeJWT(user.ID, cfg.JWTSecret)
	if aErr != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't create JWT auth", aErr)
	}

	// need to get access token to also include in response
	refreshToken, _ := auth.MakeRefreshTokenString()

	// need to also store this new refresh token in the db.
	// if this is not possible, response with error

	_, rTerr := cfg.DB.CreateRefreshToken(
		req.Context(),
		database.CreateRefreshTokenParams{
			Token: refreshToken,
			UserID: user.ID,
			ExpiresAt: time.Now().Add(time.Hour * 24 * 60),
			RevokedAt: sql.NullTime{},
		},
	)

	if rTerr != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't create a refresh token", rTerr)
	}

	respondWithJSON(w, http.StatusOK, 
		responseValues{
			Id: user.ID,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
			Email: user.Email,
			Token: token,
			RefreshToken: refreshToken,
			IsChirpyRed: user.IsChirpyRed,
		},
	)

}