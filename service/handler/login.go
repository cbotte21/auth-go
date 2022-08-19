package handler

import (
	"github.com/cbotte21/games-auth/datastore"
	"github.com/cbotte21/games-auth/schema"
	"net/http"
)

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm() //Populate PostForm
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Please try again later.\n"))
		return
	}

	credentials := r.PostForm

	if !credentials.Has("email") || !credentials.Has("password") { //HAS email and password
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Request must contain an email and password.\n"))
		return
	}

	query := schema.User{
		Email: credentials.Get("email"),
	}

	candideUser, result := datastore.Find(query)
	if !result {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Email does not exist.\n"))
		return
	}

	if !candideUser.VerifyPassword(credentials.Get("password")) { //Validate password
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Username and password do not match."))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("example jwt"))
}
