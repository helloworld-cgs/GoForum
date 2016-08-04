package m

import (
)
import (
	"time"
)


type Notification struct {
	Id        uint        `orm:"pk"`
	CreatedAt time.Time   `orm:"auto_now_add"`
	UpdatedAt time.Time   `orm:"auto_now"`
	DeletedAt time.Time
	User        *User    `orm:"rel(fk)"`
	Related     *User    `orm:"rel(fk)"`
	UserID      uint
	RelatedID   uint     `orm:"default(0)"`
	SubjectType int
	IsRead      bool     `orm:"default(false)"`

	Data string   `orm:"default('')"`
}

func (Notification) TableName() string {
	return "notification"
}

type PostMessage struct {
	Id        uint        `orm:"pk"`
	CreatedAt time.Time   `orm:"auto_now_add"`
	UpdatedAt time.Time   `orm:"auto_now"`
	DeletedAt time.Time
	User        *User    `orm:"rel(fk)"`
	Related     *User    `orm:"rel(fk)"`
	UserID      uint
	RelatedID   uint     `orm:"default(0)"`
	SubjectType int
	IsRead      bool     `orm:"default(false)"`

	RelatedUsername string   `orm:"default('')"`
	PostID          uint
	PostTitle       string
	Summary         string  `orm:"size(255)"`
	Quote           string  `orm:"size(255)"`
}

func (PostMessage) TableName() string {
	return "post_message"
}
