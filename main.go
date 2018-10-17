package main

import (
	"log"
	"net/http"
	"os"
	"rbrband/rbrband_feedback/creds"

	"github.com/gorilla/mux"
	mgo "gopkg.in/mgo.v2"
)

//Page Sample struct
type Page struct {
	Title      string
	Text       string
	StaticHost string
	Json       string
	Email      string
	Config     string
}
type Mongo struct {
	MongoUserID   string
	MongoPassword string
	MongoURL      string
}

var mongoSession *mgo.Session

func getEnv() string {
	env := os.Getenv("RUBBERBAND_ENV_NAME")
	if env == "" {

		return "local"
	}

	return env
}

func getStaticHost() string {
	env := getEnv()

	if env == "local" {
		return "http://localhost:5011"
	}
	return "production static link"

}
func getDBEnv(mongo Mongo) Mongo {
	log.Println("Inside getDBEnv")
	log.Println(getEnv())
	switch getEnv() {
	case "local":
		log.Println("local")
		mongo.MongoURL = creds.MONGO_LOCAL_URL
		mongo.MongoUserID = creds.MONGO_LOCAL_USERNAME
		mongo.MongoPassword = creds.MONGO_LOCAL_PASSWORD
		break
	case "staging":
		mongo.MongoURL = creds.MONGO_STAGING_URL
		mongo.MongoUserID = creds.MONGO_STAGING_USERNAME
		mongo.MongoPassword = creds.MONGO_STAGING_PASSWORD
		break
	case "prod":
		mongo.MongoURL = creds.MONGO_PROD_URL
		mongo.MongoUserID = creds.MONGO_PROD_USERNAME
		mongo.MongoPassword = creds.MONGO_PROD_PASSWORD
		break
	}
	log.Println(mongo.MongoPassword, mongo.MongoPassword, mongo.MongoURL)
	return mongo

}

//Person struct dummy, will remove later
type Person struct {
	Name  string
	Phone string
}

//InitDB function to make Mongo connection
func InitDB() *mgo.Session {
	var mongo Mongo
	mongo = getDBEnv(mongo)
	log.Println("Function to connect to DB")
	log.Println(mongo.MongoUserID, " is user")
	log.Println(mongo.MongoPassword, "is password")
	log.Println(mongo.MongoURL, "is URL")
	session, err := mgo.Dial("mongodb://" + mongo.MongoUserID + ":" + mongo.MongoPassword + "@" + mongo.MongoURL)
	mongoSession = session
	log.Println("session", session, "error:", err)
	//defer session.Close()

	// Optional. Switch the session to a monotonic behavior.
	session.SetMode(mgo.Monotonic, true)
	// c := session.DB("test").C("people")
	// err = c.Insert(&Person{"Ale", "+55 53 8116 9639"},
	// 	&Person{"Cla", "+55 53 8402 8510"})
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// result := Person{}
	// err = c.Find(bson.M{"name": "Ale"}).One(&result)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// fmt.Println("Phone:", result.Phone)
	return session

}
func indexHandler(session *mgo.Session) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		session := mongoSession.Copy()
		defer session.Close()
		defer log.Println("closing evverything")
		// c := session.DB("test").C("people")
		// err := c.Insert(&Person{"Prasang", "+91 76 9818 9874"},
		// 	&Person{"Disha", "+91 99 6086 0260"})
		// if err != nil {
		// 	log.Fatal(err)
		// }
		t, _ := ParseTemplate("goadmin")
		page := &Page{
			Title:      "Welcome to Rubberband",
			Text:       "Rubberband home",
			StaticHost: getStaticHost(),
			Json:       "{}",
			Email:      "prasang@rbrband.in",
			Config:     "index",
		}
		log.Println("going to execute template")
		log.Println(getStaticHost())
		t.ExecuteTemplate(w, "home", page)
		return
	}
}

func signupHandler(w http.ResponseWriter, r *http.Request) {
	t, _ := ParseTemplate("goadmin")
	log.Println("Signup called")
	page := &Page{
		Title:      "Welcome to Rubberband",
		Text:       "Signup/Login using",
		StaticHost: getStaticHost(),
		Json:       "{}",
		Email:      "prasang@rbrband.in",
		Config:     "index",
	}
	log.Println("going to execute template")
	log.Println(getStaticHost())
	t.ExecuteTemplate(w, "signup", page)

}
func main() {
	session := InitDB()
	defer mongoSession.Close()
	log.Println("MakeFile is running now,we can connect DB here")
	r := mux.NewRouter()

	r.HandleFunc("/", indexHandler(session))
	r.HandleFunc("/signup", signupHandler)
	http.ListenAndServe(":8000", r)
}
