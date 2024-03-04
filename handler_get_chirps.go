package main

import (
	"net/http"
	"sort"
	"strconv"

	"github.com/go-chi/chi/v5"
)

func (cfg *apiConfig) handlerGetChirps(w http.ResponseWriter, req *http.Request) {
	dbChirps, err := cfg.DB.GetChirps()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't retrieve chirps")
		return
	}

	authorID := req.URL.Query().Get("author_id")

	chirps := []Chirp{}
	if authorID != "" {
		authorIDInt, err := strconv.Atoi(authorID)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Couldn't parse author ID")
			return
		}

		for _, dbChirp := range dbChirps {
			if dbChirp.Author_ID == authorIDInt {
				chirps = append(chirps, Chirp{
					ID:        dbChirp.ID,
					Body:      dbChirp.Body,
					Author_ID: dbChirp.Author_ID,
				})
			}
		}
	} else {
		for _, dbChirp := range dbChirps {
			chirps = append(chirps, Chirp{
				ID:        dbChirp.ID,
				Body:      dbChirp.Body,
				Author_ID: dbChirp.Author_ID,
			})
		}
	}

	sort.Slice(chirps, func(i, j int) bool {
		return chirps[i].ID < chirps[j].ID
	})

	respondWithJSON(w, http.StatusOK, chirps)
}

func (cfg *apiConfig) handlerGetChirpByID(w http.ResponseWriter, req *http.Request) {
	params := chi.URLParam(req, "chirpID")
	id, err := strconv.Atoi(params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Invalid chirp ID")
		return
	}

	dbChirp, err := cfg.DB.GetChirp(id)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Couldn't get chirp")
		return
	}

	respondWithJSON(w, http.StatusOK, Chirp{
		ID:        dbChirp.ID,
		Body:      dbChirp.Body,
		Author_ID: dbChirp.Author_ID,
	})
}
