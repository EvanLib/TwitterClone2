package models

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

var userPwPerpper = "SOMETHING DUMB AND SCRETE"

type User struct {
	gorm.Model
	Email          string  `gorm:"not null;unique_index"`
	Password       string  `gorm:"-"`
	HashedPassword string  `gorm:"not null"`
	Profile        Profile `gorm:"not null"`
}

type UserService interface {
	ByID(id uint) *User
	ByEmail(email string) *User
	Create(user *User) error
	Update(user *User) error
	Delete(id uint) error
	Authenticate(email, password string) *User
}

type UserGorm struct {
	*gorm.DB
}

func (ug *UserGorm) DestructiveReset() {
	//ug.DropTable(&User{})
	ug.AutoMigrate(&User{})
	ug.AutoMigrate(&Profile{})
}

func NewUserGorm(db *gorm.DB) *UserGorm {
	return &UserGorm{db}
}

func (ug *UserGorm) ByID(id uint) *User {
	return ug.byQuery(ug.DB.Where("id = ?", id))

}

func (ug *UserGorm) ByEmail(email string) *User {
	return ug.byQuery(ug.DB.Where("email = ?", email))
}

func (ug *UserGorm) byQuery(query *gorm.DB) *User {
	ret := &User{}
	err := query.First(ret).Error
	switch err {
	case nil:
		ug.Model(&ret).Related(&ret.Profile, "Profile")
		return ret
	case gorm.ErrRecordNotFound:
		return nil
	default:
		panic(err)
	}
}

func (ug *UserGorm) Create(user *User) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password+userPwPerpper), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.HashedPassword = string(hashedPassword)
	user.Password = ""

	return ug.DB.Create(user).Error
}
func (ug *UserGorm) Update(user *User) error {

	return ug.DB.Save(user).Error
}
func (ug *UserGorm) Delete(id uint) error {
	user := &User{Model: gorm.Model{ID: id}}
	return ug.DB.Delete(user).Error
}

func (ug *UserGorm) Authenticate(email, password string) *User {
	foundUser := ug.ByEmail(email)
	if foundUser == nil {
		//No User found
		return nil
	}

	err := bcrypt.CompareHashAndPassword([]byte(foundUser.HashedPassword), []byte(password+userPwPerpper))
	if err != nil {
		// Invalid password
		return nil
	}
	return foundUser
}
