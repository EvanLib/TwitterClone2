package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/EvanLib/TwitterClone2/authentication"
	"github.com/EvanLib/TwitterClone2/models"
	"github.com/gorilla/context"
)

type Tweets struct {
	models.TweetService
	models.UserService
}

func NewTweetsController(tweetsService models.TweetService, userService models.UserService) *Tweets {
	return &Tweets{
		TweetService: tweetsService,
		UserService:  userService,
	}

}
func (t Tweets) Index(w http.ResponseWriter, r *http.Request) {
	tweets := t.TweetService.GetAll()
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
	tweet := &models.Tweet{
		Tweet: "This is a tester tweet",
		Likes: 10,
	}

	err := t.TweetService.CreateTweet(tweet)
	if err != nil {
		fmt.Println(err)
	}
	//Append to user
	user := t.UserService.ByID(userJWT.ID)
	if user == nil {
		fmt.Println("Create Tweet: User not found.")
		return
	}

	t.UserService.AddTweet(user, tweet)
}
