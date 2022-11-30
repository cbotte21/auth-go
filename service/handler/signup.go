package handler

import (
	"github.com/cbotte21/auth-go/datastore"
	"github.com/cbotte21/auth-go/schema"
	"github.com/cbotte21/auth-go/utilities"
	"go.mongodb.org/mongo-driver/mongo"
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
	_, err = datastore.Find(emailCheckQuery)
	if err != mongo.ErrNoDocuments {
		w.WriteHeader(http.StatusConflict)
		w.Write([]byte("Email is already registered.\n"))
		return
	}
	_, err = datastore.Find(usernameCheckQuery)
	if err != mongo.ErrNoDocuments {
		w.WriteHeader(http.StatusConflict)
		w.Write([]byte("Username is already registered.\n"))
		return
	}

	candideUser := schema.User{
		Email:            credentials.Get("email"),
		Username:         credentials.Get("username"),
		InitialTimestamp: strconv.FormatInt(time.Now().Unix(), 10),
	}

	if candideUser.SetPassword(credentials.Get("password")) != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Please try again later.\n"))
		return
	}

	err = datastore.Create(candideUser)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Please try again later.\n"))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{ "status": "account created" }`))
}
