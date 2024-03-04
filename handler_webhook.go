package main

import (
	"encoding/json"
	"net/http"
)

func (cfg *apiConfig) handlerUpgrageUserWebhook(w http.ResponseWriter, req *http.Request) {
	type parameters struct {
		Event string `json:"event"`
		Data  struct {
			UserId int `json:"user_id"`
		} 
	}

	decoder := json.NewDecoder(req.Body)
	params := parameters{}
	err := decoder.Decode(&params)
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
		respondWithError(w, http.StatusNotFound, "User Not Found")
		return
	}

	respondWithJSON(w, http.StatusOK, struct{}{})
}
