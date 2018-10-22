package main

import (
	"log"
	"os"
	"rbrband/rbrband_feedback/creds"
)

//getEnv function will be responsible to set keys and passwords stored in creds
func getEnv() string {
	env := os.Getenv("RUBBERBAND_ENV_NAME")
	if env == "" {
		return "local"
	}
	return env
}

//getStaticHost will define the URL of JS, CSS and images for every environment
func getStaticHost() string {
	env := getEnv()
	if env == "local" {

		return "http://localhost:5011"
	}
	return "production static link"

}

//getDBEnv will set the URL, id and password for mongo connection
func getDBEnv(mongo Mongo) Mongo {
	switch getEnv() {
	case "local":
		log.Println("local")
		mongo.MongoURL = creds.MONGO_LOCAL_URL
		mongo.MongoUserID = creds.MONGO_LOCAL_USERNAME
		mongo.MongoPassword = creds.MONGO_LOCAL_PASSWORD
		break
	case "staging":
		mongo.MongoURL = creds.MONGO_STAGING_URL
		mongo.MongoUserID = creds.MONGO_STAGING_USERNAME
		mongo.MongoPassword = creds.MONGO_STAGING_PASSWORD
		break
	case "prod":
		mongo.MongoURL = creds.MONGO_PROD_URL
		mongo.MongoUserID = creds.MONGO_PROD_USERNAME
		mongo.MongoPassword = creds.MONGO_PROD_PASSWORD
		break
	}
	return mongo

}
