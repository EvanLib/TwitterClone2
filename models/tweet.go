package models

import (
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

type Tweet struct {
	//gorm.Model        //Gives id and interaction
	ID        uint       `gorm:"primary_key" json:"id"`
	CreatedAt time.Time  ``
	UpdatedAt time.Time  `json:"-"`
	DeletedAt *time.Time `sql:"index"`
	Tweet     string     `json:"tweet"`
	Likes     int        `json:"likes"`
	ProfileID uint
	Profile   Profile
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
func NewTweetGorm(db *gorm.DB) *TweetGorm {
	return &TweetGorm{db}
}

func (tg *TweetGorm) DesctructiveReset() {
	//tg.DropTable(Tweet{})
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
