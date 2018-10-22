package main

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	mgo "gopkg.in/mgo.v2"
)

type ProfileResponse struct {
	Success bool `json:"success"`
}

//ProfileLogout is called when the profile of the musician is logged out
func ProfileLogout(session *mgo.Session) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		profileID := r.FormValue("fbId")
		profileExists := CheckMusicianExistByID(session, profileID)
		if profileExists {
			//User is found, now we can mark them as logged out
			UpdateLoggedStatusByID(session, profileID, false)
			resp := &BasicResponse{
				Success: true}
			Data, _ := json.Marshal(resp)
			JSON(string(Data))(w, r)
		} else {
			//Could not find id to logout
			JSON("{\"success\":false, \"error\":{\"code\":1004}}")(w, r)

		}

		return
	}
}

//ProfileCheck checks if a profile exists
func ProfileCheck(session *mgo.Session) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		Vars := mux.Vars(r)
		musicianID := Vars["id"]
		var profileCheckResponse ProfileResponse

		_, result := CheckIfLoggedIn(session, musicianID)
		if result {
			profileCheckResponse.Success = true
		} else {
			profileCheckResponse.Success = false
		}

		Data, _ := json.Marshal(profileCheckResponse)
		JSON(string(Data))(w, r)

		return
	}
}

//SignupMusician signs up a musician
func SignupMusician(session *mgo.Session) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		Vars := mux.Vars(r)
		fbID := Vars["id"]
		fbToken := Vars["token"]
		var profileSignupResponse ProfileResponse

		//Check if musician profile already exist
		if !CheckMusicianExistByID(session, fbID) {
			//So the profile does not exists, not we will create new
			_, err := InsertMusicianFromFB(session, fbToken, fbID)
			if err == nil {
				profileSignupResponse.Success = true
			} else {
				profileSignupResponse.Success = false
			}
		} else {
			toggle := UpdateLoggedStatusByID(session, fbID, true)
			if toggle {
				profileSignupResponse.Success = true
			} else {
				profileSignupResponse.Success = false
			}

		}

		Data, _ := json.Marshal(profileSignupResponse)
		JSON(string(Data))(w, r)

		return
	}
}

//LogoutMusician marks the musician as logged out
func LogoutMusician(session *mgo.Session) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		Vars := mux.Vars(r)
		musicianID := Vars["id"]
		var profileLogoutResponse ProfileResponse

		profileExists := CheckMusicianExistByID(session, musicianID)
		if profileExists {
			//User is found, now we can mark them as logged out
			UpdateLoggedStatusByID(session, musicianID, false)
			profileLogoutResponse.Success = true

		} else {
			//Could not find id to logout
			profileLogoutResponse.Success = false

		}

		Data, _ := json.Marshal(profileLogoutResponse)
		JSON(string(Data))(w, r)

		return
	}
}
