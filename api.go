package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/semidesnatada/chirpy/user"
)


func healthzHandler(w http.ResponseWriter, req *http.Request) {

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(200)

	w.Write([]byte(http.StatusText(200)))
}


func chirpValidationHandler(w http.ResponseWriter, req *http.Request) {

	// checks whether a chirp is of the correct length

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

	responseBytes := responseValues{}

	if decodeErr != nil {
		log.Printf("Error decoding request: %s", decodeErr)
		w.WriteHeader(500)
		responseBytes.Error = "Something went wrong"
	} else if length := len(reqParams.Body) ; length > 140 {
		log.Printf("Chirp is too long: %d", length)
		w.WriteHeader(400)
		responseBytes.Error = "Chirp is too long"
	} else {
		log.Printf("Chirp is valid")
		w.WriteHeader(200)
		// responseBytes.Valid = true
		
		words := strings.Split(reqParams.Body, " ")
		for i, word := range words {
			if lowered := strings.ToLower(word) ; lowered == "kerfuffle" || lowered == "sharbert" || lowered == "fornax" {
				words[i] = "****"
			}
		}
		cleanedBody := strings.Join(words, " ")

		responseBytes.CleanedBody = cleanedBody
	}

	data, encodeErr := json.Marshal(responseBytes)

	if encodeErr != nil {
		log.Printf("Error encoding response: %s", encodeErr)
		w.WriteHeader(500)
		// responseBytes.Error = "Something went wrong"
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
	
}

func newUserHandler(w http.ResponseWriter, req *http.Request) {
	U := user.NewUser("slowery9", "slowery9@gmail.com", "beansontoast69")
	fmt.Println(U.Email)
	w.WriteHeader(500)
}