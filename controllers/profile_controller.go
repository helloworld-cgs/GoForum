package controllers

import (
	"strconv"
	"encoding/json"
	"../models/forms"
	"./../middleware/event"
	identify "../middleware/values"
)

type ProfileController struct {
	BaseController
}

var profile_rules = map[string]int{
	"Follow": identify.Login | identify.JumpBack,
	"Collection": identify.Login | identify.JumpBack,
}

func (this *ProfileController) getRules(action string) int {
	return profile_rules[action]
}

func (this *ProfileController) Person() {
	id, _ := strconv.Atoi(this.Ctx.Input.Param(":uid"))
	uid := uint(id)
	profile := GetProfileById(this.getUserId(),uid)   //todo use NewRecord
	if profile.Profile.UserRefer != 0 {
		json, err := json.Marshal(profile)
		if err == nil {
			this.Data["person"] = string(json)
			this.TplName = "profile/person.html"
			return
		}
	}
	this.Abort("404")
}

/*those persons i'am focusing or be followed */
func (this *ProfileController) Follow() {
	follows := findFollowsById(this.getUserId())
	json, err := json.Marshal(follows)
	if err == nil {
		this.Data["follows"] = string(json)
		this.TplName = "profile/follow.html"
		return
	}
	this.Abort("404")
}

func (this *ProfileController) FollowAdd() {
	var result *forms.PostResult
	if !this.IsUserLogin() {
		result = &forms.PostResult{Status:3, Error:"用户未登录"}
	} else {
		id, err := this.GetInt("id")
		if err == nil {
			faf := forms.FollowAddForm{PersonID:uint(id)}
			myID := this.getUserId()
			create_result, is_created := faf.Add(myID)
			result = create_result
			if is_created {
				event.OnFollowed(uint(id), myID, this.getUsername())
			}
		} else {
			result = &forms.PostResult{Status:0, Error:"ID不合法"}
		}
	}
	this.Data["json"] = result
	this.ServeJSON()
}

func (this *ProfileController) Collection() {
	this.TplName = "profile/collection.html"
}