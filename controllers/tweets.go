package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/EvanLib/TwitterClone2/authentication"
	"github.com/EvanLib/TwitterClone2/models"
	"github.com/gorilla/context"
	"github.com/gorilla/mux"
)

type Tweets struct {
	models.TweetService
	models.ProfileService
}

type TweetJSON struct {
	Tweet string `json:"tweet"`
}

func NewTweetsController(tweetsService models.TweetService, profileService models.ProfileService) *Tweets {
	return &Tweets{
		TweetService:   tweetsService,
		ProfileService: profileService,
	}

}
func (t Tweets) Get(w http.ResponseWriter, r *http.Request) {
	//Grab the ID
	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["id"], 10, 0)
	if err != nil {
		fmt.Println(err)
		return
	}

	//Find user from ID
	tweet := t.TweetService.ByID(uint(id))
	if tweet == nil {
		http.Error(w, "ID not found", http.StatusNotFound)
		return
	}

	// Attach Profile
	profile := t.ProfileService.ByID(tweet.ProfileID)
	tweet.Profile = *profile

	//Json Marshal
	jasonData, err := json.Marshal(tweet)
	if err != nil {
		panic(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jasonData)
}
func (t Tweets) Index(w http.ResponseWriter, r *http.Request) {
	tweets := t.TweetService.GetAll()

	for i := range tweets {
		tweets[i].Profile = *t.ProfileService.ByID(tweets[i].ProfileID)
	}
	jasonData, err := json.Marshal(tweets)
	if err != nil {
		panic(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*") //Needed for cross domain.
	w.Write(jasonData)
}

func (t Tweets) Create(w http.ResponseWriter, r *http.Request) {
	userContext := context.Get(r, "user_info")
	userJWT := userContext.(*authentication.UserTest)

	//Create Tweet
	var tweetJSON TweetJSON
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&tweetJSON)

	if err != nil {
		log.Println("API Tweet ERROR: ", err)
		return
	}

	tweet := models.Tweet{
		Tweet: tweetJSON.Tweet,
		Likes: 0,
	}
	err = t.TweetService.CreateTweet(&tweet)
	if err != nil {
		fmt.Println(err)
	}

	//Append to user's profile
	profile := t.ProfileService.ByID(userJWT.ID)
	if profile == nil {
		fmt.Println("Create Tweet: User not found.")
		return
	}

	profile.Tweets = append(profile.Tweets, tweet)

	t.ProfileService.Update(profile)

}
