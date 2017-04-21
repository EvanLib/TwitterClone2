package models

import (
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

type publicTweet struct {
	//Used for Json encoding

}
type Tweet struct {
	//gorm.Model        //Gives id and interaction
	ID        uint       `gorm:"primary_key"json:"id"`
	CreatedAt time.Time  ``
	UpdatedAt time.Time  `json:"-"`
	DeletedAt *time.Time `sql:"index"`
	Tweet     string     `json:"tweet"`
	Likes     int        `json:"likes"`
}

type TweetGorm struct {
	*gorm.DB //Create once and interact with it
}

type TweetService interface {
	ByID(id uint) *Tweet
	GetAll() []Tweet
	CreateTweet(tweet *Tweet) error
}

func (tg *TweetGorm) CreateMockTweets() {
	for i := 0; i < 50; i++ {
		message := fmt.Sprintf("Tweet string: %d", i)
		tweet := &Tweet{
			Tweet: message,
			Likes: i,
		}
		tg.DB.Create(tweet)
	}
}

func (tg *TweetGorm) DatabaseStuff() {
	fmt.Println("Tweets Database loaded....")
	//tg.CreateMockTweets()
}
func NewTweetGorm(connectionInfo string) (*TweetGorm, error) {
	db, err := gorm.Open("mysql", connectionInfo)
	if err != nil {
		return nil, err
	}
	return &TweetGorm{db}, nil
}

func (tg *TweetGorm) DesctructiveReset() {
	tg.AutoMigrate(Tweet{})
}

//Crud stuff
func (tg *TweetGorm) GetAll() []Tweet {
	allTweets := []Tweet{}
	tg.DB.Find(&allTweets)
	return allTweets
}
func (tg *TweetGorm) ByID(id uint) *Tweet {
	return tg.byQuery(tg.DB.Where("id = ?", id))
}

func (tg *TweetGorm) CreateTweet(tweet *Tweet) error {
	return tg.DB.Create(tweet).Error
}
func (tg *TweetGorm) byQuery(query *gorm.DB) *Tweet {
	ret := &Tweet{}
	err := query.First(ret).Error
	switch err {
	case nil:
		return ret
	case gorm.ErrRecordNotFound:
		return nil
	default:
		panic(err)
	}
}
