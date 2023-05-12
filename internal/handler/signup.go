package handler

import (
	"github.com/cbotte21/auth-go/internal/validator"
	"github.com/cbotte21/microservice-common/pkg/datastore"
	"github.com/cbotte21/microservice-common/pkg/schema"
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

	if !credentials.Has("email") || !credentials.Has("password") { //HAS email and password
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Request must contain an email and password.\n"))
		return
	}

	if !validator.ValidateEmail(credentials.Get("email")) || !validator.ValidatePassword(credentials.Get("password")) { //Validate username and password
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Email or password does not meet requirements.\n"))
		return
	}

	//Check if an account already exists with email or username
	emailCheckQuery := schema.User{
		Email: credentials.Get("email"),
	}
	_, err = datastore.Find(emailCheckQuery)
	if err != mongo.ErrNoDocuments {
		w.WriteHeader(http.StatusConflict)
		w.Write([]byte("Email is already registered.\n"))
		return
	}

	//Register
	currTime := strconv.FormatInt(time.Now().Unix(), 10)
	candideUser := schema.User{
		Email:            credentials.Get("email"),
		InitialTimestamp: currTime,
		RecentTimestamp:  currTime,
		Role:             0,
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
