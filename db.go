package main

import (
	"encoding/base64"
	"log"
	"net/http"
	"rbrband/rbrband_feedback/creds"
	"strconv"
	"time"

	fb "github.com/huandu/facebook"
	qrcode "github.com/skip2/go-qrcode"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

//InitDB function to make Mongo connection during startup
func InitDB() *mgo.Session {
	var mongo Mongo
	mongo = getDBEnv(mongo)
	//Function to connect to DB")
	session, _ := mgo.Dial("mongodb://" + mongo.MongoUserID + ":" + mongo.MongoPassword + "@" + mongo.MongoURL)
	mongoSession = session
	session.SetMode(mgo.Monotonic, true)
	return session
}

//InsertDetailsFromFB function will insert the values to mongo, first it will check whether a user exists by calling another function,
//if yes, only update will be called, else a new document will be inserted in the collection "musicians" in DB "users"
func InsertMusicianFromFB(session *mgo.Session, token string, id string) (string, error) {
	log.Println("Creating a new musician")
	log.Println("id is :", id)
	var err error
	res, err := fb.Get("/"+id, fb.Params{
		"fields":       "first_name,name,id,birthday,email,gender,location,link,picture{url},videos{id}",
		"access_token": token,
	})
	log.Println(res)
	log.Println("Error to create the profile is :", err)
	var user User
	res.Decode(&user)
	if user.Name != "" {
		user.IsLoggedIn = true
		user.UTS = time.Now()
		user.Rating = 0
		user.GuestsRated = 0
		c := session.DB(creds.MONGO_DB).C(creds.MONGO_MUSICIAN_COLLECTION)
		err = c.Insert(user)
		if err != nil {
			log.Print(err)
		}
	} else {
		user.IsLoggedIn = false
	}

	return user.ID, err
}

//InsertBand function will insert the values to mongo, first it will check whether a user exists by calling another function,
//if yes, only update will be called, else a new document will be inserted in the collection "musicians" in DB "users"
func InsertBand(session *mgo.Session, band Band) (Band, error) {
	log.Println("Creating a new band")
	var err error

	band.UTS = time.Now()
	c := session.DB(creds.MONGO_DB).C(creds.MONGO_BANDS_COLLECTION)
	err = c.Insert(band)
	if err != nil {
		log.Print(err)
	}
	return band, err
}

//CheckMusicianExistByID function
func CheckMusicianExistByID(session *mgo.Session, id string) bool {
	result := User{}
	c := session.DB(creds.MONGO_DB).C(creds.MONGO_MUSICIAN_COLLECTION)
	err := c.Find(bson.M{"id": id}).One(&result)
	if err != nil {
		log.Println(err)
	}
	if result.ID != "" {
		return true
	}
	return false
}

//UpdateLoggedStatusByID function to mark the loggedIn feature
func UpdateLoggedStatusByID(session *mgo.Session, id string, status bool) bool {
	log.Println("inside UpdateLoggedStatusByID for ", id, status)
	c := session.DB(creds.MONGO_DB).C(creds.MONGO_MUSICIAN_COLLECTION)
	err := c.Update(bson.M{"id": id}, bson.M{"$set": bson.M{"isLoggedIn": status}})
	log.Println("Error:", err)
	if err != nil {
		return false
	}

	return true
}

//UpdateBandRating function to update the rating of the band
func UpdateBandRating(session *mgo.Session, band Band, rating string) error {
	log.Println("inside UpdateBandRating for ", band.ID, rating)
	log.Println("Current Rating is :", band.Rating)
	log.Println("Guests rated:", band.GuestsRated)
	//log.Println("Full band details are:", band)
	log.Println("Increasing guests rated:")
	var err error
	newGuestsRated := band.GuestsRated + 1
	ratingFloat, _ := strconv.ParseFloat(rating, 64)
	c := session.DB(creds.MONGO_DB).C(creds.MONGO_BANDS_COLLECTION)
	if band.GuestsRated == 0 {
		log.Println("So we are here")
		err = c.Update(bson.M{"id": band.ID}, bson.M{"$set": bson.M{"guests_rated": 1}})
		//log.Println("First error:", err)
		ratingFloat, _ := strconv.ParseFloat(rating, 64)
		err = c.Update(bson.M{"id": band.ID}, bson.M{"$set": bson.M{"rating": FloatToTwo(ratingFloat)}})
		//log.Println("Second error:", err)
	} else {
		massiveRating := float64(band.GuestsRated) * band.Rating
		log.Println("Guests Rated:", band.GuestsRated)
		log.Println("Float guests rated:", float64(band.GuestsRated))
		log.Println("Band Rating:", band.Rating)
		log.Println("MassiveRating is :", massiveRating)
		newMassiveRating := massiveRating + ratingFloat

		newRating := newMassiveRating / float64(newGuestsRated)
		//newRating := ((band.Rating * float64(band.GuestsRated)) + ratingFloat) / float64(band.GuestsRated+1)

		err = c.Update(bson.M{"id": band.ID}, bson.M{"$set": bson.M{"guests_rated": newGuestsRated}})
		err = c.Update(bson.M{"id": band.ID}, bson.M{"$set": bson.M{"rating": FloatToTwo(newRating)}})

	}

	return err
}

//UpdateMusiciansRating function to update the rating of the musician
func UpdateMusiciansRating(session *mgo.Session, bandID string, rating string) error {
	band := GetBand(session, bandID)
	log.Println("Musicians in the band are:", len(band.Members))
	var err error
	for _, musicianID := range band.Members {
		log.Println("Updating rating for :", musicianID)
		log.Println(musicianID)
		musician := GetMusician(session, musicianID)
		log.Println("Current Musician Rating is :", musician.Rating)
		log.Println("Name :", musician.FirstName)
		log.Println("Guests rated:", musician.GuestsRated)
		//log.Println("Full band details are:", band)
		log.Println("Increasing guests rated:")

		newGuestsRated := musician.GuestsRated + 1
		ratingFloat, _ := strconv.ParseFloat(rating, 64)
		c := session.DB(creds.MONGO_DB).C(creds.MONGO_MUSICIAN_COLLECTION)
		if musician.GuestsRated == 0 {
			log.Println("So we are here")
			err = c.Update(bson.M{"id": musician.ID}, bson.M{"$set": bson.M{"guests_rated": 1}})
			//log.Println("First error:", err)
			ratingFloat, _ := strconv.ParseFloat(rating, 64)
			err = c.Update(bson.M{"id": musician.ID}, bson.M{"$set": bson.M{"rating": FloatToTwo(ratingFloat)}})
			//log.Println("Second error:", err)
		} else {
			massiveRating := float64(musician.GuestsRated) * musician.Rating
			log.Println("Guests Rated:", musician.GuestsRated)
			log.Println("Float guests rated:", float64(musician.GuestsRated))
			log.Println("Band Rating:", musician.Rating)
			log.Println("MassiveRating is :", massiveRating)
			newMassiveRating := massiveRating + ratingFloat

			newRating := newMassiveRating / float64(newGuestsRated)
			//newRating := ((band.Rating * float64(band.GuestsRated)) + ratingFloat) / float64(band.GuestsRated+1)

			err = c.Update(bson.M{"id": musician.ID}, bson.M{"$set": bson.M{"guests_rated": newGuestsRated}})
			err = c.Update(bson.M{"id": musician.ID}, bson.M{"$set": bson.M{"rating": FloatToTwo(newRating)}})
		}

	}

	return err
}

//GetMusician finds and sends the musician document from the collection
func GetMusician(session *mgo.Session, id string) User {
	result := User{}
	c := session.DB(creds.MONGO_DB).C(creds.MONGO_MUSICIAN_COLLECTION)
	c.Find(bson.M{"id": id}).One(&result)

	return result
}

//GetBand function to be called to check if a band exists in DB
func GetBand(session *mgo.Session, id string) Band {
	result := Band{}
	c := session.DB(creds.MONGO_DB).C(creds.MONGO_BANDS_COLLECTION)
	c.Find(bson.M{"id": id}).One(&result)

	return result
}

//UpdateBandCountMusician function to update the bands associated of the musician
func UpdateBandCountMusician(session *mgo.Session, id string, band Band) error {
	user := GetMusician(session, id)
	//return error
	c := session.DB(creds.MONGO_DB).C(creds.MONGO_MUSICIAN_COLLECTION)
	err := c.Update(bson.M{"id": id}, bson.M{"$set": bson.M{"bandsassociated": user.BandsAssociated + 1}})
	err = c.Update(bson.M{"id": id}, bson.M{"$set": bson.M{"bands": band}})
	log.Println("Error:", err)
	return err

}

//GetQRCodeStringByID function gets QR code string of the id
func GetQRCodeStringByID(session *mgo.Session, id string) string {
	result := Band{}
	c := session.DB(creds.MONGO_DB).C(creds.MONGO_BANDS_COLLECTION)
	c.Find(bson.M{"id": id}).One(&result)

	return result.QRCode
}

//GenerateQRCodeString function to generate and save the QR code in DB
func GenerateQRCodeString(session *mgo.Session, band_id string, r *http.Request) {
	url := r.Host + "/feedback" + "/" + band_id
	png, _ := qrcode.Encode(url, qrcode.Medium, 256)
	encoded := base64.StdEncoding.EncodeToString(png)
	c := session.DB(creds.MONGO_DB).C(creds.MONGO_BANDS_COLLECTION)
	c.Update(bson.M{"id": band_id}, bson.M{"$set": bson.M{"qrcode": encoded}})

}

//GetBandByUserId finds and sends the band document from the collection
func GetBandByUserId(session *mgo.Session, id string) []Band {
	result := []Band{}
	c := session.DB(creds.MONGO_DB).C(creds.MONGO_BANDS_COLLECTION)
	c.Find(bson.M{"members": id}).All(&result)
	log.Println(len(result))
	return result
}
