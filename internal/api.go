package service

import (
	"github.com/cbotte21/auth-go/internal/handler"
	"github.com/cbotte21/microservice-common/pkg/datastore"
	"github.com/cbotte21/microservice-common/pkg/jwtParser"
	"github.com/cbotte21/microservice-common/pkg/schema"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

type Api struct {
	port       int
	router     *mux.Router
	userClient *datastore.MongoClient[schema.User]
	jwtSecret  *jwtParser.JwtSecret
}

func NewApi(port int, userClient *datastore.MongoClient[schema.User], jwtSecret *jwtParser.JwtSecret) (*Api, bool) {
	api := &Api{}
	api.port = port
	api.router = mux.NewRouter()
	api.userClient = userClient
	api.jwtSecret = jwtSecret
	api.RegisterHandlers()
	return api, true
}

func (api *Api) Start() error {
	return http.ListenAndServe(":"+strconv.Itoa(api.port), api.router)
}

func (api *Api) RegisterHandlers() { //Add all API handlers here
	api.router.HandleFunc("/", handler.IndexHandler).Methods("GET")
	//User lifecycle
	api.router.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		handler.LoginHandler(w, r, api.userClient, api.jwtSecret)
	}).Methods("POST")
	api.router.HandleFunc("/signup", func(w http.ResponseWriter, r *http.Request) {
		handler.SignupHandler(w, r, api.userClient)
	}).Methods("POST")
	api.router.HandleFunc("/verify", func(w http.ResponseWriter, r *http.Request) {
		handler.VerifyHandler(w, r, api.jwtSecret)
	}).Methods("POST")
}
