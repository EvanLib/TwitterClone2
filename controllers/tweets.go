package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/EvanLib/Twitter_Clone_Backend/models"
)

type Tweets struct {
	models.TweetService
}

func NewTweetsController(tweetsService models.TweetService) *Tweets {
	return &Tweets{
		TweetService: tweetsService,
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
