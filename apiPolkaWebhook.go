package main

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/google/uuid"
	"github.com/semidesnatada/chirpy/internal/auth"
)

func (cfg *apiConfig) polkaWebhookHandler(w http.ResponseWriter, req *http.Request) {

	type requestValues struct {
		Event string `json:"event"`
		Data struct {
			UserID uuid.UUID `json:"user_id"`
		} `json:"data"`
	}

	apiKey, aPerr := auth.GetAPIKey(req.Header)
	if aPerr != nil {
		respondWithError(w, http.StatusUnauthorized, "can't detect an api key in this request", aPerr)
		return
	}
	if apiKey != cfg.PolkaSecret {
		respondWithError(w, http.StatusUnauthorized, "this is not a valid api key", errors.New("you are not authorized to make this request"))
		return
	}

	decoder := json.NewDecoder(req.Body)
	reqParams := requestValues{}
	decodeErr := decoder.Decode(&reqParams)

	if decodeErr != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't decode the request", decodeErr)
		return
	}

	if reqParams.Event != "user.upgraded" {
		respondWithError(w, http.StatusNoContent, "only care about user upgrades", errors.New("event is not 'user.upgraded'"))
		return
	}

	_, updateErr := cfg.DB.UpgradeUser(req.Context(), reqParams.Data.UserID)

	if updateErr != nil {
		respondWithError(w, http.StatusNotFound, "couldn't update user status in DB", updateErr)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)


}