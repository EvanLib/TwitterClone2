package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/EvanLib/TwitterClone2/authentication"
	"github.com/EvanLib/TwitterClone2/models"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/context"
	"github.com/gorilla/mux"
)

type Users struct {
	models.UserService
}

type UserJSON struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

//CRUD

func NewUsers(us models.UserService) *Users {
	return &Users{

		UserService: us,
	}
}

//Post /signup
func (u *Users) Create(w http.ResponseWriter, r *http.Request) {

	user := &models.User{
		Email:    "Evan@evandoesstuff.com",
		Password: "Somepass",
		Profile: models.Profile{
			Name:     "Evan Grayson",
			UserName: "Twiggy_io",
		},
	}

	if err := u.UserService.Create(user); err != nil {
		panic(err)
	}
	if err := u.signIn(w, user); err != nil {
		http.Error(w, "Something went wrong.", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/cookietest", http.StatusFound)
}

//Authentication
func (u *Users) Login(w http.ResponseWriter, r *http.Request) {

	var userJSON UserJSON
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&userJSON)

	if err != nil {
		log.Println("API LOGIN ERROR: ", err)
		return
	}

	user := u.Authenticate(userJSON.Email, userJSON.Password)
	if user == nil {
		http.Error(w, "No user found.", http.StatusUnauthorized)
		return
	}

	tokenString, err := authentication.CreateToken(user.Email, user.ID, user.Profile.ID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, "Sorry, error while Signing Token!")
		log.Printf("Token Signing error: %v\n", err)
		return
	}
	w.Header().Set("Content-Type", "application/jwt")
	w.WriteHeader(http.StatusOK)

	b, err := json.Marshal(map[string]string{
		"token":    tokenString,
		"username": user.Profile.UserName, // Fix the Profile
	})
	w.Write(b)

}
func (u *Users) Logout(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	fmt.Println("Logout Called")
	//User info from context
	userContext := context.Get(r, "user_info")
	fmt.Println(userContext)
	// TODO: write a handler for invalidating tokens. Redis backend maybe.
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

func (u *Users) signIn(w http.ResponseWriter, user *models.User) error {
	//Strip JSON login package

	return nil
}

//CRUD For Profiles
func (u *Users) GetProfile(w http.ResponseWriter, r *http.Request) {
	//Grab the ID
	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["id"], 10, 0)
	if err != nil {
		fmt.Println(err)
		return
	}

	//Find user from ID
	user := u.ByID(uint(id))
	if user == nil {
		http.Error(w, "ID not found", http.StatusNotFound)
		return
	}

	//Put into json
	jasonData, err := json.Marshal(user.Profile)
	if err != nil {
		panic(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jasonData)
}
