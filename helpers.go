package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"rbrband/rbrband_feedback/creds"
	"strconv"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

//JSON encoder
func JSON(str string) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		io.WriteString(w, str)
	}
}

//CheckIfLoggedIn function to be called after every handler to check if FB login is there
func CheckIfLoggedIn(session *mgo.Session, musicianID string) (User, bool) {
	result := User{}
	c := session.DB(creds.MONGO_DB).C(creds.MONGO_MUSICIAN_COLLECTION)
	c.Find(bson.M{"id": musicianID, "isLoggedIn": true}).One(&result)
	if !result.IsLoggedIn {
		log.Println("User not logged in")
		return result, false
	}
	return result, true
}

func FloatToTwo(number float64) float64 {
	temp := fmt.Sprintf("%.2f", number)
	finalFloat, _ := strconv.ParseFloat(temp, 64)
	return finalFloat
}
