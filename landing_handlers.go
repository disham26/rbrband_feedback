package main

import (
	"net/http"

	mgo "gopkg.in/mgo.v2"
)

//IndexHandler is the first function for the landing page
func IndexHandler(session *mgo.Session) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		session := mongoSession.Copy()
		defer session.Close()

		t, _ := ParseTemplate("goadmin")
		page := &Page{
			Title:       "Welcome to Rubberband",
			Text:        "Rubberband home",
			StaticHost:  getStaticHost(),
			JSON:        "{}",
			Email:       "prasang@rbrband.in",
			Config:      "index",
			LoggedIn:    false,
			CurrentPage: "index",
			UserID:      "",
		}
		t.ExecuteTemplate(w, "home", page)
		return
	}
}

//SignupHandler is called when Faceboon token and ID are fetched and inserted
func SignupHandler(session *mgo.Session) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		//t, _ := ParseTemplate("goadmin")
		fbToken := r.URL.Query().Get("fbToken")
		fbID := r.URL.Query().Get("id")
		id := fbID
		//Check if musician profile already exist
		if !CheckMusicianExistByID(session, fbID) {
			//So the profile does not exists, not we will create new
			id, _ = InsertMusicianFromFB(session, fbToken, fbID)
		} else {
			UpdateLoggedStatusByID(session, fbID, true)
		}

		http.Redirect(w, r, "/profile/"+id, http.StatusSeeOther)

		return
	}
}
