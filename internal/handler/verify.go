package handler

import (
	"github.com/cbotte21/microservice-common/pkg/common_errors"
	"github.com/cbotte21/microservice-common/pkg/jwtParser"
	"github.com/cbotte21/microservice-common/pkg/validate"
	"net/http"
)

func VerifyHandler(w http.ResponseWriter, r *http.Request, jwtSecret *jwtParser.JwtSecret) {
	err := r.ParseForm() //Populate PostForm
	if err != nil {
		common_errors.InternalServiceError(&w)
		return
	}

	payload := r.PostForm

	errMsg, validReq := validate.ValidateRequestWithErrorMessage(payload, "jwt")
	if !validReq {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write(errMsg)
		return
	}

	//Parse JWT
	err = jwtSecret.ValidateJWT(payload.Get("jwt"))

	if err != nil {
		common_errors.AccountNotAuthorizedError(&w)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{ "status": "account authorized" }`))
}
