package handler

import (
	"fmt"
	"github.com/cbotte21/microservice-common/pkg/common_errors"
	"github.com/cbotte21/microservice-common/pkg/datastore"
	"github.com/cbotte21/microservice-common/pkg/jwtParser"
	"github.com/cbotte21/microservice-common/pkg/schema"
	"github.com/cbotte21/microservice-common/pkg/validate"
	"net/http"
	"strconv"
	"time"
)

func LoginHandler(w http.ResponseWriter, r *http.Request, userClient *datastore.MongoClient[schema.User], jwtSecret *jwtParser.JwtSecret) { //TODO: Update last login
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

	query := schema.User{
		Email: payload.Get("email"),
	}

	candideUser, err := userClient.Find(query) //Find user with matching email
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte("Email does not exist.\n"))
		return
	}

	if candideUser.VerifyPassword(payload.Get("password")) != nil { //Validate password
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte("Username and password do not match."))
		return
	}

	updatedUser := candideUser
	updatedUser.RecentTimestamp = strconv.FormatInt(time.Now().Unix(), 10)
	_ = userClient.Update(candideUser, updatedUser)

	tokenString, err := jwtSecret.GenerateJWT(candideUser)

	if err != nil {
		common_errors.InternalServiceError(&w)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(fmt.Sprintf(`{ "jwt": "%s" }`, tokenString)))
}
