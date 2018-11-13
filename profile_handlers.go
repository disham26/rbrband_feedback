package main

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	mgo "gopkg.in/mgo.v2"
)

//BasicResponse struct
type BasicResponse struct {
	Success bool `json:"success"`
}

//ProfileHandler is called when the profile of the musician is created
func ProfileHandler(session *mgo.Session) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		Vars := mux.Vars(r)
		musicianID := Vars["id"]
		var bandResponse string
		if r.Method == "POST" {
			err := CreateBand(session, r.Form.Get("band_name"), r.Form.Get("band_location"), r.Form.Get("band_contact"), r.Form.Get("band_charge"), r.Form.Get("band_genre"), r.Form.Get("band_description"), musicianID)
			if err == nil {
				bandResponse = "Band Successfully Onboarded"
			}
		}

		result, goAhead := CheckIfLoggedIn(session, musicianID)
		if goAhead {
			QrString := GetQRCodeStringByID(session, musicianID)
			if QrString == "" {
				result = GenerateQRCodeString(session, musicianID, r)
			}
			// png, _ := qrcode.Encode(r.Host+r.URL.String(), qrcode.Medium, 256)
			// encoded := base64.StdEncoding.EncodeToString(png)
			t, _ := ParseTemplate("goadmin")
			log.Println("UserID is :", result.ID)
			page := &ProfilePage{
				Title:        "My Profile",
				Text:         result.FirstName,
				StaticHost:   getStaticHost(),
				JSON:         "{}",
				Email:        result.Email,
				Config:       "profile",
				LoggedIn:     result.IsLoggedIn,
				CurrentPage:  "profile",
				UserID:       result.ID,
				User:         result,
				BandResponse: bandResponse,
			}
			t.ExecuteTemplate(w, "profile", page)

		} else {
			//Profile logged out, going to landing page
			http.Redirect(w, r, "/", http.StatusSeeOther)
		}
		return
	}
}

//BandProfileHandler is called when the profile of the musician is created
func BandProfileHandler(session *mgo.Session) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		Vars := mux.Vars(r)
		musicianID := Vars["id"]
		if r.Method == "POST" {
			var band Band
			var users []User
			users = append(users, GetMusician(session, musicianID))
			r.ParseForm()
			band.BandCreator = musicianID
			band.Members = users
			band.BandName = r.Form.Get("band_name")
			band.Location = r.Form.Get("band_location")
			band.Contact = r.Form.Get("band_contact")
			band.Charges, _ = strconv.Atoi(r.Form.Get("band_charge"))
			band.Genre = r.Form.Get("band_genre")
			band.Description = r.Form.Get("band_description")
			new_band, err := InsertBand(session, band)
			if err == nil {
				//Update the Musician with bands associated
				//user := GetMusician(session, musicianID)
				err := UpdateBandCountMusician(session, musicianID, new_band)
				log.Println("Error to update bands in User is :", err)
			}
			log.Println(err)
		}

		result, goAhead := CheckIfLoggedIn(session, musicianID)
		if goAhead {
			QrString := GetQRCodeStringByID(session, musicianID)
			if QrString == "" {
				result = GenerateQRCodeString(session, musicianID, r)
			}
			// png, _ := qrcode.Encode(r.Host+r.URL.String(), qrcode.Medium, 256)
			// encoded := base64.StdEncoding.EncodeToString(png)
			t, _ := ParseTemplate("goadmin")
			log.Println("UserID is :", result.ID)
			page := &ProfilePage{
				Title:       "My Profile",
				Text:        result.FirstName,
				StaticHost:  getStaticHost(),
				JSON:        "{}",
				Email:       result.Email,
				Config:      "band_profile",
				LoggedIn:    result.IsLoggedIn,
				CurrentPage: "band_profile",
				UserID:      result.ID,
				User:        result,
			}
			t.ExecuteTemplate(w, "band_profile", page)

		} else {
			//Profile logged out, going to landing page
			http.Redirect(w, r, "/", http.StatusSeeOther)
		}
		return
	}
}

//GetQRHandler is called when the profile of the musician is created
func GetQRHandler(session *mgo.Session) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		Vars := mux.Vars(r)
		musicianID := Vars["id"]

		result, goAhead := CheckIfLoggedIn(session, musicianID)
		if goAhead {
			QrString := GetQRCodeStringByID(session, musicianID)
			if QrString == "" {
				result = GenerateQRCodeString(session, musicianID, r)
			}
			// png, _ := qrcode.Encode(r.Host+r.URL.String(), qrcode.Medium, 256)
			// encoded := base64.StdEncoding.EncodeToString(png)
			t, _ := ParseTemplate("goadmin")
			log.Println("UserID is :", result.ID)
			page := &ProfilePage{
				Title:       "My Profile",
				Text:        result.FirstName,
				StaticHost:  getStaticHost(),
				JSON:        "{}",
				Email:       result.Email,
				Config:      "qr_profile",
				LoggedIn:    result.IsLoggedIn,
				CurrentPage: "qr_profile",
				UserID:      result.ID,
				User:        result,
			}
			t.ExecuteTemplate(w, "qr_profile", page)

		} else {
			//Profile logged out, going to landing page
			http.Redirect(w, r, "/", http.StatusSeeOther)
		}
		return
	}
}
