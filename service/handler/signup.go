package handler

import (
	"github.com/cbotte21/games-auth/datastore"
	"github.com/cbotte21/games-auth/schema"
	"github.com/cbotte21/games-auth/utilities"
	"net/http"
	"strconv"
	"time"
)

func SignupHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm() //Populate PostForm
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Please try again later.\n"))
		return
	}

	credentials := r.PostForm

	if !credentials.Has("email") || !credentials.Has("username") || !credentials.Has("password") { //HAS email and password
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Request must contain an email, username, and password.\n"))
		return
	}

	if !utilities.ParseEmail(credentials.Get("email")) || !utilities.ParseUsername(credentials.Get("username")) || !utilities.ParsePassword(credentials.Get("password")) { //Validate username and password
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Email, username, or password does not meet requirements.\n"))
		return
	}

	//Call signup logic
	//if USER EXISTS
	//if created successfully
	//else internal error

	emailCheckQuery := schema.User{
		Email: credentials.Get("email"),
	}

	usernameCheckQuery := schema.User{
		Username: credentials.Get("username"),
	}

	//Check if an account already exists with email or username
	_, result := datastore.Find(emailCheckQuery)
	if result {
		w.WriteHeader(http.StatusConflict)
		w.Write([]byte("Email is already registered.\n"))
		return
	}
	_, result = datastore.Find(usernameCheckQuery)
	if result {
		w.WriteHeader(http.StatusConflict)
		w.Write([]byte("Username is already registered.\n"))
		return
	}

	candideUser := schema.User{
		Email:            credentials.Get("email"),
		Username:         credentials.Get("username"),
		Password:         credentials.Get("password"),
		InitialTimestamp: strconv.FormatInt(time.Now().Unix(), 10),
	}

	status := datastore.Create(candideUser)

	if !status {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Please try again later. (Could not create account)\n"))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("account created"))
}
