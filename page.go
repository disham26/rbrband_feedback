package main

import (
	"html/template"
	"log"
	"os"
)

func getTemplatesDir() string {
	env := os.Getenv("ENV_NAME")
	log.Println("env is :", env)
	if env == "" {
		env = "local"
	}
	log.Println("env is :", env)
	return env
}

//ParseTemplate makes it easy to map the templates
func ParseTemplate(path string) (*template.Template, error) {
	//reads all files
	p, er := template.New("mustache").Delims("<<", ">>").ParseGlob(getTemplatesDir() + "static/" + path + "/templates/[a-z]*.mustache")
	//p, er := template.ParseGlob("views/admin/*.html")
	//p, er := template.ParseFiles(path, "views/admin/header.html", "views/admin/footer.html")
	return template.Must(p, er), er
}
