package controllers

import (
	//"../models/forms"
)

type AboutController struct {
	BaseController
}

var about_rules = map[string]int{
}

func (this *AboutController) getRules(action string) int {
	return about_rules[action]
}

func (this *AboutController) Index() {
	this.TplName = "about/index.html"
}

func (this *AboutController) Feedback() {
	this.TplName = "about/feedback.html"
}

func (this *AboutController)POST_Feedback(){
//forms.PostResult{}
}
