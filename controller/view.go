/*
 * 'View' Provides basic functionality for constructing a view using Go's template engine.
 */
package controller

import (
	"AfkChampFrontend/model/user"
	"html/template"
	"net/http"
)

var TemplateMapping map[string]*template.Template

type BaseTemplateData struct {
	WebsiteName  string
	CurrentUser  *user.UserEntry
	UserLoggedIn bool
}

func init() {
	TemplateMapping = make(map[string]*template.Template)

	// Store a mapping of all templates here
	TemplateMapping["main/home.html"] = template.Must(template.New("main/home.html").Delims("<<", ">>").ParseFiles("html/base.html", "html/main/home.html"))
	TemplateMapping["about/about.html"] = template.Must(template.New("about/about.html").Delims("<<", ">>").ParseFiles("html/base.html", "html/about/about.html"))
	TemplateMapping["login/login.html"] = template.Must(template.New("login/login.html").Delims("<<", ">>").ParseFiles("html/base.html", "html/login/login.html"))
	TemplateMapping["login/register.html"] = template.Must(template.New("login/register.html").Delims("<<", ">>").ParseFiles("html/base.html", "html/login/register.html"))

	// Admin Templates
	TemplateMapping["admin/admin.html"] = template.Must(template.New("login/register.html").Delims("<<", ">>").ParseFiles("html/adminBase.html", "html/admin/admin.html"))

	// Admin Games
	TemplateMapping["admin/game/index.html"] = template.Must(template.New("admin/game/index.html").Delims("<<", ">>").ParseFiles("html/adminBase.html", "html/admin/game/index.html"))
	TemplateMapping["admin/game/newedit.html"] = template.Must(template.New("admin/game/newedit.html").Delims("<<", ">>").ParseFiles("html/adminBase.html", "html/admin/game/newedit.html"))

	// Admin Events
	TemplateMapping["admin/event/index.html"] = template.Must(template.New("admin/event/index.html").Delims("<<", ">>").ParseFiles("html/adminBase.html", "html/admin/event/index.html"))
	TemplateMapping["admin/event/newedit.html"] = template.Must(template.New("admin/event/newedit.html").Delims("<<", ">>").ParseFiles("html/adminBase.html", "html/admin/event/newedit.html"))

	// Admin Users
	TemplateMapping["admin/user/index.html"] = template.Must(template.New("admin/user/index.html").Delims("<<", ">>").ParseFiles("html/adminBase.html", "html/admin/user/index.html"))
	TemplateMapping["admin/user/newedit.html"] = template.Must(template.New("admin/user/newedit.html").Delims("<<", ">>").ParseFiles("html/adminBase.html", "html/admin/user/newedit.html"))

	// Events
	TemplateMapping["event/main.html"] = template.Must(template.New("event/main.html").Delims("<<", ">>").ParseFiles("html/base.html", "html/event/main.html"))
    
    // Matches
    TemplateMapping["match/match.html"] = template.Must(template.New("match/match.html").Delims("<<", ">>").ParseFiles("html/base.html", "html/match/match.html"))
}

func CreateTemplateData(w http.ResponseWriter, r *http.Request) BaseTemplateData {
	newTempData := BaseTemplateData{WebsiteName: "Raid Boss Down"}
	currentUser, err := GetCurrentUser(w, r, false)
	newTempData.CurrentUser = currentUser
	newTempData.UserLoggedIn = (err != nil)
	return newTempData
}
