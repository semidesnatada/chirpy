package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"sort"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/semidesnatada/chirpy/internal/auth"
	"github.com/semidesnatada/chirpy/internal/database"
)

func (cfg *apiConfig) chirpsCreationHandler(w http.ResponseWriter, req *http.Request) {

	// checks whether a chirp is of the correct length

	const maxChirpLength = 140

	type requestValues struct {
		Body string `json:"body"`
		// UserID uuid.UUID `json:"user_id"`
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
	token, tErr := auth.GetBearerToken(req.Header)
	if tErr != nil {
		respondWithError(w, http.StatusUnauthorized, "could not authorise this chirp", tErr)
	}
	uId, pErr := auth.ValidateJWT(token, cfg.JWTSecret)
	if pErr != nil {
		respondWithError(w, http.StatusUnauthorized, "could not parse this authorisation token", pErr)
	}

	chirp, creErr := cfg.DB.CreateChirp(
		req.Context(),
		database.CreateChirpParams{
			Body: cleanedBody,
			UserID: uId,
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

	s := req.URL.Query().Get("author_id")

	var chirps []database.Chirp

	if s == "" {
		chirpsResult, cErr := cfg.DB.GetAllChirps(req.Context())
		if cErr != nil {
			respondWithError(w, http.StatusInternalServerError, "couldn't get all the chirps for you", cErr)
			return
		}
		chirps = chirpsResult
	} else {
		sID, parsErr := uuid.Parse(s)
		if parsErr != nil {
			respondWithError(w, http.StatusBadRequest, "couldn't parse an id from this search query", errors.New("can't process this request"))
			return
		}
		chirpsResult, cErr := cfg.DB.GetAllChirpsByAuthor(req.Context(), sID)
		if cErr != nil {
			respondWithError(w, http.StatusInternalServerError, "couldn't get all the chirps for you", cErr)
			return
		}
		chirps = chirpsResult
	}

	sortDirection := req.URL.Query().Get("sort")
	if sortDirection == "asc" {
		sort.Slice(chirps, func(i, j int) bool { return chirps[i].CreatedAt.Before(chirps[j].CreatedAt) })
	} else if sortDirection == "desc" {
		sort.Slice(chirps, func(i, j int) bool { return chirps[j].CreatedAt.Before(chirps[i].CreatedAt) })
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
		respondWithError(w, http.StatusNotFound, "couldn't get that chirp from the db", cErr)
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

func (cfg *apiConfig) chirpsDeleteSingleHandler(w http.ResponseWriter, req *http.Request) {

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

	ch, chErr := cfg.DB.GetSingleChirps(req.Context(), chirpID)
	if chErr != nil {
		respondWithError(w, http.StatusNotFound, "couldn't identify this chirp at all", chErr)
		return
	}
	if ch.UserID != userId {
		respondWithError(w, http.StatusForbidden, "you can't delete someone else's chirp", errors.New("no no not allowed"))
		return
	}

	cErr := cfg.DB.DeleteSingleChirps(req.Context(),
	database.DeleteSingleChirpsParams{
		ID: chirpID,
		UserID: userId,
	})

	if cErr != nil {
		respondWithError(w, http.StatusNotFound, "error deleting this chirp", cErr)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
}