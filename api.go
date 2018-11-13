package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	mgo "gopkg.in/mgo.v2"
)

//ProfileResponse struct
type ProfileResponse struct {
	Success bool `json:"success"`
}

//BandResponse struct to show the bands with which a user is associated
type BandResponse struct {
	Success bool   `json:"success"`
	Data    []Band `json:"data"`
	User    User   `json:"user"`
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

//CreateBand api to create a band
func CreateBand(session *mgo.Session, bandName string, bandLocation string, bandContact string, bandCharge string, bandGenre string, bandDesc string, musicianID string) error {
	var band Band
	var users []User
	users = append(users, GetMusician(session, musicianID))
	band.BandCreator = musicianID
	band.Members = users
	band.BandName = bandName
	band.Location = bandLocation
	band.Contact = bandContact
	band.Charges, _ = strconv.Atoi(bandCharge)
	band.Genre = bandGenre
	band.Description = bandDesc
	new_band, err := InsertBand(session, band)
	if err == nil {
		//Update the Musician with bands associated
		//user := GetMusician(session, musicianID)
		err := UpdateBandCountMusician(session, musicianID, new_band)
		log.Println("Error to update bands in User is :", err)
	}
	log.Println(err)
	return err
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

//GetBandsAssociated gets all the bands with which the user is associated
func GetBandsAssociated(session *mgo.Session) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		Vars := mux.Vars(r)
		musicianID := Vars["id"]
		log.Println("BandsAssociated called")
		var bandResponse BandResponse
		bandResponse.User = GetMusician(session, musicianID)
		bands := GetBandByUserId(session, musicianID)
		log.Println("Bands are :", bands)
		log.Println(len(bands))
		if len(bands) > 0 {
			bandResponse.Success = true
			bandResponse.Data = bands
		} else {
			bandResponse.Success = false
		}
		Data, _ := json.Marshal(bandResponse)
		JSON(string(Data))(w, r)

		return
	}
}
