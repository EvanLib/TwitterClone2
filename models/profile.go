package models

import (
	"time"

	"github.com/jinzhu/gorm"
)

type Profile struct {
	ID        uint       `gorm:"primary_key" json:"id"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"-"`
	DeletedAt *time.Time `sql:"index"`
	UserID    uint       `json:"-"`
	Name      string     `gorm:"not null" json:"name"`
	UserName  string     `gorm:"not null" json:"username"`
	Tweets    []Tweet    `json:"-"`
}

type ProfileService interface {
	ByID(id uint) *Profile
	Update(profile *Profile) error
	Delete(id uint) error
}

type ProfileGorm struct {
	*gorm.DB
}

func NewProfileGorm(db *gorm.DB) *ProfileGorm {
	return &ProfileGorm{db}
}

//ProfileService
func (pg *ProfileGorm) ByID(id uint) *Profile {
	return pg.byQuery(pg.DB.Where("id = ?", id))

}

func (pg *ProfileGorm) byQuery(query *gorm.DB) *Profile {
	ret := &Profile{}
	err := query.First(ret).Error
	switch err {
	case nil:
		// ug.Model(&ret).Related(&ret.Profile, "Profile")
		return ret
	case gorm.ErrRecordNotFound:
		return nil
	default:
		panic(err)
	}
}

func (pg *ProfileGorm) Update(profile *Profile) error {
	return pg.DB.Save(profile).Error
}
func (pg *ProfileGorm) Delete(id uint) error {
	profile := &Profile{ID: id}
	return pg.DB.Delete(profile).Error
}
