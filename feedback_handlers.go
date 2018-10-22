package main

import (
	"net/http"

	mgo "gopkg.in/mgo.v2"
)

//GuestFeedback is the first function for the landing page
func GuestFeedback(session *mgo.Session) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		session := mongoSession.Copy()
		defer session.Close()

		t, _ := ParseTemplate("goadmin")
		page := &Page{
			Title:       "",
			Text:        "Rubberband home",
			StaticHost:  getStaticHost(),
			JSON:        "{}",
			Email:       "prasang@rbrband.in",
			Config:      "index",
			LoggedIn:    false,
			CurrentPage: "feedback",
		}
		t.ExecuteTemplate(w, "feedback", page)
		return
	}
}
