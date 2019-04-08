package main

import (
	"log"
	"net/http"

	"github.com/akshitgrover/saas_k8s/controllers"
	"github.com/globalsign/mgo"
)

func main() {
	session, err := mgo.Dial("mongodb://localhost:27017")
	db := session.DB("saas_k8s")
	if err != nil {
		log.Fatal(err)
	}
	http.HandleFunc("/register", controllers.Register(db))
	http.HandleFunc("/create", controllers.Create(db))
	http.ListenAndServe(":3004", nil)
}
