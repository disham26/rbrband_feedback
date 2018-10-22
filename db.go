package main

import (
	"encoding/base64"
	"log"
	"net/http"
	"rbrband/rbrband_feedback/creds"
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
	var err error
	res, _ := fb.Get("/"+id, fb.Params{
		"fields":       "first_name,name,id,birthday,email,gender,location,link,picture{url},videos{id}",
		"access_token": token,
	})
	var user User
	res.Decode(&user)
	if user.Name != "" {
		user.IsLoggedIn = true
		user.UTS = time.Now()
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
	err := c.Update(bson.M{"id": id}, bson.M{"$set": bson.M{"isloggedin": status}})
	log.Println("Error:", err)
	if err != nil {
		return false
	}

	return true
}

//GetMusician finds and sends the musician document from the collection
func GetMusician(session *mgo.Session, id string) User {
	result := User{}
	c := session.DB(creds.MONGO_DB).C(creds.MONGO_MUSICIAN_COLLECTION)
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
	result := User{}
	c := session.DB(creds.MONGO_DB).C(creds.MONGO_MUSICIAN_COLLECTION)
	c.Find(bson.M{"id": id}).One(&result)

	return result.QR
}

//GenerateQRCodeString function to generate and save the QR code in DB
func GenerateQRCodeString(session *mgo.Session, id string, r *http.Request) User {
	png, _ := qrcode.Encode(r.Host+r.URL.String(), qrcode.Medium, 256)
	encoded := base64.StdEncoding.EncodeToString(png)
	c := session.DB(creds.MONGO_DB).C(creds.MONGO_MUSICIAN_COLLECTION)
	c.Update(bson.M{"id": id}, bson.M{"$set": bson.M{"qr": encoded}})
	return GetMusician(session, id)
}
