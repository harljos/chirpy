package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/harljos/chirpy/internal/auth"
)

func (cfg *apiConfig) handlerLoginUser(w http.ResponseWriter, req *http.Request) {
	type parameters struct {
		Email            string `json:"email"`
		Password         string `json:"password"`
		ExpiresInSeconds int    `json:"expires_in_seconds"`
	}
	type response struct {
		User
		Token string `json:"token"`
	}

	decoder := json.NewDecoder(req.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
		return
	}

	user, err := cfg.DB.LoginUser(params.Email, params.Password)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Email or Password don't match")
		return
	}

	const twentyFourHrs = 86400

	if params.ExpiresInSeconds == 0 {
		params.ExpiresInSeconds = twentyFourHrs
	} else if params.ExpiresInSeconds > twentyFourHrs {
		params.ExpiresInSeconds = twentyFourHrs
	}

	token, err := auth.MakeJWT(user.ID, cfg.jwtSecret, time.Duration(params.ExpiresInSeconds)*time.Second)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't make JWT")
		return
	}

	respondWithJSON(w, http.StatusOK, response{
		User: User{
			Email: user.Email,
			ID:    user.ID,
		},
		Token: token,
	})
}
