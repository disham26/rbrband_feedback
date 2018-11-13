package main

import (
	"log"
	"net/http"
	"os"
	"rbrband/rbrband_feedback/creds"

	"github.com/gorilla/mux"
	fb "github.com/huandu/facebook"
	mgo "gopkg.in/mgo.v2"
)

//Only mongo session
var mongoSession *mgo.Session

//Our main function
func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = creds.LOCAL_PORT
	}
	session := InitDB()
	fbApp := fb.New(creds.FACEBOOK_APP_ID, creds.FACEBOOK_SECRET_KEY)
	log.Println("app_access_token is:", fbApp.AppAccessToken())

	//Calling defer to close MongoDB connection
	defer mongoSession.Close()
	r := mux.NewRouter()
	api := r.PathPrefix("/api").Subrouter()

	//Redirections here
	r.HandleFunc("/", IndexHandler(session))
	r.HandleFunc("/signup", SignupHandler(session)).Methods("GET")
	r.HandleFunc("/profile/{id}", ProfileHandler(session))
	r.HandleFunc("/feedback/{id}", GuestFeedback(session))
	r.HandleFunc("/super", Index2Handler(session))
	r.HandleFunc("/band/{id}/qr", GetQRHandler(session))
	r.HandleFunc("/band/{id}", BandProfileHandler(session))

	//APIs here
	api.HandleFunc("/profileLogout", ProfileLogout(session)).Methods("POST")
	api.HandleFunc("/profileLogout/{id}", LogoutMusician(session))
	api.HandleFunc("/profileCheck/{id}", ProfileCheck(session))
	api.HandleFunc("/signup/{id}/{token}", ProfileCheck(session))
	api.HandleFunc("/bandsAssociated/{id}", GetBandsAssociated(session))

	//ListenAndServe
	http.ListenAndServe(":"+port, r)
}
