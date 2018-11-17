package main

import (
	"html/template"
	"os"
	"time"
)

//Page Sample struct
type Page struct {
	Title       string
	Text        string
	StaticHost  string
	JSON        string
	Email       string
	Config      string
	LoggedIn    bool
	CurrentPage string
	FbID        string
	UserID      string
}

//FeedbackPage Sample struct
type FeedbackPage struct {
	Title       string
	Text        string
	StaticHost  string
	JSON        string
	Config      string
	LoggedIn    bool
	CurrentPage string
	FbID        string
	UserID      string
	Band        Band
	Message     string
}

//ProfilePage Sample struct
type ProfilePage struct {
	Title        string
	Text         string
	StaticHost   string
	JSON         string
	Email        string
	Config       string
	LoggedIn     bool
	CurrentPage  string
	FbID         string
	QR           string
	UserID       string
	User         User
	BandResponse string
}

//BandPage Sample struct
type BandPage struct {
	Title        string
	Text         string
	StaticHost   string
	JSON         string
	Email        string
	Config       string
	LoggedIn     bool
	CurrentPage  string
	FbID         string
	QR           string
	UserID       string
	User         User
	Band         Band
	BandResponse string
}

//Mongo struct
type Mongo struct {
	MongoUserID   string
	MongoPassword string
	MongoURL      string
}

//FBLocation struct to identify the location of musician after signup
type FBLocation struct {
	ID   string
	Name string
}

//ProfilePicture struct to get the URL of FB DP
type ProfilePicture struct {
	Data ImageURL `json:"data"`
}

//ImageURL is facebook profile picture URL
type ImageURL struct {
	URL string `json:"url"`
}

//User struct has all the details of a musician profile
type User struct {
	FirstName       string         `param:"first_name" json:"first_name" bson:"first_name"`
	Name            string         `json:"name" bson:"name"`
	ID              string         `json:"id" bson:"id"`
	Birthday        string         `json:"birthday" bson:"birthday"`
	Email           string         `json:"email" bson:"email"`
	Gender          string         `json:"gender" bson:"gender"`
	Location        FBLocation     `json:"location" bson:"location"`
	Link            string         `json:"link" bson:"link"`
	IsLoggedIn      bool           `json:"isLoggedIn" bson:"isLoggedIn"`
	UTS             time.Time      `json:"uts" bson:"uts"`
	BandsAssociated int            `json:"bands_associated" uts:"bands_associated"`
	ProfilePic      ProfilePicture `json:"picture" bson:"profile_pic" `
	Bands           []Band         `json:"bands" bson:"bands"`
	Rating          float64        `json:"rating" bson:"rating"`
	GuestsRated     int            `json:"guests_rated" bson:"guests_rated"`
}

//Band struct
type Band struct {
	ID          string    `json:"id" bson:"id"`
	BandName    string    `json:"band_name" bson:"band_name"`
	Members     []string  `json:"members" bson:"members"`
	Genre       string    `json:"genre" bson:"genre"`
	Description string    `json:"description" bson:"description"`
	Location    string    `json:"location" bson:"location"`
	Contact     string    `json:"contact" bson:"contect"`
	Age         time.Time `json:"age" bson:"age"`
	Charges     int       `json:"charges" bson:"charges" `
	BandCreator string    `json:"band_creator" bson:"band_creator"`
	UTS         time.Time `json:"uts" bson:"uts"`
	Rating      float64   `json:"rating" bson:"rating"`
	QRCode      string    `json:"qrcode" bson:"qrcode"`
	GuestsRated int       `json:"guests_rated" bson:"guests_rated"`
}

//Gig struct
type Gig struct {
	GigName      string
	Date         time.Time
	VenueName    string
	VenueCity    string
	VenueContact string
	Band         Band
}

//Blog struct
type Blog struct {
	Title   string
	Date    time.Time
	Content string
	Author  string
}

func getTemplatesDir() string {
	env := os.Getenv("ENV_NAME")
	if env == "" {
		env = "local"
	}
	return env
}

//ParseTemplate makes it easy to map the templates
func ParseTemplate(path string) (*template.Template, error) {
	//reads all files
	p, er := template.New("mustache").Delims("<<", ">>").ParseGlob(getTemplatesDir() + "static/" + path + "/templates/[a-z]*.mustache")
	return template.Must(p, er), er
}
