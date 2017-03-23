package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/EvanLib/Twitter_Clone_Backend/authentication"
	"github.com/EvanLib/Twitter_Clone_Backend/controllers"
	"github.com/EvanLib/Twitter_Clone_Backend/models"
	"github.com/gorilla/mux"
)

func main() {
	//Twitter Clone by Evan Grayson
	//Auth testing

	//Database gorm stuff
	tg, err := models.NewTweetGorm("root:somepassword@/twitter_clone?charset=utf8&parseTime=True&loc=Local")
	ug, err := models.NewUserGorm("root:somepassword@/twitter_clone?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		panic(err)
	}
	//Create controllers
	tweetsController := controllers.NewTweetsController(tg)
	usersController := controllers.NewUsers(ug)

	tg.DatabaseStuff()
	//Create mux
	r := mux.NewRouter()
	r.HandleFunc("/", Index)
	r.HandleFunc("/api", Index)
	r.HandleFunc("/api/tweets", tweetsController.Index)

	//Authentication stuff
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
