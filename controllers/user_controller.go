package controllers

import (
	"html/template"
	"github.com/astaxie/beego"
	"./../models/forms"
	form_check "../middleware/form"
	identify "../middleware/auth"
)

type UserController struct {
	BaseController
}

var rules = map[string]int{
	"Login":   0,
	"Logout": identify.Login,
}

type SignResult struct {
	Next   string `json:"next"`
	Status bool   `json:"status"`
}

func (this *UserController) getRules(action string) int {
	return rules[action]
}

func (this *UserController) SignIn() {
	if (this.IsUserLogin()) {
		//if has login,then go home
		this.Redirect("/", 302)
		return
	}
	this.Data["xsrf_token"] = template.HTML(this.XSRFFormHTML())
	this.TplName = "account/signin.html"
}

func (this *UserController) POST_SignIn() {
	if (this.IsUserLogin()) {
		this.Redirect("/", 302)
		return
	}
	username := this.GetString("username")
	password := this.GetString("password")
	sign_in_form := forms.SignInForm{Username:username, Password: password}
	if errs, userID := sign_in_form.LoginVerify(); errs == nil {
		//验证通过
		this.LoginUser(userID, sign_in_form.Username)
		next := this.GetString("next")
		if len(next) > 0 && next[0] != '/' {
			next = "/" + next
		}
		this.Data["json"] = &SignResult{Status:true, Next:next}
	} else {
		s := form_check.NewInstant(errs, map[string]string{"name":  username, "pass": ""})
		this.Data["json"] = &s
	}
	this.ServeJSON()
}

func (this *UserController) SignUp() {
	if (this.IsUserLogin()) {
		this.Redirect("/", 302)
		return
	}
	this.Data["xsrf_token"] = template.HTML(this.XSRFFormHTML())
	this.TplName = "account/signup.html"
}

func (this *UserController) POST_SignUp() {
	if (this.IsUserLogin()) {
		this.Redirect("/", 302)
		return
	}
	email := this.GetString("email")
	nickname := this.GetString("nickname")
	password := this.GetString("password")
	sign_up_form := forms.SignUpForm{Email:email, Nickname:nickname, Password: password}
	if errs, status := sign_up_form.Valid(); status {
		this.Data["json"] = &SignResult{Status:true}
		flash := beego.NewFlash()
		flash.Success(email)
		flash.Store(&this.Controller)
	} else {
		s := form_check.NewInstant(errs, map[string]string{"email":  email, "password": ""})
		this.Data["json"] = &s
	}
	this.ServeJSON()
}

func (this *UserController) SignUpSuccess() {
	flash := beego.ReadFromRequest(&this.Controller)
	if s, ok := flash.Data["success"]; ok {
		this.Data["email"] = s
		this.TplName = "account/signup_success.html"
	} else {
		this.Abort("404")
	}
}

func (this *UserController) SignOut() {
	this.LogoutUser()
	this.Redirect("/account/signin", 302)
}