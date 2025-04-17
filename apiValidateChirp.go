package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func chirpValidationHandler(w http.ResponseWriter, req *http.Request) {

	// checks whether a chirp is of the correct length

	const maxChirpLength = 140

	type requestValues struct {
		Body string `json:"body"`
	}
	type responseValues struct {
		// Valid bool `json:"valid,omitempty"`
		Error string `json:"error,omitempty"`
		CleanedBody string `json:"cleaned_body,omitempty"`
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

	respondWithJSON(w, http.StatusOK, responseValues{
		CleanedBody: cleanedBody,
	})
	
}

// func getCleanBody(body string) string {
// 	words := strings.Split(body, " ")
// 	for i, word := range words {
// 		if lowered := strings.ToLower(word) ; lowered == "kerfuffle" || lowered == "sharbert" || lowered == "fornax" {
// 			words[i] = "****"
// 		}
// 	}
// 	cleanedBody := strings.Join(words, " ")
// 	return cleanedBody
// }
