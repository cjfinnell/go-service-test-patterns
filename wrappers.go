package main

import (
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

const authHeader = "Authorization"

func authWrap(next httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		authValue := r.Header.Get(authHeader)

		if authValue != "" {
			log.Print("auth successful")
			next(w, r, ps)

			return
		}

		log.Print("auth failed")

		w.WriteHeader(http.StatusUnauthorized)

		if _, err := w.Write([]byte("unauthorized")); err != nil {
			log.Printf("failed to write 'unauthorized' response: %s", err)
		}
	}
}

func requestLogger(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s", r.Method, r.URL)
		next.ServeHTTP(w, r)
	}
}
