package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/semidesnatada/chirpy/internal/database"
)

func (cfg *apiConfig) chirpsCreationHandler(w http.ResponseWriter, req *http.Request) {

	// checks whether a chirp is of the correct length

	const maxChirpLength = 140

	type requestValues struct {
		Body string `json:"body"`
		UserID uuid.UUID `json:"user_id"`
	}
	type responseValues struct {
		Id uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Body string `json:"body"`
		UserId uuid.UUID `json:"user_id"`
	}

	decoder := json.NewDecoder(req.Body)
	reqParams := requestValues{}
	decodeErr := decoder.Decode(&reqParams)

	if decodeErr != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't decode the request", decodeErr)
		return
	} else if length := len(reqParams.Body) ; length > maxChirpLength {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("chirp is too long: %d", length), nil)
		return
	}

	cleanedBody := getCleanBody(reqParams.Body)

	chirp, creErr := cfg.DB.CreateChirp(
		req.Context(),
		database.CreateChirpParams{
			Body: cleanedBody,
			UserID: reqParams.UserID,
		},
	)
	if creErr != nil {
		respondWithError(w, http.StatusBadRequest, "couldn't create this chirp", creErr)
		return
	}

	respondWithJSON(w, http.StatusCreated, responseValues{
		Id: chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body: chirp.Body,
		UserId: chirp.UserID,
	})
	
}

func getCleanBody(body string) string {
	words := strings.Split(body, " ")
	for i, word := range words {
		if lowered := strings.ToLower(word) ; lowered == "kerfuffle" || lowered == "sharbert" || lowered == "fornax" {
			words[i] = "****"
		}
	}
	cleanedBody := strings.Join(words, " ")
	return cleanedBody
}


func (cfg *apiConfig) chirpsGetHandler(w http.ResponseWriter, req *http.Request) {

	type responseItem struct {
		Id uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Body string `json:"body"`
		UserId uuid.UUID `json:"user_id"`
	}
	type responseValues []responseItem

	chirps, cErr := cfg.DB.GetAllChirps(req.Context())
	if cErr != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't get all the chirps for you", cErr)
		return
	}

	responsePayload := make(responseValues, len(chirps))

	for i, chirp := range chirps {
		responsePayload[i] = responseItem{
			Id: chirp.ID,
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
			Body: chirp.Body,
			UserId: chirp.UserID,
		}
	}

	respondWithJSON(w, http.StatusOK, responsePayload)
}

func (cfg *apiConfig) chirpsGetSingleHandler(w http.ResponseWriter, req *http.Request) {

	type responseItem struct {
		Id uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Body string `json:"body"`
		UserId uuid.UUID `json:"user_id"`
	}

	chirpIDstring := req.PathValue("chirpID")
	if chirpIDstring == "" {
		respondWithError(w, http.StatusBadRequest, "that's not a valid chirp ID", errors.New("couldn't identify chirp ID from url path"))
		return
	}

	chirpID, parsErr := uuid.Parse(chirpIDstring)
	if parsErr != nil {
		respondWithError(w, http.StatusBadRequest, "that's not a valid chirp ID", errors.New("couldn't parse the chirp ID"))
		return
	}

	chirp, cErr := cfg.DB.GetSingleChirps(req.Context(), chirpID)
	if cErr != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't get that chirp from the db", cErr)
		return
	}

	respondWithJSON(w, http.StatusOK, 
		responseItem{
			Id: chirp.ID,
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
			Body: chirp.Body,
			UserId: chirp.UserID,
		},
	)
}