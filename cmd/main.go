/*
*	Author: Cody Botte
*	Purpose: An authentication microservice using gorilla/mux. Users saved via mongodb, passwords encrypted using
* 		 bcrypt, and users identified using json web tokens.
 */

package main

import (
	"github.com/cbotte21/auth-go/internal"
	"github.com/cbotte21/microservice-common/pkg/datastore"
	"github.com/cbotte21/microservice-common/pkg/environment"
	"github.com/cbotte21/microservice-common/pkg/jwtParser"
	"github.com/cbotte21/microservice-common/pkg/schema"
	"log"
	"os"
	"strconv"
)

func main() {
	// Verify enviroment variables exist
	environment.VerifyEnvVariable("mongo_uri")
	environment.VerifyEnvVariable("mongo_db")
	environment.VerifyEnvVariable("port")
	environment.VerifyEnvVariable("jwt_secret")
	// Get port
	port, err := strconv.Atoi(environment.GetEnvVariable("port"))
	if err != nil {
		log.Fatalf("could not parse {port} enviroment variable")
	}
	// Initialize variables to attach to api
	jwtSecret := jwtParser.JwtSecret(os.Getenv("jwt_secret"))
	userClient := datastore.MongoClient[schema.User]{}
	err = userClient.Init()
	if err != nil {
		panic(err)
	}
	// Start API
	api, res := service.NewApi(port, &userClient, &jwtSecret)
	if !res || api.Start() != nil { //Start API Listener
		log.Fatal("Failed to initialize API.")
	}
}
