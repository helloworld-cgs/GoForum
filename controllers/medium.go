package controllers

import (
	"time"
	"gensh.me/goforum/models/m"
	"gensh.me/goforum/models/database"
)

type Person struct {
	ID     uint
	Name   string
	Avatar string
}

type Profile struct {
	m.Profile
	Name        string
	HasFollowed bool
	Edit        bool
	IsLogin     bool
}

func GetProfileById(my_id uint, uid uint) (profile Profile) {
	user := m.User{}
	//RW database.DB.Preload("Profile").First(&user, uid)
	database.O.QueryTable("user").Filter("id",uid).RelatedSel("Profile").One(&user)
	follow_count ,_ := database.O.QueryTable("follows").Filter("follower_id",my_id).Filter("following_id",uid).Count()
	//RW database.DB.Model(&m.Follow{}).Where("follower_id = ? AND following_id = ?", my_id, uid).Count(&follow_count)
	profile = Profile{Name:user.Name, Profile:*user.Profile, Edit:(my_id == uid),
		HasFollowed:follow_count != 0, IsLogin:!(my_id == 0)}
	return
}

//<post item >
type PostItem struct {
	Person
	PostID       uint
	Title        string
	ViewCount    int
	CommentCount int
}

func DBHotPostsConvert(dbHotPosts *[]m.Posts) (*[]PostItem) {
	postItems := make([]PostItem, 0, len(*dbHotPosts))  //dbHotPosts to mHotPosts
	for _, db_hot := range *dbHotPosts {
		postItems = append(postItems, PostItem{PostID:db_hot.Id, Title:db_hot.Title,
			ViewCount:db_hot.ViewCount, CommentCount:db_hot.CommentCount,
			Person:Person{ID:db_hot.Author.Id, Name:db_hot.Author.Name, Avatar:db_hot.Author.Profile.Avatar}});
	}
	return &postItems
}
//</post item >


// <structs for Category >
//type Tag struct {
//
//}

//type Topic struct {
//
//}

type Category struct {
	Tags   []m.Tag
	Topics []m.Topic
}

func (this *Category)NewInstant() {
	database.O.QueryTable("topic").Filter("visible", true).All(&this.Topics)
	database.O.QueryTable("tags").Filter("visible", true).All(&this.Tags)
	//RW database.DB.Where("visible = ?", true).Find(&this.Topics)
	//RW database.DB.Where("visible = ?", true).Find(&this.Tags)
}
// </structs for Category >
//  </struct for me >
type UserStatus struct {
	Person
	IsLogin bool
}
//< /struct for me page >

/**<used for Post detail> */
type PostDetail struct {
	ID           uint
	//Topic
	Title        *string
	Content      *string
	IsMobile     bool
	Sticky       bool
	CommentCount int
	ViewCount    int
	LastReplayAt *time.Time
	CreatedAt    *time.Time
	UpdatedAt    *time.Time
}

func (this *PostDetail) NewInstant(p *m.Posts) {
	this.ID = p.Id
	this.Title = &p.Title
	this.Content = &p.Content
	this.IsMobile = p.IsMobile
	this.Sticky = p.Sticky
	this.CommentCount = p.CommentCount
	this.ViewCount = p.ViewCount
	this.LastReplayAt = &p.LastReplayAt
	this.CreatedAt = &p.CreatedAt
	this.UpdatedAt = &p.UpdatedAt
}

type PostView struct {
	IsLogin bool
	Post    PostDetail
	Author  Person
}
/*</used for Post detail> */

/*follow person*/
type PersonFollow struct {
	Person
	Bio string
}
type AllFollows struct {
	Following *[]PersonFollow
	Followed  *[]PersonFollow
}
/*find all person who is following */
/*find all person who is followed by*/
func findFollowsById(id uint) *AllFollows {
	db_following := []m.Follow{}
	db_followed := []m.Follow{}
	//todo Limit(256).
	database.O.QueryTable("follows").Filter("follower_id", id).Limit(256).RelatedSel().All(&db_following)
	//database.DB.Where("follower_id = ?", id).Limit(256).Preload("Following").Preload("Following.Profile").Find(&db_following)
	personFollowing := make([]PersonFollow, 0, len(db_following))
	for _, follow := range db_following {
		personFollowing = append(personFollowing, PersonFollow{Bio:follow.Following.Profile.Bio,
			Person:Person{ID:follow.Following.Id, Name:follow.Following.Name, Avatar:follow.Following.Profile.Avatar}});
	}
	//todo Limit(256)
	database.O.QueryTable("follows").Filter("following_id", id).Limit(256).RelatedSel().All(&db_followed)
	//database.DB.Where("following_id = ?", id).Limit(256).Preload("Follower").Preload("Follower.Profile").Find(&db_followed)
	personFollowed := make([]PersonFollow, 0, len(db_followed))  //dbHotPosts to mHotPosts
	for _, follow := range db_followed {
		personFollowed = append(personFollowed, PersonFollow{Bio:follow.Follower.Profile.Bio,
			Person:Person{ID:follow.Follower.Id, Name:follow.Follower.Name, Avatar:follow.Follower.Profile.Avatar}});
	}
	return &AllFollows{Following:&personFollowing, Followed:&personFollowed}
}

type Notification struct {
	ID          uint
	RelatedID   uint
	SubjectType int
	Data        string
	IsRead      bool
}

type PostMessage struct {
	ID              uint
	RelatedID       uint
	SubjectType     int
	RelatedUsername string
	PostID          uint
	PostTitle       string
	IsRead          bool
}

func findLatestNotifications(uid uint, types []int) ([]Notification) {
	//todo user id
	notices := []m.Notification{}
	database.O.QueryTable("notification").Filter("subject_type__in",types).All(&notices)
	//database.DB.Where("subject_type in (?)", types).Find(&notices)
	n := make([]Notification, 0, len(notices))
	for _, no := range notices {
		n = append(n, Notification{ID:no.Id, RelatedID:no.RelatedID,
			SubjectType:no.SubjectType, Data:no.Data, IsRead:no.IsRead});
	}
	return n
}

func findLatestPostMessages(uid uint, types []int) ([]PostMessage) {
	//todo user id
	db_messages := []m.PostMessage{}
	database.O.QueryTable("post_message").Filter("subject_type__in",types).All(&db_messages)
	//database.DB.Where("subject_type in (?)", types).Find(&db_messages)
	messages := make([]PostMessage, 0, len(db_messages))
	for _, msg := range db_messages {
		messages = append(messages, PostMessage{ID:msg.Id, RelatedID:msg.RelatedID, SubjectType:msg.SubjectType,
			RelatedUsername:msg.RelatedUsername, PostID:msg.PostID, PostTitle:msg.PostTitle, IsRead:msg.IsRead});
	}
	return messages
}