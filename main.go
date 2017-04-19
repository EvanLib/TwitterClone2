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
	tg, err := models.NewTweetGorm("root:somepassword@/twitter_clone?charset=utf8&parseTime=True&loc=Local")
	ug, err := models.NewUserGorm("root:somepassword@/twitter_clone?charset=utf8&parseTime=True&loc=Local")
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
	r.HandleFunc("/api", Index)
	r.Handle("/api/tweets", negroni.New(
		negroni.HandlerFunc(authentication.HandlerWithNext),
		negroni.Wrap(http.HandlerFunc(tweetsController.Index)),
	)).Methods("GET")

	r.Handle("/api/tweets", negroni.New(
		negroni.HandlerFunc(authentication.HandlerWithNext),
		negroni.Wrap(http.HandlerFunc(tweetsController.Create)),
	)).Methods("POST")

	//Accounnt CRUD
	r.HandleFunc("/api/auth/create", usersController.Create)
	r.HandleFunc("/api/auth/login", usersController.Login).Methods("POST")

	r.Handle("/api/auth/logout", negroni.New(
		negroni.HandlerFunc(authentication.HandlerWithNext),
		negroni.Wrap(http.HandlerFunc(usersController.Logout)),
	)).Methods("DELETE")

	//JWT Token Auth
	r.HandleFunc("/api/auth/tokenget", authentication.AuthHandler)
	r.HandleFunc("/api/auth/login2", authentication.RestrictedHandler).Methods("POST")
	r.HandleFunc("/api/auth/tokentest", authentication.RestrictedHandler)

	//Create the http server
	http.Handle("/", &MyServer{r}) //Five days worth of searching

	srv := &http.Server{
		Handler: nil,
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

//Preflight Server stuff for CORS
//http://stackoverflow.com/questions/12830095/setting-http-headers-in-golang
type MyServer struct {
	r *mux.Router
}

func (s *MyServer) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	if origin := req.Header.Get("Origin"); origin != "" {
		rw.Header().Set("Access-Control-Allow-Origin", origin)
		rw.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		rw.Header().Set("Access-Control-Allow-Headers",
			"Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	}
	// Stop here if its Preflighted OPTIONS request
	if req.Method == "OPTIONS" {
		return
	}
	// Lets Gorilla work
	s.r.ServeHTTP(rw, req)
}
