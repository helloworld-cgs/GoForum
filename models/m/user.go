package m

import (
	"github.com/jinzhu/gorm"
	"./../database"
	"fmt"
)

const (
	STATUS_ACTIVE int = 10
)

type User struct {
	gorm.Model
	Email      string
	Tel        string
	Name       string
	Password   string  `gorm:"column:password_hash"`
	AuthKey    string  `gorm:"column:auth_key"`
	ResetToken string  `gorm:"column:password_reset_token"`
	Status     int     `gorm:"default:10"`
	Profile    Profile `gorm:"ForeignKey:UserRefer"`
}

type Profile struct {
	//ID uint `gorm:"primary_key"`
	UserRefer      uint
	Avatar         string  `gorm:"default:'default.png'"`
	Coins          int     `gorm:"default:0"`
	PostCount      int     `gorm:"default:0"`
	CommentCount   int     `gorm:"default:0"`
	FollowedCount  int     `gorm:"default:0"`
	FollowingCount int     `gorm:"default:0"`
	Bio            string  `gorm:"default:'这家伙很懒,什么都没有'"`
}

func (Profile) TableName() string {
	return "profile"
}

func (User) TableName() string {
	return "user"
}

func (u *User) GetUserById(id uint) {
	fmt.Println(id)
	database.DB.Preload("Profile").First(&u, id)
}

type Follow struct {
	Follower    User
	Following   User
	FollowerID  uint
	FollowingID uint
}