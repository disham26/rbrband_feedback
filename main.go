package main

import (
	"log"
	"net/http"
	"os"
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

func getEnv() string {
	env := os.Getenv("ENV_NAME")
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
	return "http://localhost:5011"

}

func indexHandler(w http.ResponseWriter, r *http.Request) {
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
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/signup", signupHandler)
	http.ListenAndServe(":8000", nil)
}
