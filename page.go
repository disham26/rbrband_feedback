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

//Page Sample struct
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
	FirstName       string         `param:"first_name" json:"first_name"`
	Name            string         `json:"name"`
	ID              string         `json:"id"`
	Birthday        string         `json:"birthday"`
	Email           string         `json:"email"`
	Gender          string         `json:"gender"`
	Location        FBLocation     `json:"location"`
	Link            string         `json:"link"`
	IsLoggedIn      bool           `json:"isLoggedIn"`
	UTS             time.Time      `json:"uts"`
	QR              string         `json:"QR"`
	BandsAssociated int            `json:"bandsAssociated"`
	ProfilePic      ProfilePicture `json:"picture"`
	Bands           []Band
}

//Band struct
type Band struct {
	BandName    string    `json:"band_name"`
	Members     []User    `json:"user"`
	Genre       string    `json:"genre"`
	Description string    `json:"descroption"`
	Location    string    `json:"location"`
	Contact     string    `json:"contact"`
	Age         time.Time `json:"age"`
	Charges     int       `json:"charges"`
	BandCreator string    `json:"band_creator"`
	UTS         time.Time `json:"uts"`
}

//Gig struct
type Gig struct {
	GigName      string
	Date         time.Time
	VenueName    string
	VenueCity    string
	VenueContact string
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
