package handler

import (
	"github.com/cbotte21/auth-go/internal/validator"
	"github.com/cbotte21/microservice-common/pkg/common_errors"
	"github.com/cbotte21/microservice-common/pkg/datastore"
	"github.com/cbotte21/microservice-common/pkg/schema"
	"github.com/cbotte21/microservice-common/pkg/validate"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
	"strconv"
	"time"
)

func SignupHandler(w http.ResponseWriter, r *http.Request, userClient *datastore.MongoClient[schema.User]) {
	err := r.ParseForm() //Populate PostForm
	if err != nil {
		common_errors.InternalServiceError(&w)
		return
	}

	payload := r.PostForm

	errMsg, validReq := validate.ValidateRequestWithErrorMessage(payload, "email", "password")
	if !validReq {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write(errMsg)
		return
	}

	if !validator.ValidateEmail(payload.Get("email")) || !validator.ValidatePassword(payload.Get("password")) { //Validate username and password
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte("Email or password does not meet requirements.\n"))
		return
	}

	//Check if an account already exists with email or username
	emailCheckQuery := schema.User{
		Email: payload.Get("email"),
	}
	_, err = userClient.Find(emailCheckQuery)
	if err != mongo.ErrNoDocuments {
		w.WriteHeader(http.StatusConflict)
		_, _ = w.Write([]byte("Email is already registered.\n"))
		return
	}

	//Register
	currTime := strconv.FormatInt(time.Now().Unix(), 10)
	candideUser := schema.User{
		Email:            payload.Get("email"),
		InitialTimestamp: currTime,
		RecentTimestamp:  currTime,
		Role:             0,
	}

	if candideUser.SetPassword(payload.Get("password")) != nil {
		common_errors.InternalServiceError(&w)
		return
	}

	err = userClient.Create(candideUser)

	if err != nil {
		common_errors.InternalServiceError(&w)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(`{ "status": "account created" }`))
}
