package main

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/harljos/chirpy/internal/auth"
	"github.com/harljos/chirpy/internal/database"
)

func (cfg *apiConfig) handlerUpgrageUserWebhook(w http.ResponseWriter, req *http.Request) {
	type parameters struct {
		Event string `json:"event"`
		Data  struct {
			UserId int `json:"user_id"`
		}
	}

	apiKey, err := auth.GetApiKey(req.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't find api key")
		return
	}

	if apiKey != cfg.polkaKey {
		respondWithError(w, http.StatusUnauthorized, "Api key is invalid")
		return
	}

	decoder := json.NewDecoder(req.Body)
	params := parameters{}
	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
		return
	}

	if params.Event != "user.upgraded" {
		respondWithJSON(w, http.StatusOK, struct{}{})
		return
	}

	err = cfg.DB.UpgradeUser(params.Data.UserId)
	if err != nil {
		if errors.Is(err, database.ErrNotExist) {
			respondWithError(w, http.StatusNotFound, "Couldn't find user")
			return
		}
		respondWithError(w, http.StatusInternalServerError, "Couldn't update user")
		return
	}

	respondWithJSON(w, http.StatusOK, struct{}{})
}
