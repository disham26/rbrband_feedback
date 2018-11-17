package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	mgo "gopkg.in/mgo.v2"
)

//GuestFeedback is the first function for the landing page
func GuestFeedback(session *mgo.Session) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		session := mongoSession.Copy()
		defer session.Close()
		Vars := mux.Vars(r)
		message := ""
		bandID := Vars["band_id"]
		if r.Method == "POST" {
			r.ParseForm()
			r.Form.Get("band_rating")
			log.Println("Band Rating is :", r.Form.Get("band_rating"))
			band := GetBand(session, bandID)
			err := UpdateBandRating(session, band, r.Form.Get("band_rating"))
			log.Println("Updating musician rating")
			err = UpdateMusiciansRating(session, band.ID, r.Form.Get("band_rating"))
			log.Println("Error to add rating is :", err)
			if err == nil {
				message = "Thank you for your review, this helps a lot to the musicians"
			}
		}

		band := GetBand(session, bandID)
		t, _ := ParseTemplate("goadmin")
		page := &FeedbackPage{
			Title:       "",
			Text:        "Rubberband home",
			StaticHost:  getStaticHost(),
			JSON:        "{}",
			Config:      "index",
			LoggedIn:    false,
			CurrentPage: "feedback",
			Band:        band,
			Message:     message,
		}
		t.ExecuteTemplate(w, "feedback", page)
		return
	}
}
