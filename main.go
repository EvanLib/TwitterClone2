package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/EvanLib/TwitterClone2/authentication"
	"github.com/EvanLib/TwitterClone2/controllers"
	"github.com/EvanLib/TwitterClone2/models"
	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
)

func main() {
	//Twitter Clone by Evan Grayson
	//Auth testing

	//Database gorm stuff
	tg, err := models.NewTweetGorm("root:lol626465@/twitter_clone?charset=utf8&parseTime=True&loc=Local")
	ug, err := models.NewUserGorm("root:lol626465@/twitter_clone?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		panic(err)
	}
	tg.Debug()
	ug.DestructiveReset()
	//Create controllers
	tweetsController := controllers.NewTweetsController(tg, ug)
	usersController := controllers.NewUsers(ug)
	//Negroni middleware

	tg.DatabaseStuff()
	//Create mux
	r := mux.NewRouter()
	r.HandleFunc("/", Index)
	r.HandleFunc("/api", Index)
	r.Handle("/api/tweets", negroni.New(
		negroni.HandlerFunc(authentication.HandlerWithNext),
		negroni.Wrap(http.HandlerFunc(tweetsController.Index)),
	)).Methods("GET")

	r.Handle("/api/tweets", negroni.New(
		negroni.HandlerFunc(authentication.HandlerWithNext),
		negroni.Wrap(http.HandlerFunc(tweetsController.Create)),
	)).Methods("POST")

	r.HandleFunc("/api/auth/create", usersController.Create)
	r.HandleFunc("/api/auth/login", usersController.Login).Methods("POST")

	r.HandleFunc("/api/auth/tokenget", authentication.AuthHandler)
	r.HandleFunc("/api/auth/login2", authentication.RestrictedHandler).Methods("POST")
	r.HandleFunc("/api/auth/tokentest", authentication.RestrictedHandler)

	//Create the http server
	srv := &http.Server{
		Handler: r,
		Addr:    "127.0.0.1:3000",
		//Good practive: enforce timeouts for servers you create...
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}

func Index(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprint(w, "API FOR A TWITTER CLONE. SO FUCK OFF.")
}
