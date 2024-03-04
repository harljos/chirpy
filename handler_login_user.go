package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/harljos/chirpy/internal/auth"
)

func (cfg *apiConfig) handlerLoginUser(w http.ResponseWriter, req *http.Request) {
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	type response struct {
		User
		Token        string `json:"token"`
		RefreshToken string `json:"refresh_token"`
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

	accessToken, err := auth.MakeJWT(user.ID, cfg.jwtSecret, time.Hour, auth.TokenTypeAccess)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create access JWT")
		return
	}
	refreshToken, err := auth.MakeJWT(user.ID, cfg.jwtSecret, time.Duration(1440)*time.Hour, auth.TokenTypeRefresh)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create refresh JWT")
		return
	}

	respondWithJSON(w, http.StatusOK, response{
		User: User{
			Email:       user.Email,
			ID:          user.ID,
			IsChirpyRed: user.IsChirpyRed,
		},
		Token:        accessToken,
		RefreshToken: refreshToken,
	})
}
